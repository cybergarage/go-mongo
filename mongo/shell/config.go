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

package shell

const (
	defaultHost = "localhost"
	defaultPort = 27017
)

// Config stores server configuration parammeters.
type Config struct {
	Host                  string
	Port                  int
	TLSEnabled            bool
	TLSCertificateKeyFile string
	TLSCAFile             string
}

// NewDefaultConfig returns a default configuration instance.
func NewDefaultConfig() *Config {
	config := &Config{
		Host:                  defaultHost,
		Port:                  defaultPort,
		TLSEnabled:            false,
		TLSCertificateKeyFile: "",
		TLSCAFile:             "",
	}
	return config
}

// SetHost sets a host address.
func (config *Config) SetHost(host string) {
	config.Host = host
}

// SetPort sets a listen port.
func (config *Config) SetPort(port int) {
	config.Port = port
}

// SetTLSEnabled sets a TLS enabled flag.
func (config *Config) SetTLSEnabled(enabled bool) {
	config.TLSEnabled = enabled
}

// SetTLSCertificateKeyFile sets a TLS certificate key file.
func (config *Config) SetTLSCertificateKeyFile(file string) {
	config.TLSCertificateKeyFile = file
}

// SetTLSCAFile sets a TLS CA file.
func (config *Config) SetTLSCAFile(file string) {
	config.TLSCAFile = file
}
