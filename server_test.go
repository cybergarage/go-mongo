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
	"fmt"
	"testing"

	"github.com/cybergarage/go-mongo/bson"
)

type testServer struct {
	*Server
	*BaseMessageHandler
	*BaseQueryExecutor
	documents []bson.Document
}

// newTestServer returns a test server instance.
func newTestServer() *testServer {
	server := &testServer{
		Server:             NewServer(),
		BaseMessageHandler: NewBaseMessageHandler(),
		BaseQueryExecutor:  NewBaseQueryExecutor(),
		documents: make([]bson.Document, 0),
	}

	server.SetMessageListener(server)
	server.SetMessageHandler(server)
	server.SetQueryExecutor(server)

	return server
}

func TestNewServer(t *testing.T) {
	server := NewServer()

	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	err = server.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}

//////////////////////////////////////////////////
// MessageListener
//////////////////////////////////////////////////

func (server *testServer) MessageReceived(msg OpMessage) {
	fmt.Printf("%s\n", msg.String())
}

//////////////////////////////////////////////////
// MessageHandler
//////////////////////////////////////////////////

/*
func (server *testServer) Update(q *Update) (bson.Document, error) {
	fmt.Printf("%s\n", q.String())
	return nil, nil
}

func (server *testServer) Insert(q *Insert) (bson.Document, error) {
	fmt.Printf("%s\n", q.String())
	return nil, nil
}
*/
// Disable Query to handle by QueryExecutor
/*
func (server *testServer) Query(q *Query) error {
	fmt.Printf("%s\n", q.String())
	return nil
}
*/
/*
func (server *testServer) GetMore(q *GetMore) (bson.Document, error) {
	fmt.Printf("%s\n", q.String())
	return nil, nil
}

func (server *testServer) Delete(q *Delete) (bson.Document, error) {
	fmt.Printf("%s\n", q.String())
	return nil, nil
}

func (server *testServer) KillCursors(q *KillCursors) (bson.Document, error) {
	fmt.Printf("%s\n", q.String())
	return nil, nil
}

func (server *testServer) Msg(q *Msg) (bson.Document, error) {
	fmt.Printf("%s\n", q.String())
	return nil, nil
}
*/

//////////////////////////////////////////////////
// CommandExecutor
//////////////////////////////////////////////////

func (server *testServer) ExecuteCommand(cmd *Command) (bson.Document, error) {
	/*
		key := cmd.Key()
		fmt.Printf("%s %s\n", key, cmd.String())
		keyValue := cmd.Value()
		fmt.Printf("%s [%d] %s\n", key, keyValue.Type, keyValue.String())
	*/
	/*
		keyBytes := cmd.KeyBytes()
			doc, err := protocol.NewDocumentWithBytes(keyBytes)
			if err != nil {
				return err
			}
			elements, err := doc.Elements()
			if err != nil {
				return err
			}
			for n, element := range elements {
				fmt.Printf("  [%d] %s %s\n", n, element.Key(), element.String())
			}
	*/
	return server.BaseQueryExecutor.ExecuteCommand(cmd)
}

//////////////////////////////////////////////////
// QueryExecutor
//////////////////////////////////////////////////

// Insert hadles OP_INSERT and 'insert' query of OP_MSG.
func (server *testServer) Insert(q *Query) (int32, bool) {
	docs := q.GetDocuments()
	if 0 < len(docs) {
		server.documents = append(server.documents, docs...)
	}

	return 1, true
}

// Update hadles OP_UPDATE and 'update' query of OP_MSG.
func (server *testServer) Update(q *Query) (int32, bool) {
	return 1, true
}

// Find hadles 'find' query of OP_MSG.
func (server *testServer) Find(q *Query) ([]bson.Document, bool) {
	return server.documents, true
}

// Delete hadles OP_DELETE and 'delete' query of OP_MSG.
func (server *testServer) Delete(q *Query) (int32, bool) {
	return 1, true
}
