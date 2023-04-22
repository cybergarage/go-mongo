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
	"fmt"
	"math"
	"net"
	"os"
	"strconv"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-mongo/mongo/message"
	"github.com/cybergarage/go-mongo/mongo/protocol"
	"github.com/cybergarage/go-tracing/tracer"
)

// MessageListener represents a listener for MongoDB messages.
type MessageListener interface {
	protocol.MessageListener
}

// Server is an instance for MongoDB protocols.
type Server struct {
	*Config
	tracer.Tracer
	Addr                 string
	Port                 int
	messageListener      MessageListener
	tcpListener          net.Listener
	MessageHandler       OpMessageHandler
	lastMessageRequestID int32
	*BaseMessageHandler
	*BaseCommandExecutor
}

// NewServer returns a new server instance.
func NewServer() *Server {
	server := &Server{
		Config:               NewDefaultConfig(),
		Tracer:               nil,
		Addr:                 "",
		Port:                 DefaultPort,
		messageListener:      nil,
		MessageHandler:       nil,
		tcpListener:          nil,
		lastMessageRequestID: 0,
		BaseMessageHandler:   NewBaseMessageHandler(),
		BaseCommandExecutor:  NewBaseCommandExecutor(),
	}

	server.SetMessageHandler(server)
	server.SetCommandExecutor(server)
	server.SetMessageExecutor(server)
	server.SetDatabaseCommandExecutor(server)
	server.SetUserCommandExecutor(server)

	return server
}

// SetPort sets a listen port.
func (server *Server) SetPort(port int) {
	server.Port = port
}

// GetPort returns a listent port.
func (server *Server) GetPort() int {
	return server.Port
}

// SetTracer sets a tracing tracer.
func (server *Server) SetTracer(t tracer.Tracer) {
	server.Tracer = t
}

// SetMessageListener sets a message listener.
func (server *Server) SetMessageListener(l MessageListener) {
	server.messageListener = l
}

// SetMessageHandler sets a message handler.
func (server *Server) SetMessageHandler(h OpMessageHandler) {
	server.MessageHandler = h
}

// Start starts the server.
func (server *Server) Start() error {
	err := server.open()
	if err != nil {
		return err
	}

	go server.serve()

	log.Infof("%s/%s (PID:%d) started", PackageName, Version, os.Getpid())

	return nil
}

// Stop stops the server.
func (server *Server) Stop() error {
	if err := server.close(); err != nil {
		return err
	}

	log.Infof("%s/%s (PID:%d) terminated", PackageName, Version, os.Getpid())

	return nil
}

// Restart restarts the server.
func (server *Server) Restart() error {
	if err := server.Stop(); err != nil {
		return err
	}
	return server.Start()
}

// open opens a listen socket.
func (server *Server) open() error {
	var err error
	addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.Port))
	server.tcpListener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return nil
}

// close closes a listening socket.
func (server *Server) close() error {
	if server.tcpListener != nil {
		err := server.tcpListener.Close()
		if err != nil {
			return err
		}
	}

	server.tcpListener = nil

	return nil
}

// serve handles client requests.
func (server *Server) serve() error {
	defer server.close()

	l := server.tcpListener
	for {
		if l == nil {
			break
		}
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go server.receive(conn)
	}

	return nil
}

// receive handles client messages.
func (server *Server) receive(conn net.Conn) error {
	defer conn.Close()

	var err error
	var reqMsg, resMsg protocol.Message

	log.Debugf("%s/%s (%s) accepted", PackageName, Version, conn.RemoteAddr().String())

	handlerConn := newConn()
	for err == nil {
		reqMsg, err = server.readMessage(conn)
		if err != nil {
			break
		}
		resMsg, err = server.handleMessage(handlerConn, reqMsg)
		if err != nil {
			// FIXME : Check MongoDB implementation, and update to return a more standard error response
			badReply, _ := message.NewBadResponse().BSONBytes()
			switch reqMsg.GetOpCode() {
			case protocol.OpMsg:
				resMsg = protocol.NewMsgWithBody(badReply)
			default:
				resMsg = protocol.NewReplyWithDocument(badReply)
			}
		}

		resMsg.SetRequestID(server.nextMessageRequestID())
		resMsg.SetResponseTo(reqMsg.GetRequestID())

		err = server.responseMessage(conn, resMsg)
		if err != nil {
			break
		}
	}

	return err
}

// nextMessageRequestID returns a next message request identifier.
func (server *Server) nextMessageRequestID() int32 {
	server.lastMessageRequestID++
	if math.MaxInt32 == server.lastMessageRequestID {
		server.lastMessageRequestID = 0
	}
	return server.lastMessageRequestID
}

// responseMessage returns a specified message to the request connection.
func (server *Server) responseMessage(conn net.Conn, msg protocol.Message) error {
	msgBytes := msg.Bytes()
	_, err := conn.Write(msgBytes)

	if server.messageListener != nil {
		server.messageListener.MessageRespond(msg)
	}

	return err
}

// readMessage handles client messages.
func (server *Server) readMessage(conn net.Conn) (protocol.Message, error) {
	headerBytes := make([]byte, protocol.HeaderSize)
	nRead, err := conn.Read(headerBytes)
	if err != nil {
		if nRead <= 0 {
			return nil, err
		}
		log.Fatalf(err.Error())
		return nil, err
	}

	header, err := protocol.NewHeaderWithBytes(headerBytes)
	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}

	bodyBytes := make([]byte, header.GetBodySize())
	_, err = conn.Read(bodyBytes)
	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}

	msg, err := protocol.NewMessageWithHeaderAndBytes(header, bodyBytes)
	if err != nil {
		log.Fatalf(err.Error())
		return nil, err
	}

	return msg, nil
}

// handleMessage handles client messages.
func (server *Server) handleMessage(conn *Conn, reqMsg protocol.Message) (protocol.Message, error) {
	// MessageListener

	if server.messageListener != nil {
		server.messageListener.MessageReceived(reqMsg)
	}

	// MessageHandler

	var resDoc bson.Document
	var resDocs []bson.Document

	if server.MessageHandler == nil {
		return nil, fmt.Errorf(errorMessageHanderNotImplemented)
	}

	var err error

	switch reqMsg.GetOpCode() {
	case protocol.OpUpdate:
		msg, _ := reqMsg.(*OpUpdate)
		resDoc, err = server.MessageHandler.OpUpdate(conn, msg)
	case protocol.OpInsert:
		msg, _ := reqMsg.(*OpInsert)
		resDoc, err = server.MessageHandler.OpInsert(conn, msg)
	case protocol.OpQuery:
		msg, _ := reqMsg.(*OpQuery)
		resDoc, err = server.MessageHandler.OpQuery(conn, msg)
	case protocol.OpGetMore:
		msg, _ := reqMsg.(*OpGetMore)
		resDoc, err = server.MessageHandler.OpGetMore(conn, msg)
	case protocol.OpDelete:
		msg, _ := reqMsg.(*OpDelete)
		resDoc, err = server.MessageHandler.OpDelete(conn, msg)
	case protocol.OpKillCursors:
		msg, _ := reqMsg.(*OpKillCursors)
		resDoc, err = server.MessageHandler.OpKillCursors(conn, msg)
	case protocol.OpMsg:
		msg, _ := reqMsg.(*OpMsg)
		resDoc, err = server.MessageHandler.OpMsg(conn, msg)
	default:
		err = fmt.Errorf(errorMessageHandeUnknownOpCode, reqMsg.GetOpCode())
	}

	if err != nil {
		return nil, err
	}

	var resMsg protocol.Message

	switch reqMsg.GetOpCode() {
	case protocol.OpMsg:
		resMsg = protocol.NewMsgWithBody(resDoc)
	case protocol.OpQuery:
		replyMsg := protocol.NewReplyWithDocument(resDoc)
		replyMsg.SetResponseFlags(protocol.AwaitCapable)
		resMsg = replyMsg
	default:
		resMsg = protocol.NewReplyWithDocuments(resDocs)
	}

	return resMsg, nil
}
