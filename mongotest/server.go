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
	"github.com/cybergarage/go-sasl/sasl/auth"
)

type Server struct {
	*server.Server
}

// NewServer returns a test server instance.
func NewServer() *Server {
	server := &Server{
		Server: server.NewServer(),
	}
	server.SetCredentialStore(server)
	return server
}

// LookupCredential looks up a credential by the query.
func (server *Server) LookupCredential(q auth.Query) (auth.Credential, bool, error) {
	username := q.Username()
	if username != TestUsername {
		return nil, false, auth.ErrNoCredential
	}
	passwod := TestPassword
	mech := q.Mechanism()
	if mech == "SCRAM-SHA-1" {
		var err error
		passwod, err = sasl.MongoPasswordDigest(TestUsername, TestPassword)
		if err != nil {
			return nil, false, auth.ErrNoCredential
		}
	}
	cred := auth.NewCredential(
		auth.WithCredentialUsername(username),
		auth.WithCredentialPassword(passwod),
	)
	return cred, true, nil
}
