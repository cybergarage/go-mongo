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
	"github.com/cybergarage/go-sqltest/sqltest"
)

type MongoTest = sqltest.SQLTest

const (
	MongoTestSuiteDefaultTestDirectory = "./test"
)

// MongoTestSuite represents a test suite for MongoDB.
type MongoTestSuite struct {
	Tests []*MongoTest
}

// NewMongoTestSuite returns a new test suite.
func NewMongoTestSuite() *MongoTestSuite {
	suite := &MongoTestSuite{
		Tests: make([]*MongoTest, 0),
	}
	return suite
}
