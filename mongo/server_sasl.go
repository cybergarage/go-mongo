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
	"strconv"

	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-sasl/sasl"
)

const (
	saslSupportedMechs           = "saslSupportedMechs"
	saslStart                    = "saslStart"
	saslMechanism                = "mechanism"
	saslPayload                  = "payload"
	saslAutoAuthorize            = "autoAuthorize"
	saslOptions                  = "options"
	saslSkipEmptyExchange        = "skipEmptyExchange"
	saslContinue                 = "saslContinue"
	saslConversationId           = "conversationId" // nolint:stylecheck
	saslDone                     = "saslDone"
	saslSpececulativAuthenticate = "speculativeAuthenticate"
)

// ExecuteSaslStart handles SASLStart command.
func (server *Server) ExecuteSaslStart(conn *Conn, cmd *Command) (bson.Document, error) {
	var reqMech string
	var reqPayload []byte
	var ok bool
	for _, elem := range cmd.Elements {
		key := elem.Key()
		switch key {
		case saslMechanism:
			reqMech, ok = elem.Value().StringValueOK()
			if !ok {
				return nil, NewErrorCommand(cmd)
			}
		case saslPayload:
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

	res, err := ctx.Next(sasl.Payload(reqPayload))
	if err != nil {
		return nil, err
	}

	conn.SetSASLContext(ctx)

	// Response to the client

	doc := bson.DocumentStart()

	mechsDoc := bson.ArrayStart()
	for n, mech := range server.Mechanisms() {
		idxKey := strconv.Itoa(n)
		mechsDoc = bson.AppendStringElement(mechsDoc, idxKey, mech.Name())
	}
	mechsDoc, err = bson.ArrayEnd(mechsDoc)
	if err != nil {
		return nil, err
	}
	doc = bson.AppendDocumentElement(doc, saslSupportedMechs, mechsDoc)

	conversationID := server.saslCounter.Inc()
	ctx.SetValue(saslConversationId, conversationID)
	doc = bson.AppendInt32Element(doc, saslConversationId, conversationID)

	doc = bson.AppendDocumentElement(doc, saslPayload, res.Bytes())

	_, err = bson.DocumentEnd(doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
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
		case saslConversationId:
			clientConversationID, ok = elem.Value().Int32OK()
			if !ok {
				return nil, NewErrorCommand(cmd)
			}
		case saslPayload:
			var t byte
			t, reqPayload, ok = elem.Value().BinaryOK()
			if !ok || t != 0 {
				return nil, NewErrorCommand(cmd)
			}
		}
	}

	// Check the conversation ID

	v, ok := ctx.Value(saslConversationId)
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

	res, err := ctx.Next(sasl.Payload(reqPayload))
	if err != nil {
		return nil, err
	}

	doc := bson.DocumentStart()

	doc = bson.AppendInt32Element(doc, saslConversationId, conversationID)
	doc = bson.AppendDocumentElement(doc, saslPayload, res.Bytes())

	_, err = bson.DocumentEnd(doc)
	if err != nil {
		return nil, err
	}

	conn.SetSASLContext(nil)

	return doc, nil
}
