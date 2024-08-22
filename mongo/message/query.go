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
	DB          = "$db"
	Filter      = "filter"
	Documents   = "documents"
	KillCursors = "killCursors"
)

// Query represents a message query.
type Query struct {
	database   string
	collection string
	typ        string
	conditions []bson.Document
	documents  []bson.Document
	operator   string
	limit      int
}

// NewQuery returns a new query.
func NewQuery() *Query {
	q := &Query{
		database:   "",
		collection: "",
		typ:        "",
		conditions: make([]bson.Document, 0),
		documents:  make([]bson.Document, 0),
		operator:   "",
		limit:      0,
	}
	return q
}

// Database returns the database name.
func (q *Query) Database() string {
	return q.database
}

// Collection returns the collection name.
func (q *Query) Collection() string {
	return q.collection
}

// Documents returns the query documents.
func (q *Query) Documents() []bson.Document {
	return q.documents
}

// Type returns the query type.
func (q *Query) Type() string {
	return q.typ
}

// FullCollectionName returns the full collection name.
func (q *Query) FullCollectionName() string {
	return fmt.Sprintf("%s.%s", q.database, q.collection)
}

// Conditions returns the search conditions.
func (q *Query) Conditions() []bson.Document {
	return q.conditions
}

// HasConditions returns true if the query has conditions.
func (q *Query) HasConditions() bool {
	if len(q.conditions) == 0 {
		return false
	}
	for _, cond := range q.conditions {
		condElems, err := cond.Elements()
		if err != nil {
			return false
		}
		if 0 < len(condElems) {
			return true
		}
	}
	return false
}

// Operator returns the operator string.
func (q *Query) Operator() string {
	return q.operator
}

// Limit returns the limit value.
func (q *Query) Limit() int {
	return q.limit
}
