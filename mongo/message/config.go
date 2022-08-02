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

package message

// ServerConfig represents a limit configurations for 'isMaseter' command.
type ServerConfig interface {
	// IsMaster should return true when the instance is running as master, otherwise false
	IsMaster() bool
	// GetMaxBsonObjectSize should return a max limitation value of BSON object size.
	GetMaxBsonObjectSize() int32
	// GetMaxMessageSizeBytes should return a max limitation value of message size.
	GetMaxMessageSizeBytes() int32
	// GetMaxWriteBatchSize should return a max limitation value of write batch size.
	GetMaxWriteBatchSize() int32
	// GetLogicalSessionTimeoutMinutes should return a settion timeout value.
	GetLogicalSessionTimeoutMinutes() int32
	// GetMinWireVersion should return a min supported version.
	GetMinWireVersion() int32
	// GetMaxWireVersion should return a max supported version.
	GetMaxWireVersion() int32
	// GetReadOnly should return true when the instance does not support write operations.
	GetReadOnly() bool
	// GetCompressions should return supported compress strings.
	GetCompressions() []string
}

// BuildConfig represents a limit configurations for 'buildInfo' command.
type BuildConfig interface {
	// GetVersion should return a software version.
	GetVersion() string
}

// Config represents all configurations for MongoDB.
type Config interface {
	ServerConfig
	BuildConfig
}
