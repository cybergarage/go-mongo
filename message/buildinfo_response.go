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

const (
	BuildInfo = "buildinfo"

	version                  = "version"
	DefaultCompatibleVersion = "3.6.4"
)

// NewDefaultBuildInfoResponse returns a default response instance.
func NewDefaultBuildInfoResponse() *Response {
	defaultElements := map[string]interface{}{
		version:           DefaultCompatibleVersion,
		maxBsonObjectSize: int32(DefaultMaxBsonObjectSize),
	}

	res := NewResponseWithElements(defaultElements)
	res.SetStatus(true)

	return res
}

// NewBuildInfoResponseWithConfig returns a response instance with the specified configuration.
func NewBuildInfoResponseWithConfig(config Config) *Response {
	defaultElements := map[string]interface{}{
		version:           config.GetVersion(),
		maxBsonObjectSize: config.GetMaxBsonObjectSize(),
	}

	res := NewResponseWithElements(defaultElements)
	res.SetStatus(true)

	return res
}
