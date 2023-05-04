// Copyright (C) 2020 The go-mongo Authors. All rights reserved.
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
	ts   time.Time
	span tracer.SpanContext
}

// newConn returns a connection with a default empty connection.
func newConn() *Conn {
	return &Conn{
		Map:  sync.Map{},
		ts:   time.Now(),
		span: nil,
	}
}

// Timestamp returns the creation time of the connection.
func (conn *Conn) Timestamp() time.Time {
	return conn.ts
}

func (conn *Conn) SetSpanContext(span tracer.SpanContext) {
	conn.span = span
}

func (conn *Conn) SpanContext() tracer.SpanContext {
	return conn.span
}
