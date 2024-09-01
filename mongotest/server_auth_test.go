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
	"fmt"
	"testing"

	"github.com/cybergarage/go-mongo/mongo/sasl"
)

type sasleResponse struct {
	ConversationID int    `bson:"conversationId"`
	Code           int    `bson:"code"`
	Done           bool   `bson:"done"`
	Payload        []byte `bson:"payload"`
}

func TestSASLResponses(t *testing.T) {
	t.Run("first", func(t *testing.T) {

		tests := []struct {
			mech string
			c1   string
			c2   string
		}{
			{
				"SCRAM-SHA-256",
				"n,,n=test,r=Tle5kok6ColhgwXvl72Syw9whtQXCV3K",
				"c=biws,r=Tle5kok6ColhgwXvl72Syw9whtQXCV3KCNh9jrFMXbvK21UV,p=4txwzovBCq0pFM4J3OA2iG9WBw+ClylRRRqcRwZSEiQ=",
			},
		}

		server := NewServer()

		for _, test := range tests {
			t.Run(fmt.Sprintf("%s %s", test.mech, test.c1), func(t *testing.T) {

				mech, err := server.Mechanism(test.mech)
				if err != nil {
					t.Error(err)
					return
				}

				opts := []sasl.SASLOption{
					server.Authenticators(),
				}

				ctx, err := mech.Start(opts...)
				if err != nil {
					t.Error(err)
					return
				}

				_, err = ctx.Next(sasl.SASLPayload(test.c1))
				if err != nil {
					t.Error(err)
					return
				}

			})
		}
	})

}
