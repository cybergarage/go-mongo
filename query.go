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

// QueryHandler represens a query handler fo MongoDB queries.
type QueryHandler interface {
	// Insert handles a OP_INSERT message.
	Insert(*Query) error
	// Update handles a OP_UPDATE message.
	Update(*Query) error
	// Select handles a OP_QUERY message.
	Select(*Query) ([]Document, error)
	// Delete handles a OP_DELETE message.
	Delete(*Query) error
}

// Query represents a single query of MongoDB
type Query struct {
	collectionName string
	selector       *Document
	document       *Document
}

// NewQuery returns a query instance.
func NewQuery() *Query {
	query := &Query{
		collectionName: "",
		selector:       nil,
		document:       nil,
	}

	return query
}

// SetCollectionName sets a full collection name.
func (query *Query) SetCollectionName(name string) {
	query.collectionName = name
}

// GetCollectionName returns the full collection name.
func (query *Query) GetCollectionName() string {
	return query.collectionName
}

// SetSelector sets a selector document for OP_UPDATE and OP_DELETE.
func (query *Query) SetSelector(doc *Document) {
	query.selector = doc
}

// GetSelector returns the selector document for OP_UPDATE and OP_DELETE..
func (query *Query) GetSelector() *Document {
	return query.selector
}

// SetDocument sets a data document for OP_INSERT, OP_UPDATE and OP_QUERY.
func (query *Query) SetDocument(doc *Document) {
	query.document = doc
}

// GetDocument returns the data document for OP_INSERT, OP_UPDATE and OP_QUERY.
func (query *Query) GetDocument() *Document {
	return query.document
}
