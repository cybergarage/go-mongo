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

	"github.com/cybergarage/go-mongo/bson"
)

type MyServer struct {
	*Server
	*BaseMessageHandler
	*BaseQueryExecutor
	documents []bson.Document
}

// NewMyServer returns a test server instance.
func NewMyServer() *MyServer {
	server := &MyServer{
		Server:             NewServer(),
		BaseMessageHandler: NewBaseMessageHandler(),
		BaseQueryExecutor:  NewBaseQueryExecutor(),
		documents:          make([]bson.Document, 0),
	}

	server.SetMessageListener(server)
	server.SetMessageHandler(server)
	server.SetQueryExecutor(server)

	return server
}

//////////////////////////////////////////////////
// MessageListener
//////////////////////////////////////////////////

func (server *MyServer) MessageReceived(msg OpMessage) {
	fmt.Printf("-> %s\n", msg.String())
}

func (server *MyServer) MessageRespond(msg OpMessage) {
	fmt.Printf("<- %s\n", msg.String())
}

//////////////////////////////////////////////////
// QueryExecutor
//////////////////////////////////////////////////

// Insert hadles OP_INSERT and 'insert' query of OP_MSG.
func (server *MyServer) Insert(q *Query) (int32, bool) {
	// ====================== YOUR CODE HERE ======================
	// 説明 : 指定されたドキュメント保持してください
	// ヒント :
	//　　　追加するドキュメントはQuery::GetDocuments()で取得できます。
	//　　　返り値には追加れたドキュメント数を返します
	// ============================================================

	docs := q.GetDocuments()
	if 0 < len(docs) {
		server.documents = append(server.documents, docs...)
	}

	return int32(len(docs)), true
}

// Update hadles OP_UPDATE and 'update' query of OP_MSG.
func (server *MyServer) Update(q *Query) (int32, bool) {
	// ====================== YOUR CODE HERE ======================
	// 説明 : 更新対象にに指定されたドキュメントを更新してください
	// ヒント :
	//　　　更新対象の条件はQuery::GetFilter()で取得できます。
	//　　　更新するドキュメントはQuery::GetDocuments()で取得できます。
	// ============================================================

	_ = q.GetFilter()

	return 1, true
}

// Find hadles 'find' query of OP_MSG.
func (server *MyServer) Find(q *Query) ([]bson.Document, bool) {
	// ====================== YOUR CODE HERE ======================
	// 説明 : 検索条件に指定されたドキュメントを返してください
	// ヒント : 検索条件はQuery::GetFilter()で取得できます。
	// ============================================================
	_ = q.GetFilter()

	return server.documents, true
}

// Delete hadles OP_DELETE and 'delete' query of OP_MSG.
func (server *MyServer) Delete(q *Query) (int32, bool) {
	// ====================== YOUR CODE HERE ======================
	// 説明 : 検索条件に指定されたドキュメントを返してください
	// ヒント : 検索条件はQuery::GetFilter()で取得できます。
	// ============================================================

	_ = q.GetFilter()

	return 1, true
}
