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

package main

import (
	"fmt"

	"github.com/cybergarage/go-mongo/mongo"
	"github.com/cybergarage/go-mongo/mongo/bson"
)

type MyServer struct {
	*mongo.Server
	documents []bson.Document
}

// NewMyServer returns a test server instance.
func NewMyServer() *MyServer {
	server := &MyServer{
		Server:    mongo.NewServer(),
		documents: make([]bson.Document, 0),
	}

	server.SetMessageListener(server)
	server.SetUserCommandExecutor(server)

	return server
}

//////////////////////////////////////////////////
// MessageListener
//////////////////////////////////////////////////

func (server *MyServer) MessageReceived(msg mongo.OpMessage) {
	fmt.Printf("-> %s\n", msg.String())
}

func (server *MyServer) MessageRespond(msg mongo.OpMessage) {
	fmt.Printf("<- %s\n", msg.String())
}

//////////////////////////////////////////////////
// CommandExecutor
//////////////////////////////////////////////////

// Insert hadles OP_INSERT and 'insert' query of OP_MSG.
func (server *MyServer) Insert(q *mongo.Query) (int32, bool) {
	// ====================== YOUR CODE HERE ======================
	// 説明 : 指定されたドキュメント保持してください
	// ヒント :
	//　　　追加するドキュメントはQuery::GetDocuments()で取得できます。
	//　　　返り値には追加れたドキュメント数を返します
	// 参考 : insert - Database Commands > Query and Write Operation Commands
	//        https://docs.mongodb.com/manual/reference/command/insert/#dbcmd.insert
	// ============================================================

	nInserted := int32(0)

	docs := q.GetDocuments()
	for _, doc := range docs {
		// MongoDBドキュメントは主キーとなる_idフィールドが存在します
		// 参考 : The _id Field - Documents (https://docs.mongodb.com/manual/core/document/)
		docElem, err := doc.LookupErr("_id")
		if err != nil {
			continue
		}
		docID, ok := docElem.ObjectIDOK()
		if !ok {
			continue
		}

		// _idの重複により、既に追加されたドキュメントか確認します
		for _, dbDoc := range server.documents {
			dbDocID := dbDoc.Lookup("_id").ObjectID()
			if docID == dbDocID {
				continue
			}
		}

		// 新規ドキュメントを追加します
		server.documents = append(server.documents, doc)
		nInserted++
	}

	if len(docs) != int(nInserted) {
		return nInserted, false
	}

	return nInserted, true
}

// Update hadles OP_UPDATE and 'update' query of OP_MSG.
func (server *MyServer) Update(q *mongo.Query) (int32, bool) {
	// ====================== YOUR CODE HERE ======================
	// 説明 : 更新対象にに指定されたドキュメントを更新してください
	// ヒント :
	//　　　更新対象の条件はQuery::GetConditions()で取得できます。
	//　　　更新するドキュメントはQuery::GetDocuments()で取得できます。
	// 参考 : update - Database Commands > Query and Write Operation Commands
	//       https://docs.mongodb.com/manual/reference/command/update/#dbcmd.update
	// ============================================================

	conds := q.GetConditions()
	docs := q.GetDocuments()

	fmt.Printf("conds = %v\n", conds)
	fmt.Printf("docs = %v\n", docs)

	return 1, true
}

// Find hadles 'find' query of OP_MSG.
func (server *MyServer) Find(q *mongo.Query) ([]bson.Document, bool) {
	// ====================== YOUR CODE HERE ======================
	// 説明 : 検索条件に指定されたドキュメントを返してください
	// ヒント : 検索条件はQuery::GetConditions()で取得できます。
	// 参考 : find - Database Commands > Query and Write Operation Commands
	//       https://docs.mongodb.com/manual/reference/command/find/#dbcmd.find
	// ============================================================

	foundDoc := make([]bson.Document, 0)

	for _, doc := range server.documents {
		isMatched := true
		for _, cond := range q.GetConditions() {
			condElems, err := cond.Elements()
			if err != nil {
				return nil, false
			}
			for _, condElem := range condElems {
				docElem, err := doc.LookupErr(condElem.Key())
				if err != nil {
					isMatched = false
					break
				}
				if !condElem.Value().Equal(docElem) {
					isMatched = false
					break
				}
			}
		}

		if !isMatched {
			continue
		}

		foundDoc = append(foundDoc, doc)
	}

	return foundDoc, true
}

// Delete hadles OP_DELETE and 'delete' query of OP_MSG.
func (server *MyServer) Delete(q *mongo.Query) (int32, bool) {
	// ====================== YOUR CODE HERE ======================
	// 説明 : 検索条件に指定されたドキュメントを返してください
	// ヒント : 検索条件はQuery::GetConditions()で取得できます。
	// 参考 : delete - Database Commands > Query and Write Operation Commands
	//       https://docs.mongodb.com/manual/reference/command/delete/#dbcmd.delete
	// ============================================================

	conds := q.GetConditions()
	if len(conds) <= 0 {
		nDeleted := len(server.documents)
		server.documents = make([]bson.Document, 0)
		return int32(nDeleted), true
	}

	deletedDocs := append(make([]bson.Document, 0), server.documents...)
	nDeleted := 0

	for n := (len(server.documents) - 1); 0 <= n; n-- {
		doc := server.documents[n]
		isMatched := true
		for _, cond := range q.GetConditions() {
			condElems, err := cond.Elements()
			if err != nil {
				return 0, false
			}
			for _, condElem := range condElems {
				docElem, err := doc.LookupErr(condElem.Key())
				if err != nil {
					isMatched = false
					break
				}
				if !condElem.Value().Equal(docElem) {
					isMatched = false
					break
				}
			}
		}

		if !isMatched {
			continue
		}

		deletedDocs = append(deletedDocs[:n], deletedDocs[n+1:]...)
		nDeleted++
	}

	return int32(nDeleted), true
}
