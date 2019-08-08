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

const (
	NumberOfAffectedDocuments = "n"
)

// NewMessageReplyWithParameters returns a message response instance.
func NewMessageReplyWithParameters(ok bool, n int32) *Response {
	res := NewResponse()
	res.SetNumberOfAffectedDocuments(n)
	res.SetStatus(ok)
	return res
}

// SetNumberOfAffectedDocuments sets a number of affected documents.
func (res *Response) SetNumberOfAffectedDocuments(n int32) {
	res.SetInt32Element(NumberOfAffectedDocuments, n)
}
