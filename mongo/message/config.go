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
	// MaxBsonObjectSize should return a max limitation value of BSON object size.
	MaxBsonObjectSize() int32
	// MaxMessageSizeBytes should return a max limitation value of message size.
	MaxMessageSizeBytes() int32
	// MaxWriteBatchSize should return a max limitation value of write batch size.
	MaxWriteBatchSize() int32
	// LogicalSessionTimeoutMinutes should return a settion timeout value.
	LogicalSessionTimeoutMinutes() int32
	// MinWireVersion should return a min supported version.
	MinWireVersion() int32
	// MaxWireVersion should return a max supported version.
	MaxWireVersion() int32
	// IsReadOnly should return true when the instance does not support write operations.
	IsReadOnly() bool
	// Compressions should return supported compress strings.
	Compressions() []string
}

// BuildConfig represents a limit configurations for 'buildInfo' command.
type BuildConfig interface {
	// Version should return a software version.
	Version() string
}

// Config represents all configurations for MongoDB.
type Config interface {
	ServerConfig
	BuildConfig
}
