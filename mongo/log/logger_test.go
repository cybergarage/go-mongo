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

package log

import (
	"errors"
	"testing"
)

const (
	testLogMessage         = "hello"
	nullOutputErrorMessage = "Shared Logger is null, but message is output"
	outputErrorMessage     = "Message can't be output"
)

func TestNullLogger(t *testing.T) {
	SetSharedLogger(nil)

	nOutput := Trace(testLogMessage)
	if 0 < nOutput {
		t.Error(errors.New(nullOutputErrorMessage))
	}

	nOutput = Info(testLogMessage)
	if 0 < nOutput {
		t.Error(errors.New(nullOutputErrorMessage))
	}

	nOutput = Error(testLogMessage)
	if 0 < nOutput {
		t.Error(errors.New(nullOutputErrorMessage))
	}

	nOutput = Warn(testLogMessage)
	if 0 < nOutput {
		t.Error(errors.New(nullOutputErrorMessage))
	}

	nOutput = Fatal(testLogMessage)
	if 0 < nOutput {
		t.Error(errors.New(nullOutputErrorMessage))
	}
}

func TestStdoutLogger(t *testing.T) {
	SetSharedLogger(NewStdoutLogger(LevelTrace))
	defer SetSharedLogger(nil)

	nOutput := Trace(testLogMessage)
	if nOutput <= 0 {
		t.Error(errors.New(outputErrorMessage))
	}

	nOutput = Info(testLogMessage)
	if nOutput <= 0 {
		t.Error(errors.New(outputErrorMessage))
	}

	nOutput = Error(testLogMessage)
	if nOutput <= 0 {
		t.Error(errors.New(outputErrorMessage))
	}

	nOutput = Warn(testLogMessage)
	if nOutput <= 0 {
		t.Error(errors.New(outputErrorMessage))
	}

	nOutput = Fatal(testLogMessage)
	if nOutput <= 0 {
		t.Error(errors.New(outputErrorMessage))
	}
}
