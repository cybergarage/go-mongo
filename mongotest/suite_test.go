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

package mongotest

import (
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mongo/mongo/shell"
)

func TestEmbedSuite(t *testing.T) {
	log.EnableStdoutDebug(true)

	server := NewServer()
	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		err := server.Stop()
		if err != nil {
			t.Error(err)
			return
		}
	}()

	client := shell.NewClient()
	err = client.Open()
	if err != nil {
		t.Skip(err.Error())
		return
	}

	defer func() {
		err := client.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	RunEmbedSuite(t, client)
}

func TestTLSEmbedSuite(t *testing.T) {
	log.EnableStdoutDebug(true)

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

	defer func() {
		err := server.Stop()
		if err != nil {
			t.Error(err)
			return
		}
	}()

	client := shell.NewClient()
	client.SetTLSEnabled(true)
	client.SetTLSCertificateKeyFile(TestClientCertFile)
	client.SetTLSCAFile(TestClientCAFile)

	err = client.Open()
	if err != nil {
		t.Skip(err.Error())
		return
	}

	defer func() {
		err := client.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	RunEmbedSuite(t, client)
}
