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
	"github.com/cybergarage/go-mongo/mongo/message"
)

//////////////////////////////////////////////////
// DatabaseCommandExecutor
//////////////////////////////////////////////////

// Hello displays information about this member’s role in the replica set, including whether it is the master.
func (server *Server) Hello(conn *Conn, cmd *Command) (bson.Document, error) {
	reply, err := message.NewIsMasterResponseWithConfig(server)
	if err != nil {
		return nil, err
	}

	for _, elem := range cmd.Elements {
		key := elem.Key()
		switch key {
		case message.SASLSupportedMechs:
			collection, ok := elem.Value().StringValueOK()
			if !ok {
				collection = ""
			}
			mechs, err := server.SASLSupportedMechs(conn, collection)
			if err != nil {
				return nil, err
			}
			supportedMechs := []any{}
			for _, mech := range mechs {
				supportedMechs = append(supportedMechs, mech.Name())
			}
			reply.SetArrayElements(message.SASLSupportedMechs, supportedMechs)
		}
	}

	replyDoc, err := reply.BSONBytes()
	if err != nil {
		return nil, err
	}
	return replyDoc, nil
}

// BuildInfo returns statistics about the MongoDB build.
func (server *Server) BuildInfo(conn *Conn, cmd *Command) (bson.Document, error) {
	reply, err := message.NewBuildInfoResponseWithConfig(server)
	if err != nil {
		return nil, err
	}
	replyDoc, err := reply.BSONBytes()
	if err != nil {
		return nil, err
	}
	return replyDoc, nil
}
