// Copyright (C) 2020 Satoshi Konno. All rights reserved.
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

import (
	"sync"
	"time"

	"github.com/cybergarage/go-tracing/tracer"
)

// Conn represents a connection of Wire protocol.
type Conn struct {
	sync.Map
	ts time.Time
	tracer.Context
	authrized bool
}

func newConnWith(ctx tracer.Context) *Conn {
	return &Conn{
		Map:       sync.Map{},
		ts:        time.Now(),
		Context:   ctx,
		authrized: false,
	}
}

// Timestamp returns the creation time of the connection.
func (conn *Conn) Timestamp() time.Time {
	return conn.ts
}

// SpanContext returns the tracer span context of the connection.
func (conn *Conn) SpanContext() tracer.Context {
	return conn.Context
}

// SetAuthrized sets the authrized flag to the connection.
func (conn *Conn) SetAuthrized(authrized bool) {
	conn.authrized = authrized
}

// IsAuthrized returns true if the connection is authrized.
func (conn *Conn) IsAuthrized() bool {
	return conn.authrized
}
