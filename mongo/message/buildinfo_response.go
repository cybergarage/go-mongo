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
	"strconv"
	"strings"
)

const (
	version      = "version"
	versionArray = "versionArray"
)

const (
	// DefaultCompatibleVersion = "4.0.12".
	// DefaultCompatibleVersion = "3.6.4".
	DefaultCompatibleVersion = "3.4.22"
)

// NewDefaultBuildInfoResponse returns a default response instance.
func NewDefaultBuildInfoResponse() *Response {
	defaultElements := map[string]interface{}{
		maxBsonObjectSize: int32(DefaultMaxBsonObjectSize),
	}

	res := NewResponseWithElements(defaultElements)
	res.SetVersion(DefaultCompatibleVersion)
	res.SetStatus(true)

	return res
}

// NewBuildInfoResponseWithConfig returns a response instance with the specified configuration.
func NewBuildInfoResponseWithConfig(config Config) *Response {
	res := NewDefaultBuildInfoResponse()
	res.SetVersion(config.GetVersion())
	res.SetStatus(true)

	return res
}

// SetVersion sets a version string and array of the specified version string.
func (res *Response) SetVersion(ver string) {
	// version
	res.SetStringElement(version, ver)
	// versionArray
	verInts := make([]int32, 0)
	verStrs := strings.Split(ver, ".")
	for _, verStr := range verStrs {
		verInt, err := strconv.Atoi(verStr)
		if err != nil {
			continue
		}
		verInts = append(verInts, int32(verInt))
	}
	res.SetInt32ArrayElements(versionArray, verInts)
}
