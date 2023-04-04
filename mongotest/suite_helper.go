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

	"github.com/cybergarage/go-mongo/mongo/shell"
	"github.com/cybergarage/go-mongo/mongotest/test"
)

const ScenarioTestDatabase = "qst"

func RunEmbedSuite(t *testing.T) {
	t.Helper()

	cs, err := NeweEmbedSuite(test.EmbedTests)
	if err != nil {
		t.Error(err)
		return
	}

	client := shell.NewClient()
	err = client.Open()
	if err != nil {
		t.Skipf(err.Error())
		return
	}

	cs.SetClient(client)

	for _, test := range cs.Tests {
		t.Run(test.Name(), func(t *testing.T) {
			test.SetClient(cs.client)
			err := test.Run()
			if err != nil {
				t.Errorf("%s : %s", test.Name(), err.Error())
			}
		})
	}

	err = client.Close()
	if err != nil {
		t.Error(err)
	}
}

func RunLocalSuite(t *testing.T) {
	t.Helper()

	cs, err := NewSuiteWithDirectory(SuiteDefaultTestDirectory)
	if err != nil {
		t.Error(err)
		return
	}

	client := shell.NewClient()

	err = client.Open()
	if err != nil {
		t.Error(err)
		return
	}

	cs.SetClient(client)

	for _, test := range cs.Tests {
		t.Run(test.Name(), func(t *testing.T) {
			test.SetClient(cs.client)
			err := test.Run()
			if err != nil {
				t.Errorf("%s : %s", test.Name(), err.Error())
			}
		})
	}
}
