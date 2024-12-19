// Copyright (C) 2024 The go-mongo Authors. All rights reserved.
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

package scram

import (
	"github.com/cybergarage/go-sasl/sasl/mech"
	"github.com/cybergarage/go-sasl/sasl/scram"
)

// Message represents a SCRAM message.
type Message = scram.Message

// Response represents a SASL mechanism response.
type Response = mech.Response

// NewMessage creates a new SCRAM message.
func NewMessageWithError(err error) *Message {
	return scram.NewMessageWithError(err)
}

// IsStandardError returns true if the specified error is a standard error.
func IsStandardError(err error) bool {
	return scram.IsStandardError(err)
}
