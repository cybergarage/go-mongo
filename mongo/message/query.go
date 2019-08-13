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

package message

import (
	"fmt"

	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-mongo/mongo/protocol"
)

const (
	Delete      = "delete"
	Insert      = "insert"
	Find        = "find"
	Update      = "update"
	LsID        = "lsid"
	ID          = "id"
	Binary      = "$binary"
	Base64      = "base64"
	SubType     = "subType"
	Db          = "$db"
	Filter      = "filter"
	Documents   = "documents"
	KillCursors = "killCursors"
)

// Query represents a message query.
type Query struct {
	Database   string
	Collection string
	Type       string
	Conditions []bson.Document
	Documents  []bson.Document
	Operator   string
}

// NewQuery returns a new query.
func NewQuery() *Query {
	q := &Query{
		Database:   "",
		Collection: "",
		Type:       "",
		Conditions: make([]bson.Document, 0),
		Documents:  make([]bson.Document, 0),
		Operator:   "",
	}
	return q
}

// NewQueryWithMessage returns a new query with the specified OP_MSG.
func NewQueryWithMessage(msg *protocol.Msg) (*Query, error) {
	q := NewQuery()
	return q, q.ParseMsg(msg)
}

// GetType returns the section type.
func (q *Query) GetType() string {
	return q.Type
}

// GetFullCollectionName returns the full collection name.
func (q *Query) GetFullCollectionName() string {
	return fmt.Sprintf("%s.%s", q.Database, q.Collection)
}

// GetConditions returns the search conditions.
func (q *Query) GetConditions() []bson.Document {
	return q.Conditions
}

// GetDocuments returns the search conditions.
func (q *Query) GetDocuments() []bson.Document {
	return q.Documents
}

// GetOperator returns the operator string.
func (q *Query) GetOperator() string {
	return q.Operator
}

// ParseMsg parses the specified OP_MSG.
func (q *Query) ParseMsg(msg *protocol.Msg) error {
	err := q.parseBodyDocument(msg.GetBody())
	if err != nil {
		return err
	}
	err = q.parseDocuments(msg.GetDocuments())
	if err != nil {
		return err
	}
	return nil
}

// parseBodyDocument parses the specified BSON document.
func (q *Query) parseBodyDocument(doc bson.Document) error {
	elements, err := doc.Elements()
	if err != nil {
		return err
	}
	for _, element := range elements {
		key := element.Key()
		switch key {
		case Insert, Delete, Update, Find, KillCursors:
			q.Type = key
			col, ok := element.Value().StringValueOK()
			if ok {
				q.Collection = col
			}
		case Db:
			db, ok := element.Value().StringValueOK()
			if ok {
				q.Database = db
			}
		case Filter:
			doc, ok := element.Value().DocumentOK()
			if ok {
				q.Conditions = append(q.Conditions, doc)
			}
		}
	}

	return nil
}

// parseDocuments parses the specified BSON document.
func (q *Query) parseDocuments(docs []bson.Document) error {
	switch q.Type {
	case Insert:
		q.Documents = docs
	case Update:
		for _, doc := range docs {
			updateCondElem, err := doc.LookupErr("q")
			if err != nil {
				return err
			}
			updateCond, ok := updateCondElem.DocumentOK()
			if !ok {
				continue
			}
			updateDocElem, err := doc.LookupErr("u")
			if err != nil {
				return err
			}
			updateDoc, ok := updateDocElem.DocumentOK()
			if !ok {
				continue
			}
			updateElems, err := updateDoc.Elements()
			if err != nil {
				return err
			}
			if len(updateElems) <= 0 {
				continue
			}
			ope := updateElems[0]
			opeDoc, ok := ope.Value().DocumentOK()
			if !ok {
				continue
			}
			q.Operator = ope.Key()
			q.Conditions = append(q.Conditions, updateCond)
			q.Documents = append(q.Documents, opeDoc)
		}
	}
	return nil
}
