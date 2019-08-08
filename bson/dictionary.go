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

import "fmt"

// Datetime is UTC milliseconds since the Unix epoch.
type Datetime int64

// Dictionary represents a simple BSON document.
type Dictionary struct {
	elements map[string]interface{}
}

// NewDictionary returns a new dictionary instance.
func NewDictionary() *Dictionary {
	dict := &Dictionary{
		elements: map[string]interface{}{},
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

// SetElements sets elements.
func (dict *Dictionary) SetElements(elements map[string]interface{}) error {
	for key, element := range elements {
		switch val := element.(type) {
		case int32:
			dict.SetInt32Element(key, val)
		case int64:
			dict.SetInt64Element(key, val)
		case bool:
			dict.SetBooleanElement(key, val)
		case string:
			dict.SetStringElement(key, val)
		case float64:
			dict.SetDoubleElement(key, val)
		case Datetime:
			dict.SetDatetimeElement(key, val)
		case Document:
			dict.SetDocumentElement(key, val)
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
	elementBytes := make([]byte, 0)
	for key, element := range dict.elements {
		switch val := element.(type) {
		case int32:
			elementBytes = AppendInt32Element(elementBytes, key, val)
		case int64:
			elementBytes = AppendInt64Element(elementBytes, key, val)
		case bool:
			elementBytes = AppendBooleanElement(elementBytes, key, val)
		case string:
			elementBytes = AppendStringElement(elementBytes, key, val)
		case float64:
			elementBytes = AppendDoubleElement(elementBytes, key, val)
		case Datetime:
			elementBytes = AppendDateTimeElement(elementBytes, key, int64(val))
		case Document:
			elementBytes = AppendDocumentElement(elementBytes, key, Document(val))
		default:
			return nil, fmt.Errorf(errorDictionaryNotSupportedType, key, val)
		}
	}

	documentLength := 4 + len(elementBytes) + 1
	documentBytes := AppendInt32(make([]byte, 0), int32(documentLength))
	documentBytes = append(documentBytes, elementBytes...)
	documentBytes = append(documentBytes, 0x00)

	return documentBytes, nil
}
