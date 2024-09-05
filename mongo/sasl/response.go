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

package sasl

import (
	"fmt"

	"github.com/cybergarage/go-mongo/mongo/message"
	"github.com/cybergarage/go-mongo/mongo/sasl/scram"
)

// MongoDB Handshake
// https://github.com/mongodb/specifications/blob/master/source/mongodb-handshake/handshake.md
// Authentication - MongoDB Specifications
// https://github.com/mongodb/specifications/blob/master/source/auth/auth.md
// SASL Mechanisms
// conversationId: the conversation identifier. This will need to be remembered and used for the duration of the conversation.
// code: A response code that will indicate failure. This field is not included when the command was successful.
// done: a boolean value indicating whether or not the conversation has completed.
// payload: a sequence of bytes or a base64 encoded string (depending on input) to pass into the SASL library to transition the state machine.

// NewServerFirstResponse creates a new server first response.
func NewServerFirstResponse(conversationID int32, payload []byte) (*message.Response, error) {
	// Authentication - MongoDB Specifications
	// https://github.com/mongodb/specifications/blob/master/source/auth/auth.md
	// CMD = { saslStart: 1, mechanism: <mechanism_name>, payload: BinData(...), autoAuthorize: 1 }
	// RESP = { conversationId: <number>, code: <code>, done: <boolean>, payload: <payload> }

	firstMsgElements := map[string]any{
		ConversationId: conversationID,
		Payload:        payload,
		Done:           false,
	}

	resMsg, err := message.NewResponseWithElements(firstMsgElements)
	if err != nil {
		return nil, err
	}
	resMsg.SetStatus(true)
	return resMsg, nil
}

// NewServerFinalResponse creates a new server first response.
func NewServerFinalResponse(conversationID int32, payload []byte) (*message.Response, error) {
	// Authentication - MongoDB Specifications
	// https://github.com/mongodb/specifications/blob/master/source/auth/auth.md
	// CMD = { saslContinue: 1, conversationId: conversationId, payload: BinData(...) }
	// RESP = { conversationId: <number>, code: <code>, done: <boolean>, payload: <payload> }

	finalMsgElements := map[string]any{
		ConversationId: conversationID,
		Payload:        payload,
		Done:           true,
	}
	resMsg, err := message.NewResponseWithElements(finalMsgElements)
	if err != nil {
		return nil, err
	}
	resMsg.SetStatus(true)
	return resMsg, nil
}

// NewServerErrorResponse creates a new server error response.
func NewServerErrorResponse(conversationID int32, err error) (*message.Response, error) {
	errMsg := map[string]any{
		ConversationId: conversationID,
		Payload:        scram.NewMessageWithError(err).Bytes(),
		Done:           false,
	}
	resMsg, err := message.NewResponseWithElements(errMsg)
	if err != nil {
		return nil, err
	}
	resMsg.SetStatus(true)
	return resMsg, nil
}

// NewServerErrorlResponse creates a new server error response.
func NewServerErrorlResponse(conversationID int32, err error) (*message.Response, error) {
	res := message.NewResponse()
	errMsgElements := map[string]any{
		ConversationId: conversationID,
		Payload:        fmt.Sprintf("e=%s", err),
		Done:           false,
	}
	resMsg, err := message.NewResponseWithElements(errMsgElements)
	if err != nil {
		return nil, err
	}
	resMsg.SetStatus(false)
	return res, nil
}
