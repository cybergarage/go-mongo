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

package server

import (
	"github.com/cybergarage/go-mongo/mongo"
	"github.com/cybergarage/go-mongo/mongo/bson"
)

//////////////////////////////////////////////////
//
// Server
//
//////////////////////////////////////////////////

type Server struct {
	*mongo.Server
	documents []bson.Document
}

// NewServer returns a test server instance.
func NewServer() *Server {
	server := &Server{
		Server:    mongo.NewServer(),
		documents: make([]bson.Document, 0),
	}

	server.SetMessageListener(server)
	server.SetUserCommandExecutor(server)

	return server
}

//////////////////////////////////////////////////
//
// MessageListener
//
//////////////////////////////////////////////////

// MessageReceived passes a request message from MongoDB client
func (server *Server) MessageReceived(msg mongo.OpMessage) {
	// fmt.Printf("-> %s\n", msg.String())
	// log.Hexdump(log.LevelInfo, msg.Bytes())
}

// MessageRespond passes a response message from mongo.Server
func (server *Server) MessageRespond(msg mongo.OpMessage) {
	// fmt.Printf("<- %s\n", msg.String())
	// log.Hexdump(log.LevelInfo, msg.Bytes())
}

// Insert hadles OP_INSERT and 'insert' query of OP_MSG or OP_QUERY.
func (server *Server) Insert(q *mongo.Query) (int32, bool) {
	nInserted := int32(0)

	docs := q.GetDocuments()
	for _, doc := range docs {
		// See : The _id Field - Documents (https://docs.mongodb.com/manual/core/document/)
		docValue, err := doc.LookupErr("_id")
		if err != nil {
			continue
		}

		isInserted := false

		for _, serverDoc := range server.documents {
			serverValue, err := serverDoc.LookupErr("_id")
			if err != nil {
				continue
			}
			if serverValue.Equal(docValue) {
				isInserted = true
				break
			}
		}

		if !isInserted {
			server.documents = append(server.documents, doc)
		}

		nInserted++
	}

	if len(docs) != int(nInserted) {
		return nInserted, false
	}

	return nInserted, true
}

// Find hadles 'find' query of OP_MSG or OP_QUERY.
func (server *Server) Find(q *mongo.Query) ([]bson.Document, bool) {
	foundDoc := make([]bson.Document, 0)

	for _, doc := range server.documents {
		isMatched := true
		for _, cond := range q.GetConditions() {
			condElems, err := cond.Elements()
			if err != nil {
				return nil, false
			}
			for _, condElem := range condElems {
				docValue, err := doc.LookupErr(condElem.Key())
				if err != nil {
					isMatched = false
					break
				}
				condValue := condElem.Value()
				if !condValue.Equal(docValue) {
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

// Update hadles OP_UPDATE and 'update' query of OP_MSG or OP_QUERY.
func (server *Server) Update(q *mongo.Query) (int32, bool) {
	nUpdated := 0

	queryDocs := q.GetDocuments()
	queryConds := q.GetConditions()
	if len(queryConds) == 0 {
		return 0, true
	}

	for n := (len(server.documents) - 1); 0 <= n; n-- {
		serverDoc := server.documents[n]
		isMatched := true
		for _, cond := range q.GetConditions() {
			condElems, err := cond.Elements()
			if err != nil {
				return 0, false
			}
			for _, condElem := range condElems {
				serverFoundElem, err := serverDoc.LookupErr(condElem.Key())
				if err != nil {
					isMatched = false
					break
				}
				if !condElem.Value().Equal(serverFoundElem) {
					isMatched = false
					break
				}
			}
		}

		if !isMatched {
			continue
		}

		server.documents = append(server.documents[:n], server.documents[n+1:]...)

		serverDocElems, err := serverDoc.Elements()
		if err != nil {
			return int32(nUpdated), false
		}

		updateDoc := bson.StartDocument()
		for _, serverDocElem := range serverDocElems {
			elemKey := serverDocElem.Key()
			elemValue := serverDocElem.Value()
			for _, queryDoc := range queryDocs {
				queryValue, err := queryDoc.LookupErr(elemKey)
				if err == nil {
					elemValue = queryValue
					break
				}
			}
			updateDoc, err = bson.AppendValueElement(updateDoc, elemKey, elemValue)
			if err != nil {
				return int32(nUpdated), false
			}
		}
		updateDoc, err = bson.EndDocument(updateDoc)
		if err != nil {
			return int32(nUpdated), false
		}

		err = updateDoc.Validate()
		if err != nil {
			return int32(nUpdated), false
		}

		server.documents = append(server.documents, updateDoc)

		nUpdated++
	}

	return int32(nUpdated), true
}

// Delete hadles OP_DELETE and 'delete' query of OP_MSG or OP_QUERY.
func (server *Server) Delete(q *mongo.Query) (int32, bool) {
	nDeleted := 0

	queryConds := q.GetConditions()
	if len(queryConds) == 0 {
		nDeleted := len(server.documents)
		server.documents = make([]bson.Document, 0)
		return int32(nDeleted), true
	}

	for n := (len(server.documents) - 1); 0 <= n; n-- {
		serverDoc := server.documents[n]
		isMatched := true
		for _, cond := range q.GetConditions() {
			condElems, err := cond.Elements()
			if err != nil {
				return 0, false
			}
			for _, condElem := range condElems {
				docValue, err := serverDoc.LookupErr(condElem.Key())
				if err != nil {
					isMatched = false
					break
				}
				condValue := condElem.Value()
				if !condValue.Equal(docValue) {
					isMatched = false
					break
				}
			}
		}

		if !isMatched {
			continue
		}

		server.documents = append(server.documents[:n], server.documents[n+1:]...)
		nDeleted++
	}

	return int32(nDeleted), true
}
