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
	"strconv"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// Datetime is UTC milliseconds since the Unix epoch.
type Datetime int64

// Dictionary represents a simple BSON document.
type Dictionary struct {
	elements map[string]any
}

// NewDictionary returns a new dictionary instance.
func NewDictionary() *Dictionary {
	dict := &Dictionary{
		elements: map[string]any{},
	}
	return dict
}

// SetBooleanElement sets a boolean element.
func (dict *Dictionary) SetBooleanElement(key string, element bool) {
	dict.elements[key] = element
}

// SetInt32Element sets a int32 element.
func (dict *Dictionary) SetInt32Element(key string, element int32) {
	dict.elements[key] = element
}

// SetInt64Element sets a int64 element.
func (dict *Dictionary) SetInt64Element(key string, element int64) {
	dict.elements[key] = element
}

// SetDoubleElement sets a int64 element.
func (dict *Dictionary) SetDoubleElement(key string, element float64) {
	dict.elements[key] = element
}

// SetStringElement sets a string element.
func (dict *Dictionary) SetStringElement(key string, element string) {
	dict.elements[key] = element
}

// SetDatetimeElement sets a datetime element.
func (dict *Dictionary) SetDatetimeElement(key string, element Datetime) {
	dict.elements[key] = element
}

// SetDocumentElement sets a document element.
func (dict *Dictionary) SetDocumentElement(key string, element Document) {
	dict.elements[key] = element
}

// SetNullElement sets a null element.
func (dict *Dictionary) SetNullElement(key string) {
	dict.elements[key] = nil
}

// SetArrayElements sets an array elements.
func (dict *Dictionary) SetArrayElements(key string, elements []any) {
	dict.elements[key] = elements
}

// SetDictionaryElements sets a dictionary elements.
func (dict *Dictionary) SetDictionaryElements(key string, elements map[string]any) error {
	elemDict := NewDictionary()
	err := elemDict.SetElements(elements)
	if err != nil {
		return err
	}
	dict.elements[key] = elemDict
	return nil
}

// SetElements sets elements.
func (dict *Dictionary) SetElements(elements map[string]any) error {
	for key, element := range elements {
		switch v := element.(type) {
		case int32:
			dict.SetInt32Element(key, v)
		case int64:
			dict.SetInt64Element(key, v)
		case bool:
			dict.SetBooleanElement(key, v)
		case string:
			dict.SetStringElement(key, v)
		case float64:
			dict.SetDoubleElement(key, v)
		case Datetime:
			dict.SetDatetimeElement(key, v)
		case Document:
			dict.SetDocumentElement(key, v)
		case []byte:
			dict.SetDocumentElement(key, v)
		case []any:
			dict.SetArrayElements(key, v)
		case map[string]any:
			if err := dict.SetDictionaryElements(key, v); err != nil {
				return err
			}
		case nil:
			dict.SetNullElement(key)
		default:
			return fmt.Errorf(errorDictionaryNotSupportedType, key, element)
		}
	}
	return nil
}

// See : BSON Specification Version 1.1
// http://bsonspec.org/spec.html
// BSONBytes returns a BSON document of the dictionary.
func (dict *Dictionary) BSONBytes() (Document, error) {
	var err error
	elementBytes := make([]byte, 0)
	for key, element := range dict.elements {
		switch v := element.(type) {
		case int32:
			elementBytes = AppendInt32Element(elementBytes, key, v)
		case int64:
			elementBytes = AppendInt64Element(elementBytes, key, v)
		case bool:
			elementBytes = AppendBooleanElement(elementBytes, key, v)
		case string:
			elementBytes = AppendStringElement(elementBytes, key, v)
		case float64:
			elementBytes = AppendDoubleElement(elementBytes, key, v)
		case Datetime:
			elementBytes = AppendDateTimeElement(elementBytes, key, int64(v))
		case Document:
			elementBytes = AppendDocumentElement(elementBytes, key, v)
		case []byte:
			elementBytes = AppendDocumentElement(elementBytes, key, v)
		case []any:
			var arrayIndex int32
			arrayIndex, elementBytes = bsoncore.AppendArrayElementStart(elementBytes, key)
			for n, iv := range v {
				elementBytes, err = AppendInterfaceElement(elementBytes, strconv.Itoa(n), iv)
				if err != nil {
					return nil, err
				}
			}
			elementBytes, err = bsoncore.AppendArrayEnd(elementBytes, arrayIndex)
			if err != nil {
				return nil, err
			}
		case nil:
			elementBytes = AppendNullElement(elementBytes, key)
		default:
			return nil, fmt.Errorf(errorDictionaryNotSupportedType, key, v)
		}
	}

	documentLength := 4 + len(elementBytes) + 1
	documentBytes := AppendInt32(make([]byte, 0), int32(documentLength))
	documentBytes = append(documentBytes, elementBytes...)
	documentBytes = append(documentBytes, 0x00)

	return documentBytes, nil
}
