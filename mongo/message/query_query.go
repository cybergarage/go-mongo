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
	"github.com/cybergarage/go-mongo/mongo/protocol"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// NewQueryWithQuery returns a new query with the specified OP_QUERY.
func NewQueryWithQuery(msg *protocol.Query) (*Query, error) {
	q := NewQuery()
	return q, q.ParseQuery(msg)
}

// ParseQuery parses the specified OP_MSG.
func (q *Query) ParseQuery(msg *protocol.Query) error {
	query := msg.Document()
	if query == nil {
		return nil
	}
	elements, err := query.Elements()
	if err != nil {
		return err
	}
	for _, element := range elements {
		key := element.Key()
		val := element.Value()
		switch key {
		case Insert, Delete, Update, Find:
			q.typ = key
			col, ok := val.StringValueOK()
			if ok {
				q.collection = col
			}
		case Documents:
			switch val.Type {
			case bsontype.Array:
				docs, ok := val.ArrayOK()
				if ok {
					docsElems, err := docs.Values()
					if err == nil {
						for _, docElem := range docsElems {
							doc, ok := docElem.DocumentOK()
							if ok {
								q.documents = append(q.documents, doc)
							}
						}
					}
				}
			case bsontype.EmbeddedDocument:
				doc, ok := val.DocumentOK()
				if ok {
					q.documents = append(q.documents, doc)
				}
			}
		case Filter:
			switch val.Type {
			case bsontype.Array:
				conds, ok := val.ArrayOK()
				if ok {
					condsElems, err := conds.Values()
					if err == nil {
						for _, condElem := range condsElems {
							cond, ok := condElem.DocumentOK()
							if ok {
								q.conditions = append(q.conditions, cond)
							}
						}
					}
				}
			case bsontype.EmbeddedDocument:
				cond, ok := val.DocumentOK()
				if ok {
					q.conditions = append(q.conditions, cond)
				}
			}
		}
	}

	return nil
}
