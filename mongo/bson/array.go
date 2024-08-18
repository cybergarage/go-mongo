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
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// AppendArrayStart returns a new array which has only the reserved length header.
func ArrayStart() []byte {
	_, bytes := bsoncore.AppendArrayStart(nil)
	return bytes
}

// AppendArrayEnd writes the null byte for an array and updates the length of the array.
func ArrayEnd(dst []byte) ([]byte, error) {
	return bsoncore.AppendArrayEnd(dst, 0)
}

// AppendArrayElement appends an element to dst and return the extended buffer.
func AppendArrayElement(dst []byte, key string, value []byte) []byte {
	return bsoncore.AppendDocumentElement(dst, key, value)
}
