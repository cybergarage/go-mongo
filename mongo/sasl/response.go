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
	"github.com/cybergarage/go-mongo/mongo/message"
)

// MongoDB Handshake
// https://github.com/mongodb/specifications/blob/master/source/mongodb-handshake/handshake.md
// specifications/source/auth/auth.md at master · mongodb/specifications
// https://github.com/mongodb/specifications/blob/master/source/auth/auth.md
// Round Trips to Authenticate a MongoDB Client Connection | ALEX BEVILACQUA
// https://alexbevi.com/blog/2020/03/03/mongodb-wire-protocol-sasl-authentication/

// NewServerFirstResponse creates a new server first response.
func NewServerFirstResponse(conversationID int32, payload []byte) (*message.Response, error) {
	spec := map[string]any{
		ConversationId: conversationID,
		Payload:        payload,
		Done:           false,
	}

	// firstMsgElements := map[string]any{
	// 	SpececulativAuthenticate: spec,
	// }

	resMsg, err := message.NewResponseWithElements(spec /*firstMsgElements*/)
	if err != nil {
		return nil, err
	}
	resMsg.SetStatus(true)
	return resMsg, nil
}

// NewServerFinalResponse creates a new server first response.
func NewServerFinalResponse(conversationID int32, payload []byte) (*message.Response, error) {
	res := message.NewResponse()
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
	return res, nil
}
