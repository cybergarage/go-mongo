// Copyright (C) 2020 The go-mongo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongotest

const (
	errorInvalidScenarioCases          = "query cases [%d] are not equal with result cases [%d]"
	errorClientNotFound                = "client for testing is not found"
	errorInvalidJSONResponse           = "JSON response (%v) is invalid"
	errorJSONResponseNotFound          = "JSON response is not found"
	errorJSONResponseRowsNotFound      = "JSON response (%v) is different from the expected response ('%v')"
	errorJSONResponseHasNoRow          = "JSON response (%v) is different from the expected rows (%v)"
	queryLinePrefix                    = "[%d] %s"
	goodQueryPrefix                    = "O " + queryLinePrefix + " : "
	errorQueryPrefix                   = "X " + queryLinePrefix + " : "
	errorJSONResponseUnmatchedRowCount = errorQueryPrefix + "JSON response row count (%v) is unmatched to (%v)"
	errorJSONResponseHasUnexpectedRows = errorQueryPrefix + "JSON response has unexpected rows (%v)"
)
