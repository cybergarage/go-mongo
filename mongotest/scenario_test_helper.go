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
	"path"
	"testing"
)

func RunScenarioTestFiles(t *testing.T, testFilenames []string) {
	t.Helper()

	client := NewMongoDBClient()
	for _, testFilename := range testFilenames {
		t.Run(testFilename, func(t *testing.T) {
			ct := NewScenarioTest()
			err := ct.LoadFile(path.Join(SuiteDefaultTestDirectory, testFilename))
			if err != nil {
				t.Error(err)
				return
			}
			ct.SetClient(client)

			err = ct.Run()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
