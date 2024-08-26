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

package bson

import (
	"bytes"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Document represents a document of MongoDB wire protocol.
type Document = bsoncore.Document

// NewDocument returns a new document.
func NewDocument() Document {
	return make(Document, 0)
}

// NewDocumentWithBytes returns a document with the specified bytes.
func NewDocumentWithBytes(src []byte) (Document, error) {
	return bsoncore.NewDocumentFromReader(bytes.NewReader(src))
}

// DocumentToJSONString returns the JSON string of the document.
func DocumentToJSONString(doc Document) (string, error) {
	decoder, err := bson.NewDecoder(bsonrw.NewBSONDocumentReader(doc))
	if err != nil {
		return "", err
	}

	var result bson.M
	err = decoder.Decode(&result)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// AppendDocumentEnd appends the end of the document.
func AppendDocumentStart(dst []byte) (int32, []byte) {
	return bsoncore.AppendDocumentStart(dst)
}

// AppendDocumentEnd appends the end of the document.
func AppendDocumentEnd(dst []byte, start int32) ([]byte, error) {
	return bsoncore.AppendDocumentEnd(dst, start)
}

// DocumentStart returns a new document which has only the reserved length header.
func DocumentStart() Document {
	_, bytes := bsoncore.AppendDocumentStart(nil)
	return bytes
}

// DocumentEnd writes the null byte for a document and updates the length of the document.
func DocumentEnd(dst []byte) ([]byte, error) {
	return bsoncore.AppendDocumentEnd(dst, 0)
}
