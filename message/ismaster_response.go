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

import (
	"time"

	"github.com/cybergarage/go-mongo/bson"
)

const (
	////////////////////////////////////////
	//Replication Commands
	////////////////////////////////////////

	// isMaster (All Instances)

	IsMaster                     = "ismaster"
	MaxBsonObjectSize            = "maxBsonObjectSize"
	MaxMessageSizeBytes          = "maxMessageSizeBytes"
	MaxWriteBatchSize            = "maxWriteBatchSize"
	LocalTime                    = "localTime"
	LogicalSessionTimeoutMinutes = "logicalSessionTimeoutMinutes"
	MinWireVersion               = "minWireVersion"
	MaxWireVersion               = "maxWireVersion"
	ReadOnly                     = "readOnly"
	Compression                  = "compression"
	SASLSupportedMechs           = "saslSupportedMechs"

	DefaultMaxBsonObjectSize            = 16 * 1024 * 1024
	DefaultMaxMessageSizeBytes          = 48000000
	DefaultMaxWriteBatchSize            = 100000
	DefaultLogicalSessionTimeoutMinutes = 30
	DefaultMinWireVersion               = 0
	DefaultMaxWireVersion               = 7
)

// NewDefaultIsMasterResponse returns a default response instance.
func NewDefaultIsMasterResponse() *Response {
	defaultElements := map[string]interface{}{
		IsMaster:                     true,
		MaxBsonObjectSize:            int32(DefaultMaxBsonObjectSize),
		MaxMessageSizeBytes:          int32(DefaultMaxMessageSizeBytes),
		MaxWriteBatchSize:            int32(DefaultMaxWriteBatchSize),
		LocalTime:                    bson.Datetime(time.Now().Unix()),
		LogicalSessionTimeoutMinutes: int32(DefaultLogicalSessionTimeoutMinutes),
		MinWireVersion:               int32(DefaultMinWireVersion),
		MaxWireVersion:               int32(DefaultMaxWireVersion),
		ReadOnly:                     false,
	}

	res := NewResponseWithElements(defaultElements)
	res.SetStatus(true)

	return res
}

// NewIsMasterResponseWithConfig returns a response instance with the specified configuration.
func NewIsMasterResponseWithConfig(config ServerConfig) *Response {
	defaultElements := map[string]interface{}{
		IsMaster:                     config.IsMaster(),
		MaxBsonObjectSize:            config.GetMaxBsonObjectSize(),
		MaxMessageSizeBytes:          config.GetMaxMessageSizeBytes(),
		MaxWriteBatchSize:            config.GetMaxWriteBatchSize(),
		LocalTime:                    bson.Datetime(time.Now().Unix()),
		LogicalSessionTimeoutMinutes: config.GetLogicalSessionTimeoutMinutes(),
		MinWireVersion:               config.GetMinWireVersion(),
		MaxWireVersion:               config.GetMaxWireVersion(),
		ReadOnly:                     config.GetReadOnly(),
	}

	res := NewResponseWithElements(defaultElements)
	res.SetStatus(true)

	return res
}
