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
	"regexp"

	"github.com/cybergarage/go-sqltest/sqltest/util"
)

const (
	SuiteDefaultTestDirectory = "./test"
)

// Suite represents a scenario test suite.
type Suite struct {
	Tests  []*ScenarioTest
	client Client
}

// NewSuite returns a scenario test suite instance.
func NewSuite() *Suite {
	suite := &Suite{
		Tests:  make([]*ScenarioTest, 0),
		client: nil,
	}
	return suite
}

// SetClient sets a client for testing.
func (suite *Suite) SetClient(c Client) {
	suite.client = c
}

// NewSuiteWithDirectory returns a scenario test suite instance which loads under the specified directory.
func NewSuiteWithDirectory(dir string) (*Suite, error) {
	suite := NewSuite()

	re := regexp.MustCompile(".*\\." + ScenarioTestFileExt)
	findPath := util.NewFileWithPath(dir)
	files, err := findPath.ListFilesWithRegexp(re)
	if err != nil {
		return nil, err
	}

	suite.Tests = make([]*ScenarioTest, 0)
	for _, file := range files {
		s, err := NewScenarioTestWithFile(file.Path)
		if err != nil {
			return nil, err
		}
		suite.Tests = append(suite.Tests, s)
	}

	return suite, nil
}

// NeweEmbedSuite returns a scenario test suite instance which loads under the specified directory.
func NeweEmbedSuite(tests map[string][]byte) (*Suite, error) {
	suite := NewSuite()
	for name, b := range tests {
		s, err := NewScenarioTestWithBytes(name, b)
		if err != nil {
			return nil, err
		}
		suite.Tests = append(suite.Tests, s)
	}
	return suite, nil
}

// Run runs all loaded scenario tests. The method stops the testing when a scenario test is aborted, and the following tests are not run.
func (suite *Suite) Run() error {
	for _, test := range suite.Tests {
		test.SetClient(suite.client)
		err := test.Run()
		if err != nil {
			return fmt.Errorf("%s : %w", test.Name(), err)
		}
	}
	return nil
}
