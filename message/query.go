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

	"github.com/cybergarage/go-mongo/bson"
	"github.com/cybergarage/go-mongo/protocol"
)

const (
	Delete    = "delete"
	Insert    = "insert"
	Find      = "find"
	Update    = "update"
	LsID      = "lsid"
	ID        = "id"
	Binary    = "$binary"
	Base64    = "base64"
	SubType   = "subType"
	Db        = "$db"
	Filter    = "filter"
	Documents = "documents"
)

// Query represents a message query.
type Query struct {
	Database   string
	Collection string
	Type       string
	Filter     bson.Document
	Documents  []bson.Document
}

// NewQuery returns a new query.
func NewQuery() *Query {
	q := &Query{
		Database:   "",
		Collection: "",
		Type:       "",
		Filter:     make(bson.Document, 0),
		Documents:  make([]bson.Document, 0),
	}
	return q
}

// NewQueryWithMessage returns a new query with the specified OP_MSG.
func NewQueryWithMessage(msg *protocol.Msg) (*Query, error) {
	q := NewQuery()
	err := q.parseDocument(msg.GetBody())
	if err != nil {
		return nil, err
	}
	q.Documents = msg.GetDocuments()
	return q, nil
}

// GetType returns the section type.
func (q *Query) GetType() string {
	return q.Type
}

// GetFullCollectionName returns the full collection name.
func (q *Query) GetFullCollectionName() string {
	return fmt.Sprintf("%s.%s", q.Database, q.Collection)
}

// GetFilter returns the search filter.
func (q *Query) GetFilter() bson.Document {
	return q.Filter
}

// GetDocuments returns the search filter.
func (q *Query) GetDocuments() []bson.Document {
	return q.Documents
}

// parseDocuments parses the specified BSON documents.
func (q *Query) parseDocuments(docs []bson.Document) error {
	for _, doc := range docs {
		err := q.parseDocument(doc)
		if err != nil {
			return err
		}
	}
	return nil
}

// parseDocument parses the specified BSON document.
func (q *Query) parseDocument(doc bson.Document) error {
	elements, err := doc.Elements()
	if err != nil {
		return err
	}
	for _, element := range elements {
		key := element.Key()
		switch key {
		case Insert, Delete, Update, Find:
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
				q.Filter = doc
			}
		}
	}

	return nil
}
