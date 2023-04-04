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

import (
	"fmt"
	"os/exec"
)

const (
	mongosh = "mongosh"
)

type Client struct {
	*Config
}

// NewClient returns a shell client instance.
func NewClient() *Client {
	return &Client{
		Config: NewDefaultConfig(),
	}
}

// Open opens a database specified by the internal configuration.
func (client *Client) Open() error {
	return nil
}

// Close closes opens a database specified by the internal configuration.
func (client *Client) Close() error {
	return nil
}

// Query executes a query that returns rows.
func (client *Client) Query(query string) (any, error) {
	var args []string
	args = append(args, "--eval", fmt.Sprintf("'%s'", query))
	out, err := exec.Command(mongosh, args...).CombinedOutput()
	if err == nil {
		// TODO: Parse the output result set response
		return nil, nil
	}
	return nil, fmt.Errorf("%w : %s", err, string(out))
}
