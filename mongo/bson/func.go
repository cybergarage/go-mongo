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
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// StartDocument returns a new document which has only the reserved length header.
func StartDocument() Document {
	_, bytes := bsoncore.AppendDocumentStart(nil)
	return bytes
}

// EndDocument writes the null byte for a document and updates the length of the document.
func EndDocument(dst []byte) ([]byte, error) {
	return bsoncore.AppendDocumentEnd(dst, 0)
}

// AppendInt32 appends an i32 value to dst and return the extended buffer.
func AppendInt32(dst []byte, value int32) []byte {
	return bsoncore.AppendInt32(dst, value)
}

// AppendCString appends a string as a cstring to dst and return the extended buffer.
func AppendCString(dst []byte, value string) []byte {
	return bsoncore.AppendKey(dst, value)
}

// AppendBooleanElement appends a boolean element to dst and return the extended buffer.
func AppendBooleanElement(dst []byte, key string, value bool) []byte {
	return bsoncore.AppendBooleanElement(dst, key, value)
}

// AppendInt32Element appends an i32 element to dst and return the extended buffer.
func AppendInt32Element(dst []byte, key string, value int32) []byte {
	return bsoncore.AppendInt32Element(dst, key, value)
}

// AppendInt64Element appends an i64 element to dst and return the extended buffer.
func AppendInt64Element(dst []byte, key string, value int64) []byte {
	return bsoncore.AppendInt64Element(dst, key, value)
}

// AppendDoubleElement appends a float64 element to dst and return the extended buffer.
func AppendDoubleElement(dst []byte, key string, value float64) []byte {
	return bsoncore.AppendDoubleElement(dst, key, value)
}

// AppendStringElement appends a string element to dst and return the extended buffer.
func AppendStringElement(dst []byte, key string, value string) []byte {
	return bsoncore.AppendStringElement(dst, key, value)
}

// AppendDateTimeElement appends an i64 datetime to dst and return the extended buffer.
func AppendDateTimeElement(dst []byte, key string, value int64) []byte {
	return bsoncore.AppendDateTimeElement(dst, key, value)
}

// AppendDocumentElement appends a document to dst and return the extended buffer.
func AppendDocumentElement(dst []byte, key string, value Document) []byte {
	return bsoncore.AppendDocumentElement(dst, key, value)
}

// AppendNullElement will append a BSON null element using key to dst and return the extended buffer.
func AppendNullElement(dst []byte, key string) []byte {
	return bsoncore.AppendNullElement(dst, key)
}

// AppendInterfaceElement will append an interface element using key to dst and return the extended buffer.
func AppendInterfaceElement(dst []byte, key string, ivalue interface{}) ([]byte, error) {
	switch value := ivalue.(type) {
	case bool:
		return bsoncore.AppendBooleanElement(dst, key, value), nil
	case int32:
		return bsoncore.AppendInt32Element(dst, key, value), nil
	case int64:
		return bsoncore.AppendInt64Element(dst, key, value), nil
	case float64:
		return bsoncore.AppendDoubleElement(dst, key, value), nil
	case string:
		return bsoncore.AppendStringElement(dst, key, value), nil
	case Document:
		return bsoncore.AppendDocumentElement(dst, key, value), nil
	case nil:
		return bsoncore.AppendNullElement(dst, key), nil
	}
	return dst, fmt.Errorf("Unkown elment type : %v", ivalue)
}

// AppendValueElement will append a value using key to dst and return the extended buffer.
func AppendValueElement(dst []byte, key string, value Value) ([]byte, error) {
	switch value.Type {
	case bsontype.Boolean:
		return bsoncore.AppendBooleanElement(dst, key, value.Boolean()), nil
	case bsontype.Int32:
		return bsoncore.AppendInt32Element(dst, key, value.Int32()), nil
	case bsontype.Int64:
		return bsoncore.AppendInt64Element(dst, key, value.Int64()), nil
	case bsontype.Double:
		return bsoncore.AppendDoubleElement(dst, key, value.Double()), nil
	case bsontype.String:
		return bsoncore.AppendStringElement(dst, key, value.StringValue()), nil
	case bsontype.EmbeddedDocument:
		return bsoncore.AppendDocumentElement(dst, key, value.Document()), nil
	case bsontype.ObjectID:
		return bsoncore.AppendObjectIDElement(dst, key, value.ObjectID()), nil
	case bsontype.Binary:
		subType, binData := value.Binary()
		return bsoncore.AppendBinaryElement(dst, key, subType, binData), nil
	case bsontype.Null:
		return bsoncore.AppendNullElement(dst, key), nil
	}
	return dst, fmt.Errorf("Unkown elment type : %v", value)
}
