// Copyright (C) 2020 The go-mongo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongotest

import (
	"fmt"

	"github.com/cybergarage/go-mongo/mongo/shell"
)

// QueryResponseRowsKey represents a query response.
type QueryResponse struct {
	Data any
}

// QueryResponseMap represents a JSON response map data type.
type QueryResponseMap = map[string]any

// QueryResponseRows represents a JSON response array type.
type QueryResponseRows = []any

// NewQueryResponse returns a blank response instance.
func NewQueryResponse() *QueryResponse {
	return &QueryResponse{
		Data: nil,
	}
}

// NewQueryResponseWithString returns a response instance of the specified JSON response.
func NewQueryResponseWithString(json string) (*QueryResponse, error) {
	res := NewQueryResponse()
	err := res.ParseString(json)
	return res, err
}

// ParseString parses a specified string response as a JSON data.
func (res *QueryResponse) ParseString(resStr string) error {
	rootObj, err := shell.DecodeResponse(resStr)
	if err != nil {
		return err
	}
	res.Data = rootObj
	return nil
}

// String returns the string representation.
func (res *QueryResponse) String() string {
	return fmt.Sprintf("%v", res.Data)
}