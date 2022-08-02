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

import "github.com/cybergarage/go-mongo/mongo/message"

// Config stores server configuration parammeters.
type Config struct {
	isMaster                     bool
	maxBsonObjectSize            int32
	maxMessageSizeBytes          int32
	maxWriteBatchSize            int32
	logicalSessionTimeoutMinutes int32
	minWireVersion               int32
	maxWireVersion               int32
	readOnly                     bool
	compressions                 []string
	version                      string
}

// NewDefaultConfig returns a default configuration instance.
func NewDefaultConfig() *Config {
	config := &Config{
		maxBsonObjectSize:            message.DefaultMaxBsonObjectSize,
		maxMessageSizeBytes:          message.DefaultMaxMessageSizeBytes,
		maxWriteBatchSize:            message.DefaultMaxWriteBatchSize,
		logicalSessionTimeoutMinutes: message.DefaultLogicalSessionTimeoutMinutes,
		minWireVersion:               message.DefaultMinWireVersion,
		maxWireVersion:               message.DefaultMaxWireVersion,
		readOnly:                     false,
		version:                      message.DefaultCompatibleVersion,
	}
	return config
}

// IsMaster should return true when the instance is running as master, otherwise false.
func (config *Config) IsMaster() bool {
	return config.isMaster
}

// GetMaxBsonObjectSize should return a max limitation value of BSON object size.
func (config *Config) GetMaxBsonObjectSize() int32 {
	return config.maxBsonObjectSize
}

// GetMaxMessageSizeBytes should return a max limitation value of message size.
func (config *Config) GetMaxMessageSizeBytes() int32 {
	return config.maxMessageSizeBytes
}

// GetMaxWriteBatchSize should return a max limitation value of write batch size.
func (config *Config) GetMaxWriteBatchSize() int32 {
	return config.maxWriteBatchSize
}

// GetLogicalSessionTimeoutMinutes should return a settion timeout value.
func (config *Config) GetLogicalSessionTimeoutMinutes() int32 {
	return config.logicalSessionTimeoutMinutes
}

// GetMinWireVersion should return a min supported version.
func (config *Config) GetMinWireVersion() int32 {
	return config.minWireVersion
}

// GetMaxWireVersion should return a max supported version.
func (config *Config) GetMaxWireVersion() int32 {
	return config.maxWireVersion
}

// GetReadOnly should return true when the instance does not support write operations.
func (config *Config) GetReadOnly() bool {
	return config.readOnly
}

// GetCompressions should return supported compress strings.
func (config *Config) GetCompressions() []string {
	return config.compressions
}

// GetVersion should return supported MongoDB version string.
func (config *Config) GetVersion() string {
	return config.version
}
