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

	"github.com/cybergarage/go-mongo/mongo/bson"
)

const (
	maxBsonObjectSize            = "maxBsonObjectSize"
	maxMessageSizeBytes          = "maxMessageSizeBytes"
	maxWriteBatchSize            = "maxWriteBatchSize"
	localTime                    = "localTime"
	logicalSessionTimeoutMinutes = "logicalSessionTimeoutMinutes"
	minWireVersion               = "minWireVersion"
	maxWireVersion               = "maxWireVersion"
	readOnly                     = "readOnly"
	compression                  = "compression"

	DefaultMaxBsonObjectSize            = 16 * 1024 * 1024
	DefaultMaxMessageSizeBytes          = 48000000
	DefaultMaxWriteBatchSize            = 100000
	DefaultLogicalSessionTimeoutMinutes = 30
	DefaultMinWireVersion               = 0
	DefaultMaxWireVersion               = 7
)

// NewDefaultIsMasterResponse returns a default response instance.
func NewDefaultIsMasterResponse() (*Response, error) {
	defaultElements := map[string]interface{}{
		IsMaster:                     true,
		maxBsonObjectSize:            int32(DefaultMaxBsonObjectSize),
		maxMessageSizeBytes:          int32(DefaultMaxMessageSizeBytes),
		maxWriteBatchSize:            int32(DefaultMaxWriteBatchSize),
		localTime:                    bson.Datetime(time.Now().Unix()),
		logicalSessionTimeoutMinutes: int32(DefaultLogicalSessionTimeoutMinutes),
		minWireVersion:               int32(DefaultMinWireVersion),
		maxWireVersion:               int32(DefaultMaxWireVersion),
		readOnly:                     false,
	}

	res, err := NewResponseWithElements(defaultElements)
	if err != nil {
		return nil, err
	}
	res.SetStatus(true)

	return res, nil
}

// NewIsMasterResponseWithConfig returns a response instance with the specified configuration.
func NewIsMasterResponseWithConfig(config ServerConfig) (*Response, error) {
	defaultElements := map[string]interface{}{
		IsMaster:                     config.IsMaster(),
		maxBsonObjectSize:            config.MaxBsonObjectSize(),
		maxMessageSizeBytes:          config.MaxMessageSizeBytes(),
		maxWriteBatchSize:            config.MaxWriteBatchSize(),
		localTime:                    bson.Datetime(time.Now().Unix()),
		logicalSessionTimeoutMinutes: config.LogicalSessionTimeoutMinutes(),
		minWireVersion:               config.MinWireVersion(),
		maxWireVersion:               config.MaxWireVersion(),
		readOnly:                     config.IsReadOnly(),
	}

	res, err := NewResponseWithElements(defaultElements)
	if err != nil {
		return nil, err
	}
	res.SetStatus(true)

	return res, nil
}
