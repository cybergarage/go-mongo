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
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	QueryResponseRowsKey = "rows"
)

// QueryResponseData defines a JSON response data type.
type QueryResponseData = map[string]interface{}

// QueryResponseRows defines a JSON response rows type.
type QueryResponseRows = []interface{}

// QueryResponseRow defines a JSON response row type.
type QueryResponseRow = map[string]interface{}

// QueryResponse represents a response of a query.
type QueryResponse struct {
	Data QueryResponseData
}

// NewQueryResponse returns a response instance.
func NewQueryResponse() *QueryResponse {
	res := &QueryResponse{
		Data: QueryResponseData{},
	}
	return res
}

// NewQueryResponseWithString returns a response instance of the specified JSON response.
func NewQueryResponseWithString(json string) (*QueryResponse, error) {
	res := NewQueryResponse()
	err := res.ParseString(json)
	return res, err
}

// ParseString parses a specified string response as a JSON data.
func (res *QueryResponse) ParseString(jsonStr string) error {
	var rootObj interface{}
	err := json.Unmarshal([]byte(jsonStr), &rootObj)
	if err != nil {
		return err
	}

	var ok bool
	res.Data, ok = rootObj.(QueryResponseData)
	if !ok {
		return fmt.Errorf(errorInvalidJSONResponse, rootObj)
	}

	return nil
}

// Rows returns response rows with true when the response has any rows, otherwise nil and false.
func (res *QueryResponse) Rows() (QueryResponseRows, error) {
	if res.Data == nil {
		return nil, fmt.Errorf(errorJSONResponseNotFound)
	}

	rowsData, ok := res.Data[QueryResponseRowsKey]
	if !ok {
		return nil, fmt.Errorf(errorJSONResponseRowsNotFound, res.Data, QueryResponseRowsKey)
	}

	rows, ok := rowsData.(QueryResponseRows)
	if !ok {
		return nil, fmt.Errorf(errorJSONResponseRowsNotFound, res.Data, QueryResponseRowsKey)
	}

	return rows, nil
}

// HasRow returns true when the response has a specified row, otherwise false.
// nolint: gocyclo
func (res *QueryResponse) HasRow(row interface{}) error {
	rowMap, ok := row.(QueryResponseRow)
	if !ok {
		return fmt.Errorf(errorJSONResponseHasNoRow, rowMap, row)
	}

	resRows, err := res.Rows()
	if err != nil {
		return err
	}

	var deepEqual func(iv1 interface{}, iv2 interface{}) bool
	deepEqual = func(iv1 interface{}, iv2 interface{}) bool {
		if reflect.DeepEqual(iv1, iv2) {
			return true
		}

		trimString := func(s string) string {
			return strings.Trim(s, "\"'")
		}

		// NOTE: DeepEqual checks the types and values strictly.
		// Therefore, support other types if needed.
		// log.Debugf("deepEqual: %v (%T) != %v (%T)", iv1, iv1, iv2, iv2)

		uint8sToString := func(ui8s []uint8) string {
			bytesLen := len(ui8s)
			bytes := make([]byte, bytesLen)
			for n := 0; n < bytesLen; n++ {
				bytes[n] = ui8s[n]
			}
			return string(bytes)
		}

		switch v1 := iv1.(type) {
		case string:
			switch v2 := iv2.(type) {
			case string:
				sv1 := trimString(v1)
				sv2 := trimString(v2)
				if sv1 == sv2 {
					return true
				}
			case []uint8:
				sv2 := uint8sToString(v2)
				return deepEqual(v1, sv2)
			case int:
				sv2 := strconv.Itoa(v2)
				return deepEqual(v1, sv2)
			case float64:
				sv2 := strconv.FormatFloat(v2, 'G', -1, 64)
				return deepEqual(v1, sv2)
			default:
				sv2 := fmt.Sprintf("%s", iv2)
				return deepEqual(v1, sv2)
			}
		case int:
			switch v2 := iv2.(type) {
			case string:
				iv2, err := strconv.Atoi(trimString(v2))
				if err == nil && v1 == iv2 {
					return true
				}
			case float64:
				if v1 == int(v2) {
					return true
				}
			case []uint8:
				sv2 := uint8sToString(v2)
				return deepEqual(v1, sv2)
			}
		case float64:
			switch v2 := iv2.(type) {
			case string:
				fv2, err := strconv.ParseFloat(trimString(v2), 64)
				if err == nil && v1 == fv2 {
					return true
				}
			case int:
				if int(v1) == v2 {
					return true
				}
			case []uint8:
				sv2 := uint8sToString(v2)
				return deepEqual(v1, sv2)
			}
		case []uint8:
			sv1 := uint8sToString(v1)
			return deepEqual(sv1, iv2)
		}

		return false
	}

	for _, resRow := range resRows {
		resMap, ok := resRow.(QueryResponseRow)
		if !ok {
			continue
		}

		if reflect.DeepEqual(rowMap, resMap) {
			return nil
		}

		hasAllColumn := true
		for rowKey, rowData := range rowMap {
			resData, ok := resMap[rowKey]
			if !ok {
				hasAllColumn = false
				break
			}

			if !deepEqual(rowData, resData) {
				hasAllColumn = false
				break
			}
		}

		if hasAllColumn {
			return nil
		}
	}

	return fmt.Errorf(errorJSONResponseHasNoRow, res.Data, row)
}

// String returns the string representation.
func (res *QueryResponse) String() string {
	return fmt.Sprintf("%v", res.Data)
}
