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

// Server is an instance for MongoDB protocols.
type Server struct {
	Addr         string
	Port         int
	queryHandler QueryHandler
}

// NewServer returns a new server instance.
func NewServer() *Server {
	server := &Server{
		Addr: "",
		Port: DefaultPort,
	}
	server.SetQueryHandler(server)
	return server
}

// SetPort sets a listen port.
func (server *Server) SetPort(port int) {
	server.Port = port
}

// GetPort returns a listent port.
func (server *Server) GetPort() int {
	return server.Port
}

// Start starts the server server.
func (server *Server) Start() error {
	return nil
}

// Stop stops the server server.
func (server *Server) Stop() error {
	return nil
}

// SetQueryHandler sets a query handler.
func (server *Server) SetQueryHandler(handler QueryHandler) {
	server.queryHandler = handler
}
