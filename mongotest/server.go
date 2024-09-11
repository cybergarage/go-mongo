// Copyright (C) 2022 The go-mongo Authors All rights reserved.
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

package mongotest

import (
	"github.com/cybergarage/go-mongo/examples/go-mongod/server"
	"github.com/cybergarage/go-mongo/mongo/sasl"
	"github.com/cybergarage/go-sasl/sasl/cred"
)

type Server struct {
	*server.Server
}

// NewServer returns a test server instance.
func NewServer() *Server {
	server := &Server{
		Server: server.NewServer(),
	}
	server.AddAuthenticator(server)
	return server
}

func (server *Server) HasCredential(username string, opts ...sasl.SASLAuthenticatorOption) (*sasl.SASLCredential, bool) {
	if username != TestUsername {
		return nil, false
	}
	passwod := TestPassword
	for _, opt := range opts {
		switch t := opt.(type) {
		case sasl.SASLMechanism:
			if t.Type().String() == "SCRAM-SHA-1" {
				var err error
				passwod, err = sasl.MongoPasswordDigest(TestUsername, TestPassword)
				if err != nil {
					return nil, false
				}
			}
		}
	}
	cred := cred.NewCredential(
		cred.WithUsername(username),
		cred.WithPassword(passwod),
	)
	return cred, true
}
