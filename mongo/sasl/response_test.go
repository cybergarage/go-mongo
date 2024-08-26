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

package sasl

import (
	"fmt"
	"testing"

	util "github.com/cybergarage/go-mongo/mongo/bson"
	"go.mongodb.org/mongo-driver/bson"
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
			mechs          []any
			conversationID int32
			payload        string
		}{
			{[]any{"PLAIN"}, 1, "abc"},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("%d %s", test.conversationID, test.payload), func(t *testing.T) {
				res, err := NewServerFirstResponse(test.mechs, test.conversationID, []byte(test.payload))
				if err != nil {
					t.Error(err)
					return
				}
				if res == nil {
					t.Error("NewServerFirstResponse is nil")
					return
				}
				bsonBytes, err := res.BSONBytes()
				if err != nil {
					t.Error(err)
					return
				}
				var saslResp sasleResponse
				err = bson.Unmarshal(bsonBytes, &saslResp)
				if err != nil {
					t.Errorf("unmarshal error: %s", err)
					return
				}
				if string(saslResp.Payload) != test.payload {
					if json, err := util.DocumentToJSONString(bsonBytes); err == nil {
						t.Logf("\n%s", json)
						t.Logf("\n%v", saslResp)
					}
					t.Errorf("payload error: %s != %s", string(saslResp.Payload), test.payload)
					return
				}
				if saslResp.ConversationID != int(test.conversationID) {
					t.Errorf("conversationID error: %d != %d", saslResp.ConversationID, test.conversationID)
					if json, err := util.DocumentToJSONString(bsonBytes); err == nil {
						t.Logf("\n%s", json)
						t.Logf("\n%v", saslResp)
					}
					return
				}
			})
		}
	})

}
