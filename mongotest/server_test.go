// Copyright (C) 2022 The go-mongo Authors All rights reserved.
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

package mongotest

import (
	"testing"

	"github.com/cybergarage/go-logger/log"
)

func TestServer(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	server := NewServer()
	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Tutorial", func(t *testing.T) {
		RunClientTest(t, server)
	})

	t.Run("DBAuth", func(t *testing.T) {
		TestDBAuth(t, server)
	})

	t.Run("YCSB", func(t *testing.T) {
		YCSBTest(t)
	})

	err = server.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTLSServer(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	server := NewServer()

	server.SetTLSEnabled(true)
	server.SetServerKey(TestSeverKey)
	server.SetServerCert(TestServerCert)
	server.SetRootCerts(TestCACert)

	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	err = server.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}
