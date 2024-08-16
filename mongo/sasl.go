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
	conversationId               = "conversationId" // nolint:stylecheck
	saslDone                     = "saslDone"
	saslSpececulativAuthenticate = "speculativeAuthenticate"
)

// ExecuteSaslStart handles SASLStart command.
func (server *Server) ExecuteSaslStart(conn *Conn, cmd *Command) (bson.Document, error) {
	var mech string
	var payload string
	for _, elem := range cmd.Elements {
		key := elem.Key()
		switch key {
		case saslMechanism:
			mech = elem.Value().StringValue()
			if len(mech) == 0 {
				return nil, nil
			}
		case saslPayload:
			payload = elem.Value().StringValue()
			if len(payload) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

// ExecuteSaslContinue handles SASLContinue command.
func (server *Server) ExecuteSaslContinue(*Conn, *Command) (bson.Document, error) {
	return nil, nil
}
