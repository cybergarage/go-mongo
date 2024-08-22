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
	"github.com/cybergarage/go-mongo/mongo/bson"
	"github.com/cybergarage/go-mongo/mongo/protocol"
)

// NewQueryWithMessage returns a new query with the specified OP_MSG.
func NewQueryWithMessage(msg *protocol.Msg) (*Query, error) {
	q := NewQuery()
	return q, q.ParseMsg(msg)
}

// ParseMsg parses the specified OP_MSG.
func (q *Query) ParseMsg(msg *protocol.Msg) error {
	err := q.parseBodyDocument(msg.Body())
	if err != nil {
		return err
	}
	err = q.parseDocuments(msg.Documents())
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
			q.typ = key
			col, ok := element.Value().StringValueOK()
			if ok {
				q.collection = col
			}
		case DB:
			db, ok := element.Value().StringValueOK()
			if ok {
				q.database = db
			}
		case Documents:
			docs, ok := element.Value().ArrayOK()
			if ok {
				vals, err := docs.Values()
				if err != nil {
					return err
				}
				for _, val := range vals {
					doc, ok := val.DocumentOK()
					if ok {
						q.documents = append(q.documents, doc)
					}
				}
			}
		case Filter:
			doc, ok := element.Value().DocumentOK()
			if ok {
				q.conditions = append(q.conditions, doc)
			}
		}
	}

	return nil
}

// parseDocuments parses the specified BSON document.
func (q *Query) parseDocuments(docs []bson.Document) error {
	switch q.typ {
	case Insert:
		q.documents = append(q.documents, docs...)
	case Delete:
		for _, doc := range docs {
			condElem, err := doc.LookupErr("q")
			if err == nil {
				cond, ok := condElem.DocumentOK()
				if ok {
					q.conditions = append(q.conditions, cond)
				}
			}

			limitElem, err := doc.LookupErr("limit")
			if err == nil {
				limit, ok := limitElem.Int32OK()
				if ok {
					q.limit = int(limit)
				}
			}
		}
	case Update:
		for _, doc := range docs {
			condElem, err := doc.LookupErr("q")
			if err != nil {
				return err
			}
			cond, ok := condElem.DocumentOK()
			if !ok {
				continue
			}
			docElem, err := doc.LookupErr("u")
			if err != nil {
				return err
			}
			doc, ok := docElem.DocumentOK()
			if !ok {
				continue
			}
			docElems, err := doc.Elements()
			if err != nil {
				return err
			}
			if len(docElems) == 0 {
				continue
			}
			ope := docElems[0]
			opeDoc, ok := ope.Value().DocumentOK()
			if !ok {
				continue
			}
			q.operator = ope.Key()
			q.conditions = append(q.conditions, cond)
			q.documents = append(q.documents, opeDoc)
		}
	}
	return nil
}
