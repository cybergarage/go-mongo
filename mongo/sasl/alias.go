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

package sasl

import (
	"github.com/cybergarage/go-sasl/sasl"
	"github.com/cybergarage/go-sasl/sasl/cred"
)

// SASLMechanism represents a SASL mechanism.
type SASLMechanism = sasl.Mechanism

// Option represents a SASL mechanism option.
type SASLOption = sasl.Option

// Group represents a group option.
type SASLGroup = sasl.Group

// AuthzID represents an authorization ID option.
type SASLAuthzID = sasl.AuthzID

// Username represents a username option.
type SASLUsername = sasl.Username

// Password represents a password option.
type SASLPassword = sasl.Password

// Token represents a token.
type SASLToken = sasl.Token

// Email represents an email.
type SASLEmail = sasl.Email

// Payload represents a payload.
type SASLPayload = sasl.Payload

// Authenticators represents a list of credential authenticators.
type SASLAuthenticators = cred.Authenticators

// HashFunc represents a hash function.
type SASLRandomSequence = sasl.RandomSequence

// IterationCount represents an iteration count.
type SASLIterationCount = sasl.IterationCount

// HashFunc represents a hash function.
type SASLHashFunc = sasl.HashFunc

// Challenge represents a challenge.
type SASLChallenge = sasl.Challenge

// Salt represents a salt.
type SASLSalt = sasl.Salt
