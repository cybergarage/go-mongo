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

// NewDefaultLastErrorResponse returns a default response instance.
func NewDefaultLastErrorResponse() *Response {
	defaultElements := map[string]interface{}{
		"err": nil,
		"n":   int32(0),
	}

	res := NewResponseWithElements(defaultElements)
	res.SetStatus(true)

	return res
}
