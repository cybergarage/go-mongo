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

// See : Update Operators
// https://docs.mongodb.com/manual/reference/operator/update/

const (
	// CurrentDate sets the value of a field to current date, either as a Date or a Timestamp.
	CurrentDate = "$currentDate"
	// Inc increments the value of the field by the specified amount.
	Inc = "$inc"
	// Min only updates the field if the specified value is less than the existing field value.
	Min = "$min"
	// Max only updates the field if the specified value is greater than the existing field value.
	Max = "$max"
	// Mul multiplies the value of the field by the specified amount.
	Mul = "$mul"
	// Rename renames a field.
	Rename = "$rename"
	// Set ets the value of a field in a document.
	Set = "$set"
	// SetOnInsert sets the value of a field if an update results in an insert of a document. Has no effect on update operations that modify existing documents.
	SetOnInsert = "$setOnInsert"
	// Unset removes the specified field from a document.
	Unset = "$unset"
)
