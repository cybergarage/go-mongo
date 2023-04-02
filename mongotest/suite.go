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
	"regexp"

	"github.com/cybergarage/go-sqltest/sqltest"
	"github.com/cybergarage/go-sqltest/sqltest/util"
)

type MongoTest = sqltest.SQLTest

const (
	SuiteDefaultTestDirectory = "./test"
	MongoTestFileExt          = "qst"
)

// Suite represents a test suite for MongoDB.
type Suite struct {
	Tests []*MongoTest
}

// NewSuite returns a new test suite.
func NewSuite() *Suite {
	suite := &Suite{
		Tests: make([]*MongoTest, 0),
	}
	return suite
}

func NewSuiteWithDirectory(dir string) (*Suite, error) {
	suite := NewSuite()
	err := suite.LoadDirectory(dir)
	return suite, err
}

func (suite *Suite) LoadDirectory(dir string) error {
	findPath := util.NewFileWithPath(dir)

	re := regexp.MustCompile(".*\\." + MongoTestFileExt)
	files, err := findPath.ListFilesWithRegexp(re)
	if err != nil {
		return err
	}

	suite.Tests = make([]*MongoTest, 0)
	for _, file := range files {
		s, err := sqltest.NewSQLTestWithFile(file.Path)
		if err != nil {
			return err
		}
		suite.Tests = append(suite.Tests, s)
	}

	return nil
}
