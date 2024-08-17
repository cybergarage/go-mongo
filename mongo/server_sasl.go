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
	conversationId               = "conversationId" // nolint:stylecheck
	saslDone                     = "saslDone"
	saslSpececulativAuthenticate = "speculativeAuthenticate"
)

// ExecuteSaslStart handles SASLStart command.
func (server *Server) ExecuteSaslStart(conn *Conn, cmd *Command) (bson.Document, error) {
	var reqMech string
	var reqPayload []byte
	for _, elem := range cmd.Elements {
		key := elem.Key()
		switch key {
		case saslMechanism:
			reqMech = elem.Value().StringValue()
		case saslPayload:
			reqPayload, err := elem.Value().Binary()
			if err != nil {
				return nil, err
			}
		}
	}

	mech, err := server.Mechanism(reqMech)
	if err != nil {
		return nil, err
	}

	opts := []sasl.Option{
		server.Authenticators(),
	}

	switch reqMech {
	case "SCRAM-SHA-1", "SCRAM-SHA-256", "SCRAM-SHA-512":
		opts = append(opts, sasl.Payload(reqPayload))
	case "PLAIN":
		opts = append(opts, reqPayload)
	}

	ctx, err := mech.Start(opts...)
	if err != nil {
		return nil, err
	}

	res, err := ctx.Next()
	if err != nil {
		return nil, err
	}

	conn.SetSASLContext(ctx)

	return nil, nil
}

// ExecuteSaslContinue handles SASLContinue command.
func (server *Server) ExecuteSaslContinue(*Conn, *Command) (bson.Document, error) {
	return nil, nil
}
