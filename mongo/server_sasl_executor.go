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
	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-mongo/mongo/sasl"
	"github.com/cybergarage/go-mongo/mongo/sasl/scram"
)

// MongoDB Handshake
// https://github.com/mongodb/specifications/blob/master/source/mongodb-handshake/handshake.md
// MongoDB : Authentication
// https://github.com/mongodb/specifications/blob/master/source/auth/auth.md

// SASLSupportedMechs returns the supported SASL mechanisms.
func (server *Server) SASLSupportedMechs(*Conn, string) ([]sasl.SASLMechanism, error) {
	return server.Mechanisms(), nil
}

// SASLStart handles SASLStart command.
func (server *Server) SASLStart(conn *Conn, cmd *Command) (bson.Document, error) {
	var reqMech string
	var reqPayload []byte
	var ok bool
	for _, elem := range cmd.Elements {
		key := elem.Key()
		switch key {
		case sasl.Mechanism:
			reqMech, ok = elem.Value().StringValueOK()
			if !ok {
				return nil, NewErrorCommand(cmd)
			}
		case sasl.Payload:
			var t byte
			t, reqPayload, ok = elem.Value().BinaryOK()
			if !ok || t != 0 {
				return nil, NewErrorCommand(cmd)
			}
		}
	}

	mech, err := server.Mechanism(reqMech)
	if err != nil {
		return nil, err
	}

	// Start the SASL context

	opts := []sasl.SASLOption{}

	ctx, err := mech.Start(opts...)
	if err != nil {
		return nil, err
	}

	mechRes, err := ctx.Next(sasl.SASLPayload(reqPayload))
	if err != nil {
		if !scram.IsStandardError(err) {
			return nil, err
		}
		mechRes = scram.NewMessageWithError(err)
	}

	// Response to the client

	conversationID := server.ConversationCounter().Inc()
	ctx.SetValue(sasl.ConversationId, conversationID)

	var resMsg *MessageResponse
	if err == nil {
		resMsg, err = sasl.NewServerFirstResponse(conversationID, mechRes.Bytes())
	} else {
		resMsg, err = sasl.NewServerErrorResponse(conversationID, err)
	}
	if err != nil {
		return nil, err
	}

	resDoc, err := resMsg.BSONBytes()
	if err != nil {
		return nil, err
	}

	conn.SetSASLContext(ctx)

	return resDoc, nil
}

// SASLContinue handles SASLContinue command.
func (server *Server) SASLContinue(conn *Conn, cmd *Command) (bson.Document, error) {
	ctx := conn.SASLContext()
	if ctx == nil {
		return nil, NewErrorCommand(cmd)
	}

	var clientConversationID int32
	var reqPayload []byte
	var ok bool
	for _, elem := range cmd.Elements {
		key := elem.Key()
		switch key {
		case sasl.ConversationId:
			clientConversationID, ok = elem.Value().Int32OK()
			if !ok {
				return nil, NewErrorCommand(cmd)
			}
		case sasl.Payload:
			var t byte
			t, reqPayload, ok = elem.Value().BinaryOK()
			if !ok || t != 0 {
				return nil, NewErrorCommand(cmd)
			}
		}
	}

	// Check the conversation ID

	v, ok := ctx.Value(sasl.ConversationId)
	if !ok {
		return nil, NewErrorCommand(cmd)
	}
	conversationID, ok := v.(int32)
	if !ok {
		return nil, NewErrorCommand(cmd)
	}

	if clientConversationID != conversationID {
		return nil, NewErrorCommand(cmd)
	}

	// Response to the client

	mechRes, err := ctx.Next(sasl.SASLPayload(reqPayload))
	if err != nil {
		if !scram.IsStandardError(err) {
			return nil, err
		}
	}

	var resMsg *MessageResponse
	if err == nil {
		resMsg, err = sasl.NewServerFinalResponse(conversationID, mechRes.Bytes())
	} else {
		resMsg, err = sasl.NewServerErrorResponse(conversationID, err)
	}
	if err != nil {
		return nil, err
	}

	resDoc, err := resMsg.BSONBytes()
	if err != nil {
		return nil, err
	}

	conn.SetSASLContext(nil)

	return resDoc, nil
}
