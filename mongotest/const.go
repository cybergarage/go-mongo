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
	_ "embed"
)

const (
	TestUsername       = "test"
	TestPassword       = "password"
	TestClientCertFile = "certs/client.pem"
	TestClientCAFile   = "certs/ca.pem"
)

//go:embed certs/key.pem
var TestSeverKey []byte

//go:embed certs/cert.pem
var TestServerCert []byte

//go:embed certs/ca.pem
var TestCACert []byte
