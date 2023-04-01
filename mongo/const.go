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

package mongo

const (
	// PackageName is the package name.
	PackageName = "gp-mongo"
	// DefaultHost is the default host for MongoDB servers.
	DefaultHost string = "localhost"
	// DefaultPort is the default port for mongod and mongos.
	DefaultPort int = 27017
	// DefaultTimeoutSecond is the default request timeout for MongoDB servers.
	DefaultTimeoutSecond = 5
)
