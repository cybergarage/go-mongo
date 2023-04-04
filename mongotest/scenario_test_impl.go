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
	"reflect"

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
		queryRes, err := client.Query(query)
		if err != nil {
			return fmt.Errorf("%s%w", errTraceMsg(n), err)
		}
		expectedRes := scenario.Expecteds[n]
		if !reflect.DeepEqual(queryRes, expectedRes) {
			return fmt.Errorf("%s", errTraceMsg(n))
		}
	}

	return nil
}
