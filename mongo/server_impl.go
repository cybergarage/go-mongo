// Copyright (C) 2019 The go-mongo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongo

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"sync"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mongo/mongo/auth"
	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-mongo/mongo/message"
	"github.com/cybergarage/go-mongo/mongo/protocol"
	"github.com/cybergarage/go-tracing/tracer"
)

type server struct {
	connectionWorkers sync.WaitGroup
	lifecycleMutex    sync.Mutex
	serveDone         chan struct{}
	requestIDMutex    sync.Mutex
	*config
	*ConnManager

	tlsConfig *tls.Config
	tracer.Tracer
	messageListener      MessageListener
	tcpListener          net.Listener
	MessageHandler       OpMessageHandler
	lastMessageRequestID int32
	*BaseMessageHandler
	*BaseCommandExecutor
	auth.Manager
}

// NewServer returns a new server instance.
func NewServer() Server {
	server := &server{
		config:               newDefaultConfig(),
		ConnManager:          NewConnManager(),
		tlsConfig:            nil,
		Tracer:               tracer.NullTracer,
		messageListener:      nil,
		MessageHandler:       nil,
		tcpListener:          nil,
		lastMessageRequestID: 0,
		BaseMessageHandler:   NewBaseMessageHandler(),
		BaseCommandExecutor:  NewBaseCommandExecutor(),
		Manager:              auth.NewManager(),
	}

	server.SetMessageHandler(server)
	server.SetCommandExecutor(server)
	server.SetMessageExecutor(server)
	server.SetDatabaseCommandExecutor(server)
	server.SetUserCommandExecutor(server)
	server.SetAuthCommandExecutor(server)

	return server
}

// Config returns a server configuration.
func (server *server) Config() Config {
	return server
}

// SetTracer sets a tracing tracer.
func (server *server) SetTracer(t tracer.Tracer) {
	server.Tracer = t
}

// SetMessageListener sets a message listener.
func (server *server) SetMessageListener(l MessageListener) {
	server.messageListener = l
}

// SetMessageHandler sets a message handler.
func (server *server) SetMessageHandler(h OpMessageHandler) {
	server.MessageHandler = h
}

// Start starts the server.
func (server *server) Start() error {
	server.lifecycleMutex.Lock()
	defer server.lifecycleMutex.Unlock()
	if server.tcpListener != nil {
		return errors.New("server already started")
	}
	if err := server.ConnManager.Start(); err != nil {
		return err
	}

	if server.IsTLSEnabled() {
		tlsConfig, err := server.TLSConfig()
		if err != nil {
			return err
		}
		server.tlsConfig = tlsConfig
	}

	err := server.open()
	if err != nil {
		return err
	}

	server.serveDone = make(chan struct{})
	go func(l net.Listener, done chan struct{}) {
		defer close(done)
		server.serve(l)
	}(server.tcpListener, server.serveDone)

	addr := net.JoinHostPort(server.Address(), strconv.Itoa(server.Port()))
	log.Infof("%s/%s (%s) started", PackageName, Version, addr)

	return nil
}

// Stop stops the server.
func (server *server) Stop() error {
	server.lifecycleMutex.Lock()
	defer server.lifecycleMutex.Unlock()

	// Release the listening socket even if a connection fails to close.
	listenErr := server.close()
	if server.serveDone != nil {
		<-server.serveDone
		server.serveDone = nil
	}
	connErr := server.ConnManager.Stop()
	server.connectionWorkers.Wait()
	if err := errors.Join(listenErr, connErr); err != nil {
		return err
	}

	addr := net.JoinHostPort(server.Address(), strconv.Itoa(server.Port()))
	log.Infof("%s/%s (%s) terminated", PackageName, Version, addr)

	return nil
}

// Restart restarts the server.
func (server *server) Restart() error {
	if err := server.Stop(); err != nil {
		return err
	}
	return server.Start()
}

// open opens a listen socket.
func (server *server) open() error {
	var err error
	addr := net.JoinHostPort(server.Address(), strconv.Itoa(server.Port()))
	server.tcpListener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	if server.Port() == 0 {
		addr, ok := server.tcpListener.Addr().(*net.TCPAddr)
		if !ok {
			return errors.New("listener has no TCP address")
		}
		server.SetPort(addr.Port)
	}
	return nil
}

// close closes a listening socket.
func (server *server) close() error {
	if server.tcpListener != nil {
		err := server.tcpListener.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
	}

	server.tcpListener = nil

	return nil
}

// serve handles client requests.
func (server *server) serve(l net.Listener) error {
	defer l.Close()
	for l != nil {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		if server.IsTLSEnabled() {
			conn = tls.Server(conn, server.tlsConfig)
		}
		handlerConn := newConnWith(conn, nil)
		server.AddConn(handlerConn)
		server.connectionWorkers.Go(func() {
			defer server.RemoveConn(handlerConn)
			defer handlerConn.Close()
			server.receive(handlerConn)
		})
	}

	return nil
}

// receive handles client messages.
func (server *server) receive(handlerConn *Conn) error {
	var err error
	var reqMsg, resMsg protocol.Message
	conn := handlerConn.Conn

	log.Debugf("%s/%s (%s) accepted", PackageName, Version, conn.RemoteAddr().String())

	if server.IsTLSEnabled() {
		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			return errors.New("TLS connection required")
		}
		if err := tlsConn.Handshake(); err != nil {
			return err
		}
		ok, err := server.Manager.VerifyCertificate(tlsConn)
		if !ok {
			log.Error(err)
			return err
		}
		tlsState := tlsConn.ConnectionState()
		handlerConn.tlsState = &tlsState
		conn = tlsConn
	}

	for err == nil {
		loopSpan := server.Tracer.StartSpan(PackageName)
		handlerConn.SetSpanContext(loopSpan)

		loopSpan.StartSpan("parse")
		reqMsg, err = server.readMessage(conn)
		loopSpan.FinishSpan()
		if err != nil {
			defer loopSpan.FinishSpan()
			break
		}

		resMsg, err = server.handleMessage(handlerConn, reqMsg)

		if err != nil {
			// FIXME : Check MongoDB implementation, and update to return a more standard error response
			badReply, _ := message.NewBadResponse().BSONBytes()
			switch reqMsg.OpCode() {
			case protocol.OpMsg:
				resMsg = protocol.NewMsgWithBody(badReply)
			default:
				resMsg = protocol.NewReplyWithDocument(badReply)
			}
		}

		resMsg.SetRequestID(server.nextMessageRequestID())
		resMsg.SetResponseTo(reqMsg.RequestID())

		loopSpan.StartSpan("response")
		err = server.responseMessage(conn, resMsg)
		loopSpan.FinishSpan()

		loopSpan.FinishSpan()
		if err != nil {
			break
		}
	}

	return err
}

// nextMessageRequestID returns a next message request identifier.
func (server *server) nextMessageRequestID() int32 {
	server.requestIDMutex.Lock()
	defer server.requestIDMutex.Unlock()
	server.lastMessageRequestID++
	if math.MaxInt32 == server.lastMessageRequestID {
		server.lastMessageRequestID = 0
	}
	return server.lastMessageRequestID
}

// responseMessage returns a specified message to the request connection.
func (server *server) responseMessage(conn net.Conn, msg protocol.Message) error {
	msgBytes := msg.Bytes()
	_, err := conn.Write(msgBytes)

	if server.messageListener != nil {
		server.messageListener.MessageRespond(msg)
	}

	return err
}

// readMessage handles client messages.
func (server *server) readMessage(conn net.Conn) (protocol.Message, error) {
	headerBytes := make([]byte, protocol.HeaderSize)
	nRead, err := conn.Read(headerBytes)
	if err != nil {
		if nRead <= 0 {
			return nil, err
		}
		log.Fatal(err.Error())
		return nil, err
	}

	header, err := protocol.NewHeaderWithBytes(headerBytes)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	bodyBytes := make([]byte, header.BodySize())
	_, err = conn.Read(bodyBytes)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	msg, err := protocol.NewMessageWithHeaderAndBytes(header, bodyBytes)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return msg, nil
}

// handleMessage handles client messages.
func (server *server) handleMessage(conn *Conn, reqMsg protocol.Message) (protocol.Message, error) {
	// MessageListener

	if server.messageListener != nil {
		server.messageListener.MessageReceived(reqMsg)
	}

	// MessageHandler

	var resDoc bson.Document

	if server.MessageHandler == nil {
		return nil, fmt.Errorf(errorMessageHanderNotImplemented)
	}

	var err error

	switch reqMsg.OpCode() {
	case protocol.OpUpdate:
		conn.StartSpan("OpUpdate")
		defer conn.FinishSpan()
		msg, _ := reqMsg.(*OpUpdate)
		resDoc, err = server.MessageHandler.OpUpdate(conn, msg)
	case protocol.OpInsert:
		conn.StartSpan("OpInsert")
		defer conn.FinishSpan()
		msg, _ := reqMsg.(*OpInsert)
		resDoc, err = server.MessageHandler.OpInsert(conn, msg)
	case protocol.OpQuery:
		conn.StartSpan("OpQuery")
		defer conn.FinishSpan()
		msg, _ := reqMsg.(*OpQuery)
		resDoc, err = server.MessageHandler.OpQuery(conn, msg)
	case protocol.OpGetMore:
		conn.StartSpan("OpGetMore")
		defer conn.FinishSpan()
		msg, _ := reqMsg.(*OpGetMore)
		resDoc, err = server.MessageHandler.OpGetMore(conn, msg)
	case protocol.OpDelete:
		conn.StartSpan("OpDelete")
		defer conn.FinishSpan()
		msg, _ := reqMsg.(*OpDelete)
		resDoc, err = server.MessageHandler.OpDelete(conn, msg)
	case protocol.OpKillCursors:
		conn.StartSpan("OpKillCursors")
		defer conn.FinishSpan()
		msg, _ := reqMsg.(*OpKillCursors)
		resDoc, err = server.MessageHandler.OpKillCursors(conn, msg)
	case protocol.OpMsg:
		conn.StartSpan("OpMsg")
		defer conn.FinishSpan()
		msg, _ := reqMsg.(*OpMsg)
		resDoc, err = server.MessageHandler.OpMsg(conn, msg)
	default:
		err = fmt.Errorf(errorMessageHandeUnknownOpCode, reqMsg.OpCode())
	}

	if err != nil {
		return nil, err
	}

	resDocs := []bson.Document{resDoc}

	var resMsg protocol.Message
	switch reqMsg.OpCode() {
	case protocol.OpMsg:
		msg, _ := reqMsg.(*OpMsg)
		q, err := message.NewQueryWithMessage(msg)
		if err != nil {
			return nil, err
		}
		queryType := q.Type()
		switch queryType {
		case message.Insert, message.Delete, message.Update, message.Find, message.KillCursors:
			resMsg = protocol.NewMsgWithBody(resDoc)
		default:
			resMsg = protocol.NewReplyWithDocuments(resDocs)
		}
	case protocol.OpQuery:
		replyMsg := protocol.NewReplyWithDocument(resDoc)
		replyMsg.SetResponseFlags(protocol.AwaitCapable)
		resMsg = replyMsg
	default:
		resMsg = protocol.NewReplyWithDocuments(resDocs)
	}

	return resMsg, nil
}
