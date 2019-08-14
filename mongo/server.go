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
	"strconv"

	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-mongo/mongo/protocol"
)

// MessageListener represents a listener for MongoDB messages.
type MessageListener interface {
	protocol.MessageListener
}

// Server is an instance for MongoDB protocols.
type Server struct {
	*Config
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
	err := server.Stop()
	if err != nil {
		return err
	}

	err = server.open()
	if err != nil {
		return err
	}

	go server.serve()

	return nil
}

// Stop stops the server.
func (server *Server) Stop() error {
	err := server.close()
	if err != nil {
		return err
	}

	return nil
}

// Restart restarts the server.
func (server *Server) Restart() error {
	err := server.Stop()
	if err != nil {
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
	for err == nil {
		err = server.readMessage(conn)
	}

	return err
}

// nextMessageRequestID returns a next message request identifer.
func (server *Server) nextMessageRequestID() int32 {
	server.lastMessageRequestID++
	if math.MaxInt32 <= server.lastMessageRequestID {
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
func (server *Server) readMessage(conn net.Conn) error {

	headerBytes := make([]byte, protocol.HeaderSize)
	_, err := conn.Read(headerBytes)
	if err != nil {
		return err
	}

	header, err := protocol.NewHeaderWithBytes(headerBytes)
	if err != nil {
		return err
	}

	bodyBytes := make([]byte, header.GetBodySize())
	_, err = conn.Read(bodyBytes)
	if err != nil {
		return err
	}

	opMsg, err := protocol.NewMessageWithHeaderAndBytes(header, bodyBytes)
	if err != nil {
		return err
	}

	// MessageListener

	if server.messageListener != nil {
		server.messageListener.MessageReceived(opMsg)
	}

	// MessageHandler

	var resDoc bson.Document
	var resDocs []bson.Document

	if server.MessageHandler == nil {
		return fmt.Errorf(errorMessageHanderNotImplemented)
	}

	switch opMsg.GetOpCode() {
	case protocol.OpUpdate:
		msg, _ := opMsg.(*OpUpdate)
		resDoc, err = server.MessageHandler.OpUpdate(msg)
	case protocol.OpInsert:
		msg, _ := opMsg.(*OpInsert)
		resDoc, err = server.MessageHandler.OpInsert(msg)
	case protocol.OpQuery:
		msg, _ := opMsg.(*OpQuery)
		resDoc, err = server.MessageHandler.OpQuery(msg)
	case protocol.OpGetMore:
		msg, _ := opMsg.(*OpGetMore)
		resDoc, err = server.MessageHandler.OpGetMore(msg)
	case protocol.OpDelete:
		msg, _ := opMsg.(*OpDelete)
		resDoc, err = server.MessageHandler.OpDelete(msg)
	case protocol.OpKillCursors:
		msg, _ := opMsg.(*OpKillCursors)
		resDoc, err = server.MessageHandler.OpKillCursors(msg)
	case protocol.OpMsg:
		msg, _ := opMsg.(*OpMsg)
		resDoc, err = server.MessageHandler.OpMsg(msg)
	default:
		err = fmt.Errorf(errorMessageHandeUnknownOpCode, opMsg.GetOpCode())
	}

	if err != nil {
		return err
	}

	var resMsg protocol.Message

	switch opMsg.GetOpCode() {
	case protocol.OpQuery:
		reply := protocol.NewReplyWithDocument(resDoc)
		reply.SetResponseFlags(protocol.AwaitCapable)
		resMsg = reply
	case protocol.OpKillCursors:
		resMsg = protocol.NewReplyWithDocuments(resDocs)
	case protocol.OpMsg:
		msg := protocol.NewMsgWithBody(resDoc)
		resMsg = msg
	default:
		err = fmt.Errorf(errorMessageHandeUnknownOpCode, opMsg.GetOpCode())
	}

	if err != nil {
		return err
	}

	if resMsg != nil {
		resMsg.SetRequestID(server.nextMessageRequestID())
		resMsg.SetResponseTo(opMsg.GetRequestID())
		err = server.responseMessage(conn, resMsg)
	}

	return err
}
