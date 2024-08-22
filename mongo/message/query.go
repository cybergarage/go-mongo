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
	Database   string
	Collection string
	Type       string
	Conditions []bson.Document
	Documents  []bson.Document
	Operator   string
	Limit      int
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
		Limit:      0,
	}
	return q
}

// GetDocuments returns the query conditions.
func (q *Query) GetDocuments() []bson.Document {
	return q.Documents
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

// HasConditions returns true if the query has conditions.
func (q *Query) HasConditions() bool {
	if len(q.Conditions) == 0 {
		return false
	}
	for _, cond := range q.GetConditions() {
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

// GetOperator returns the operator string.
func (q *Query) GetOperator() string {
	return q.Operator
}

// GetLimit returns the limit value.
func (q *Query) GetLimit() int {
	return q.Limit
}
