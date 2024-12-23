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

import (
	"github.com/cybergarage/go-mongo/mongo/auth"
	"github.com/cybergarage/go-mongo/mongo/protocol"
	"github.com/cybergarage/go-tracing/tracer"
)

// MessageListener represents a listener for MongoDB messages.
type MessageListener interface {
	protocol.MessageListener
}

// Server represents a server interface.
type Server interface {
	Config
	auth.Manager

	// Config returns a server configuration.
	Config() Config
	// SetTracer sets a tracer.
	SetTracer(tracer.Tracer)

	// SetMessageListener sets a message listener.
	SetMessageListener(l MessageListener)
	// SetMessageHandler sets a message handler.
	SetMessageHandler(h OpMessageHandler)
	// SetUserCommandExecutor sets a command exector for database operation commands.
	SetUserCommandExecutor(fn UserCommandExecutor)
	// SetDatabaseCommandExecutor sets a command exector for database operation commands.
	SetDatabaseCommandExecutor(fn DatabaseCommandExecutor)
	// SetAuthCommandExecutor  sets a command exector for auth operation commands.
	SetAuthCommandExecutor(fn AuthCommandExecutor)

	// Start starts a server.
	Start() error
	// Stop stops a server.
	Stop() error
	// Restart restarts a server.
	Restart() error
}
