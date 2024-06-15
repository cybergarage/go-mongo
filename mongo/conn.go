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
	"crypto/tls"
	"net"
	"sync"
	"time"

	"github.com/cybergarage/go-tracing/tracer"
	"github.com/google/uuid"
)

// Conn represents a connection of Wire protocol.
type Conn struct {
	net.Conn
	isClosed bool
	sync.Map
	ts time.Time
	tracer.Context
	authrized bool
	tlsState  *tls.ConnectionState
	uuid      uuid.UUID
}

func newConnWith(conn net.Conn, tlsState *tls.ConnectionState) *Conn {
	return &Conn{
		Conn:      conn,
		isClosed:  false,
		Map:       sync.Map{},
		ts:        time.Now(),
		Context:   nil,
		authrized: false,
		tlsState:  tlsState,
		uuid:      uuid.New(),
	}
}

// Close closes the connection.
func (conn *Conn) Close() error {
	if conn.isClosed {
		return nil
	}
	if err := conn.Conn.Close(); err != nil {
		return err
	}
	conn.isClosed = true
	return nil
}

// SetSpanContext sets the span context to the connection.
func (conn *Conn) SetSpanContext(span tracer.Context) {
	conn.Context = span
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

// IsTLSConnection return true if the connection is enabled TLS.
func (conn *Conn) IsTLSConnection() bool {
	return conn.tlsState != nil
}

// TLSConnectionState returns the TLS connection state.
func (conn *Conn) TLSConnectionState() (*tls.ConnectionState, bool) {
	return conn.tlsState, conn.tlsState != nil
}

// UUID returns the UUID of the connection.
func (conn *Conn) UUID() uuid.UUID {
	return conn.uuid
}
