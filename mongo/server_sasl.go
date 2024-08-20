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
	"github.com/cybergarage/go-mongo/mongo/message"
	monogo_sasl "github.com/cybergarage/go-mongo/mongo/sasl"
	"github.com/cybergarage/go-sasl/sasl"
)

// ExecuteSaslStart handles SASLStart command.
func (server *Server) ExecuteSaslStart(conn *Conn, cmd *Command) (bson.Document, error) {
	var reqMech string
	var reqPayload []byte
	var ok bool
	for _, elem := range cmd.Elements {
		key := elem.Key()
		switch key {
		case monogo_sasl.Mechanism:
			reqMech, ok = elem.Value().StringValueOK()
			if !ok {
				return nil, NewErrorCommand(cmd)
			}
		case monogo_sasl.Payload:
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

	opts := []sasl.Option{
		server.Authenticators(),
	}

	ctx, err := mech.Start(opts...)
	if err != nil {
		return nil, err
	}

	mechRes, err := ctx.Next(sasl.Payload(reqPayload))
	if err != nil {
		return nil, err
	}

	conn.SetSASLContext(ctx)

	// Response to the client

	mechs := []any{}
	for _, mech := range server.Mechanisms() {
		mechs = append(mechs, mech.Name())
	}

	conversationID := server.saslCounter.Inc()
	spec := map[string]any{
		monogo_sasl.ConversationId: conversationID,
		monogo_sasl.Payload:        mechRes.Bytes(),
		monogo_sasl.Done:           true,
	}

	firstMsgElements := map[string]any{
		monogo_sasl.SupportedMechs:           mechs,
		monogo_sasl.SpececulativAuthenticate: spec,
	}

	resMsg, err := message.NewResponseWithElements(firstMsgElements)
	if err != nil {
		return nil, err
	}
	resMsg.SetStatus(true)

	return resMsg.BSONBytes()
}

// ExecuteSaslContinue handles SASLContinue command.
func (server *Server) ExecuteSaslContinue(conn *Conn, cmd *Command) (bson.Document, error) {
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
		case monogo_sasl.ConversationId:
			clientConversationID, ok = elem.Value().Int32OK()
			if !ok {
				return nil, NewErrorCommand(cmd)
			}
		case monogo_sasl.Payload:
			var t byte
			t, reqPayload, ok = elem.Value().BinaryOK()
			if !ok || t != 0 {
				return nil, NewErrorCommand(cmd)
			}
		}
	}

	// Check the conversation ID

	v, ok := ctx.Value(monogo_sasl.ConversationId)
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

	mechRes, err := ctx.Next(sasl.Payload(reqPayload))
	if err != nil {
		return nil, err
	}

	finalMsgElements := map[string]any{
		monogo_sasl.ConversationId: conversationID,
		monogo_sasl.Payload:        mechRes.Bytes(),
		monogo_sasl.Done:           true,
	}

	resMsg, err := message.NewResponseWithElements(finalMsgElements)
	if err != nil {
		return nil, err
	}
	resMsg.SetStatus(true)

	res, err := resMsg.BSONBytes()
	if err != nil {
		return nil, err
	}

	conn.SetSASLContext(nil)

	return res, nil
}
