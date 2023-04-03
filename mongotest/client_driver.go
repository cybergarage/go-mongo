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

type MongoDBClient struct {
	*Config
}

// NewMongoDBClient returns a client instance.
func NewMongoDBClient() Client {
	client := &MongoDBClient{
		Config: NewDefaultConfig(),
	}
	return client
}

// Open opens a database specified by the internal configuration.
func (client *MongoDBClient) Open() error {
	return nil
}

// Close closes opens a database specified by the internal configuration.
func (client *MongoDBClient) Close() error {
	return nil
}

// Query executes a query that returns rows.
func (client *MongoDBClient) Query(query string, args ...interface{}) (any, error) {
	return nil, nil
}
