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
type config struct {
	Addr                         string
	Port                         int
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
	securityAuthorizationEnabled bool
}

// NewDefaultConfig returns a default configuration instance.
func NewDefaultConfig() Config {
	return newDefaultConfig()
}

// newDefaultConfig returns a default configuration instance.
func newDefaultConfig() *config {
	config := &config{
		Addr:                         "",
		Port:                         DefaultPort,
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
		securityAuthorizationEnabled: false,
	}
	return config
}

// SetAddress sets a listen address.
func (config *config) SetAddress(addr string) {
	config.Addr = addr
}

// GetAddress returns a listen address.
func (config *config) GetAddress() string {
	return config.Addr
}

// SetPort sets a listen port.
func (config *config) SetPort(port int) {
	config.Port = port
}

// GetPort returns a listent port.
func (config *config) GetPort() int {
	return config.Port
}

// IsMaster should return true when the instance is running as master, otherwise false.
func (config *config) IsMaster() bool {
	return config.isMaster
}

// MaxBsonObjectSize should return a max limitation value of BSON object size.
func (config *config) MaxBsonObjectSize() int32 {
	return config.maxBsonObjectSize
}

// MaxMessageSizeBytes should return a max limitation value of message size.
func (config *config) MaxMessageSizeBytes() int32 {
	return config.maxMessageSizeBytes
}

// MaxWriteBatchSize should return a max limitation value of write batch size.
func (config *config) MaxWriteBatchSize() int32 {
	return config.maxWriteBatchSize
}

// LogicalSessionTimeoutMinutes should return a settion timeout value.
func (config *config) LogicalSessionTimeoutMinutes() int32 {
	return config.logicalSessionTimeoutMinutes
}

// MinWireVersion should return a min supported version.
func (config *config) MinWireVersion() int32 {
	return config.minWireVersion
}

// MaxWireVersion should return a max supported version.
func (config *config) MaxWireVersion() int32 {
	return config.maxWireVersion
}

// IsReadOnly should return true when the instance does not support write operations.
func (config *config) IsReadOnly() bool {
	return config.readOnly
}

// Compressions should return supported compress strings.
func (config *config) Compressions() []string {
	return config.compressions
}

// Version should return supported MongoDB version string.
func (config *config) Version() string {
	return config.version
}

// SetAuthrizationEnabled sets the authorization flag.
func (config *config) SetAuthrizationEnabled(authorized bool) {
	config.securityAuthorizationEnabled = authorized
}

// IsAuthrizationEnabled returns true when the authorization is enabled.
func (config *config) IsAuthrizationEnabled() bool {
	return config.securityAuthorizationEnabled
}
