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
	securityAuthorization        bool
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
		isMaster:                     true,
		compressions:                 nil,
		securityAuthorization:        false,
	}
	return config
}

// IsMaster should return true when the instance is running as master, otherwise false.
func (config *Config) IsMaster() bool {
	return config.isMaster
}

// MaxBsonObjectSize should return a max limitation value of BSON object size.
func (config *Config) MaxBsonObjectSize() int32 {
	return config.maxBsonObjectSize
}

// MaxMessageSizeBytes should return a max limitation value of message size.
func (config *Config) MaxMessageSizeBytes() int32 {
	return config.maxMessageSizeBytes
}

// MaxWriteBatchSize should return a max limitation value of write batch size.
func (config *Config) MaxWriteBatchSize() int32 {
	return config.maxWriteBatchSize
}

// LogicalSessionTimeoutMinutes should return a settion timeout value.
func (config *Config) LogicalSessionTimeoutMinutes() int32 {
	return config.logicalSessionTimeoutMinutes
}

// MinWireVersion should return a min supported version.
func (config *Config) MinWireVersion() int32 {
	return config.minWireVersion
}

// MaxWireVersion should return a max supported version.
func (config *Config) MaxWireVersion() int32 {
	return config.maxWireVersion
}

// IsReadOnly should return true when the instance does not support write operations.
func (config *Config) IsReadOnly() bool {
	return config.readOnly
}

// Compressions should return supported compress strings.
func (config *Config) Compressions() []string {
	return config.compressions
}

// Version should return supported MongoDB version string.
func (config *Config) Version() string {
	return config.version
}

// SetAuthrization sets the authrized flag to the connection.
func (config *Config) SetAuthrization(authorized bool) {
	config.securityAuthorization = authorized
}

// IsAuthrized returns true if the connection is authrized.
func (config *Config) IsAuthrized() bool {
	return config.securityAuthorization
}
