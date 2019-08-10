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
	"github.com/cybergarage/go-mongo/bson"
	"github.com/cybergarage/go-mongo/message"
)

//////////////////////////////////////////////////
// DatabaseCommandExecutor
//////////////////////////////////////////////////

// ExecuteIsMaster displays information about this member’s role in the replica set, including whether it is the master.
func (server *Server) ExecuteIsMaster(cmd *Command) ([]bson.Document, error) {
	reply := message.NewIsMasterResponseWithConfig(server)
	replyDoc, err := reply.BSONBytes()
	if err != nil {
		return nil, err
	}
	return []bson.Document{replyDoc}, nil
}

// ExecuteBuildInfo returns statistics about the MongoDB build.
func (server *Server) ExecuteBuildInfo(cmd *Command) ([]bson.Document, error) {
	reply := message.NewBuildInfoResponseWithConfig(server)
	replyDoc, err := reply.BSONBytes()
	if err != nil {
		return nil, err
	}
	return []bson.Document{replyDoc}, nil
}