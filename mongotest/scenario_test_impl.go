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

	"github.com/cybergarage/go-logger/log"
)

const (
	ScenarioTestFileExt = "qst"
)

// ScenarioTest represents a scenario test.
type ScenarioTest struct {
	Scenario *Scenario
	client   Client
}

// NewScenarioTest returns a scenario test instance.
func NewScenarioTest() *ScenarioTest {
	tst := &ScenarioTest{
		Scenario: NewScenario(),
		client:   nil,
	}
	return tst
}

// NewScenarioTestWithFile return a scenario test instance for the specified test scenario file.
func NewScenarioTestWithFile(filename string) (*ScenarioTest, error) {
	tst := NewScenarioTest()
	err := tst.LoadFile(filename)
	return tst, err
}

// NewScenarioTestWithBytes return a scenario test instance for the specified test scenario bytes.
func NewScenarioTestWithBytes(name string, b []byte) (*ScenarioTest, error) {
	tst := NewScenarioTest()
	err := tst.ParseBytes(name, b)
	return tst, err
}

// SetClient sets a client for testing.
func (tst *ScenarioTest) SetClient(c Client) {
	tst.client = c
}

// Name returns the loaded senario name.
func (tst *ScenarioTest) Name() string {
	return tst.Scenario.Name()
}

// LoadFile loads a specified scenario test file.
func (tst *ScenarioTest) LoadFile(filename string) error {
	tst.Scenario = NewScenario()
	return tst.Scenario.LoadFile(filename)
}

// ParseBytes loads a specified scenario test bytes.
func (tst *ScenarioTest) ParseBytes(name string, b []byte) error {
	tst.Scenario = NewScenario()
	return tst.Scenario.ParseBytes(name, b)
}

// LoadFileWithBasename loads a scenario test file which has specified basename.
func (tst *ScenarioTest) LoadFileWithBasename(basename string) error {
	return tst.LoadFile(basename + "." + ScenarioTestFileExt)
}

// Run runs a loaded scenario test.
func (tst *ScenarioTest) Run() error {
	scenario := tst.Scenario
	if scenario == nil {
		return nil
	}

	err := scenario.IsValid()
	if err != nil {
		return err
	}

	client := tst.client
	if client == nil {
		return fmt.Errorf(errorClientNotFound)
	}

	err = client.Open()
	if err != nil {
		return err
	}
	errTraceMsg := func(n int) string {
		errTraceMsg := tst.Name() + "\n"
		for i := 0; i < n; i++ {
			errTraceMsg += fmt.Sprintf(goodQueryPrefix, i, scenario.Queries[i])
			errTraceMsg += "\n"
		}
		return errTraceMsg
	}
	for n, query := range scenario.Queries {
		log.Infof("[%d] %s", n, query)
		_, err := client.Query(query)
		if err != nil {
			return fmt.Errorf("%s%w", errTraceMsg(n), err)
		}
		/*
			defer rs.Close()

			columns, err := rs.Columns()
			if err != nil {
				return err
			}
			columnCnt := len(columns)

			columnTypes, err := rs.ColumnTypes()
			if err != nil {
				return err
			}

			// NOTE: Run() supports only the following standard column types yet.
			values := make([]interface{}, columnCnt)
			for n, columnType := range columnTypes {
				switch columnType.DatabaseTypeName() {
				case "INTEGER", "INT", "SMALLINT", "TINYINT", "MEDIUMINT", "BIGINT":
					var v int
					values[n] = &v
				case "FLOAT", "DOUBLE":
					var v float64
					values[n] = &v
				case "TEXT", "NVARCHAR":
					var v string
					values[n] = &v
				case "VARBINARY", "BINARY":
					var v string
					values[n] = &v
				default:
					var v interface{}
					values[n] = &v
				}
			}

			rsRows := make([]interface{}, 0)
			for rs.Next() {
				err = rs.Scan(values...)
				if err != nil {
					return err
				}

				row := map[string]interface{}{}
				for i := 0; i < columnCnt; i++ {
					switch v := values[i].(type) {
					case *int:
						row[columns[i]] = *v
					case *float64:
						row[columns[i]] = *v
					case *string:
						row[columns[i]] = *v
					case *interface{}:
						row[columns[i]] = *v
					default:
						row[columns[i]] = values[i]
					}
				}

				rsRows = append(rsRows, row)
			}

			expectedRes := scenario.Results[n]
			expectedRows, err := expectedRes.Rows()
			if err != nil {
				if len(rsRows) != 0 {
					return fmt.Errorf("%s"+errorJSONResponseHasUnexpectedRows, errTraceMsg(n), n, query, rsRows)
				}
			} else {
				if len(rsRows) != len(expectedRows) {
					return fmt.Errorf("%s"+errorJSONResponseUnmatchedRowCount, errTraceMsg(n), n, query, rsRows, expectedRows)
				}
			}

			for _, row := range rsRows {
				err = expectedRes.HasRow(row)
				if err != nil {
					return fmt.Errorf("%s"+errorQueryPrefix+"%w", errTraceMsg(n), n, query, err)
				}
			}
		*/
	}

	err = client.Close()
	if err != nil {
		return err
	}

	return nil
}
