// Copyright (C) 2024 The go-mongo Authors. All rights reserved.
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
	"testing"

	"github.com/cybergarage/go-mongo/mongo/sasl"
	"github.com/cybergarage/go-sasl/sasl/scram"
	xgoscram "github.com/xdg-go/scram"
)

func SCRAMServerTest(t *testing.T) {
	t.Helper()

	sha1Client, err := xgoscram.SHA1.NewClientUnprepped(TestUsername, TestPassword, "")
	if err != nil {
		t.Error(err)
		return
	}

	sha256Client, err := xgoscram.SHA256.NewClientUnprepped(TestUsername, TestPassword, "")
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name   string
		client *xgoscram.Client
		scram.HashFunc
	}{
		{
			name:     "SCRAM-SHA-1",
			client:   sha1Client,
			HashFunc: scram.HashSHA1(),
		},
		{
			name:     "SCRAM-SHA-256",
			client:   sha256Client,
			HashFunc: scram.HashSHA256(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := NewServer()

			// Mechanism

			mech, err := server.Mechanism(test.name)
			if err != nil {
				t.Error(err)
				return
			}

			// Client first message

			conv := test.client.NewConversation()
			clientMsg, err := conv.Step("")
			if err != nil {
				t.Error(err)
				return
			}

			t.Logf("[c1] %s", clientMsg)

			// Server first message

			opts := []sasl.SASLOption{
				server.CredentialStore(),
			}

			ctx, err := mech.Start(opts...)
			if err != nil {
				t.Error(err)
				return
			}

			serverMsg, err := ctx.Next(sasl.SASLPayload(clientMsg))
			if err != nil {
				t.Error(err)
				return
			}

			t.Logf("[s1] %s", serverMsg.String())

			// Client final message

			clientMsg, err = conv.Step(serverMsg.String())
			if err != nil {
				t.Error(err)
				return
			}

			t.Logf("[c2] %s", clientMsg)

			// Server final message

			serverMsg, err = ctx.Next(sasl.SASLPayload(clientMsg))
			if err != nil {
				t.Error(err)
				return
			}

			t.Logf("[s2] %s", serverMsg.String())

			// Client validation

			_, err = conv.Step(serverMsg.String())
			if err != nil {
				t.Error(err)
				return
			}
		})
	}
}

func TestAuthMechanisms(t *testing.T) {
	t.Run("SCRAM", func(t *testing.T) {
		SCRAMServerTest(t)
	})
}
