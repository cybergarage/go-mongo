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
	"fmt"
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mongo/mongo/shell"
	"github.com/cybergarage/go-mongo/mongotest/test"
)

func TestEmbedSuite(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	suite, err := NeweEmbedSuite(test.EmbedTests)
	if err != nil {
		t.Error(err)
		return
	}

	client := shell.NewClient()
	suite.SetClient(client)

	server := NewServer()
	err = server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	for _, tst := range suite.Tests {
		t.Run(tst.Name(), func(t *testing.T) {
			tst.SetClient(client)
			err := tst.Run()
			if err != nil {
				t.Error(fmt.Errorf("%s : %w", tst.Name(), err))
			}
		})
	}

	err = server.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}
