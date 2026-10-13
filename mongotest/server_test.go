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
	"context"
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mongo/mongo/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestServer(t *testing.T) {
	log.EnableStdoutDebug(true)

	server := NewServer()
	startTestServer(t, server)

	// Connect to MongoDB using the Go Driver

	clientOptions := options.Client().ApplyURI(testServerURL(server))
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		t.Error(err)
		return
	}

	cleanupClient(t, client)
	t.Run("Tutorial", func(t *testing.T) {
		RunClientTest(t, client)
	})

	t.Run("DBAuth", func(t *testing.T) {
		TestDBAuth(t, server)
	})

	t.Run("YCSB", func(t *testing.T) {
		YCSBTest(t, testServerURL(server))
	})
}

func TestTLSServer(t *testing.T) {
	log.EnableStdoutDebug(true)

	server := NewServer()

	server.SetTLSEnabled(true)
	server.SetServerKey(TestSeverKey)
	server.SetServerCert(TestServerCert)
	server.SetRootCerts(TestCACert)

	ca, err := auth.NewCertificateAuthenticator(
		auth.WithCommonNameRegexp("localhost"),
	)
	if err != nil {
		t.Error(err)
		return
	}
	server.SetCertificateAuthenticator(ca)

	startTestServer(t, server)

	// Connect to MongoDB using the Go Driver

	tlsConfig, err := server.TLSConfig()
	if err != nil {
		t.Error(err)
		return
	}

	clientOptions := options.Client().ApplyURI(testServerURL(server)).SetTLSConfig(tlsConfig)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		t.Error(err)
		return
	}

	cleanupClient(t, client)
	t.Run("Tutorial", func(t *testing.T) {
		RunClientTest(t, client)
	})
}

func TestSASLServer(t *testing.T) {
	// Authentication - MongoDB Manual v7.0
	// https://www.mongodb.com/docs/manual/core/authentication/

	log.EnableStdoutDebug(true)

	server := NewServer()

	// server.SetTLSEnabled(true)
	// server.SetServerKey(TestSeverKey)
	// server.SetServerCert(TestServerCert)
	// server.SetRootCerts(TestCACert)

	startTestServer(t, server)

	// Connect to MongoDB using the Go Driver

	clientOptions := options.Client().ApplyURI(testServerURL(server))
	// clientOptions := options.Client().ApplyURI(testServerURL(server))
	// tlsConfig, err := server.TLSConfig()
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// clientOptions.SetTLSConfig(tlsConfig)
	cred := options.Credential{
		Username: TestUsername,
		Password: TestPassword,
	}
	clientOptions.SetAuth(cred)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		t.Error(err)
		return
	}

	cleanupClient(t, client)
	t.Run("Tutorial", func(t *testing.T) {
		RunClientTest(t, client)
	})
}
