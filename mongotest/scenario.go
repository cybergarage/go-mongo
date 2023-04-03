// Copyright (C) 2020 The go-mongo Authors. All rights reserved.
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

package mongotest

import (
	"fmt"
	"os"
	"strings"
)

// Line represents a line of a scenario test file.
type Line = string

// Scenario represents a scenario.
type Scenario struct {
	Filename string
	Queries  []string
	Results  []*QueryResponse
}

// NewScenario return a scenario instance.
func NewScenario() *Scenario {
	file := &Scenario{
		Filename: "",
		Queries:  []string{},
		Results:  []*QueryResponse{},
	}
	return file
}

// NewScenarioWithFile return a scenario instance for the specified test scenario file.
func NewScenarioWithFile(filename string) (*Scenario, error) {
	scn := NewScenario()
	err := scn.LoadFile(filename)
	return scn, err
}

// NewScenarioWithBytes return a scenario instance for the specified test scenario bytes.
func NewScenarioWithBytes(name string, b []byte) (*Scenario, error) {
	scn := NewScenario()
	err := scn.ParseBytes(name, b)
	return scn, err
}

// Name returns the loaded scenario file name.
func (scn *Scenario) Name() string {
	return scn.Filename
}

// IsValid checks whether the loaded scenario is available.
func (scn *Scenario) IsValid() error {
	if len(scn.Queries) != len(scn.Results) {
		return fmt.Errorf(errorInvalidScenarioCases, len(scn.Queries), len(scn.Results))
	}
	return nil
}

// LoadFile loads the specified scenario.
func (scn *Scenario) LoadFile(filename string) error {
	lines, err := scn.loadFileLines(filename)
	if err != nil {
		return err
	}

	err = scn.ParseLineStrings(lines)
	if err != nil {
		return fmt.Errorf("%s : %w", filename, err)
	}

	scn.Filename = filename

	return nil
}

func (scn *Scenario) loadFileLines(filename string) ([]Line, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return scn.parseByteLines(fileBytes)
}

// ParseBytes parses the specified scenario bytes.
func (scn *Scenario) ParseBytes(name string, b []byte) error {
	lines, err := scn.parseByteLines(b)
	if err != nil {
		return err
	}

	err = scn.ParseLineStrings(lines)
	if err != nil {
		return fmt.Errorf("%s : %w", name, err)
	}

	scn.Filename = name
	return nil
}

func (scn *Scenario) parseByteLines(fileBytes []byte) ([]Line, error) {
	lines := make([]Line, 0)
	for _, line := range strings.Split(string(fileBytes), "\n") {
		// Skip blank or comment lines
		if len(line) == 0 || strings.HasPrefix(line, "-") {
			continue
		}
		lines = append(lines, line)
	}
	return lines, nil
}

// ParseLineStrings parses the specified scenario line strings.
func (scn *Scenario) ParseLineStrings(lines []string) error {
	var queryStr, resultStr string
	scn.Queries = make([]string, 0)
	scn.Results = make([]*QueryResponse, 0)

	appendQuery := func() {
		if len(queryStr) == 0 {
			return
		}
		scn.Queries = append(scn.Queries, strings.TrimSpace(queryStr))
		queryStr = ""
	}
	appendResult := func() error {
		if len(resultStr) == 0 {
			return nil
		}
		result, err := NewQueryResponseWithString(strings.TrimSpace(resultStr))
		if err != nil {
			return err
		}
		scn.Results = append(scn.Results, result)
		resultStr = ""
		return nil
	}

	inJSON := false
	for _, line := range lines {
		if inJSON {
			if strings.HasPrefix(line, "}") {
				resultStr += line
				err := appendResult()
				if err != nil {
					return err
				}
				inJSON = false
				continue
			}
			resultStr += " " + strings.TrimSpace(line)
		} else {
			if strings.HasPrefix(line, "{") {
				appendQuery()
				resultStr = line
				inJSON = true
				continue
			}
			queryStr += " " + strings.TrimSpace(line)
		}
	}

	return scn.IsValid()
}

// String returns the string representation.
func (scn *Scenario) String() string {
	var str string
	nResults := len(scn.Results)
	for n, query := range scn.Queries {
		str += query + "\n"
		if n < nResults {
			str += scn.Results[n].String() + "\n"
		}
	}
	return str
}
