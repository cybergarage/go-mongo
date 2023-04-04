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
	"testing"
)

func TestQueryResponses(t *testing.T) {
	testResponses := []string{
		"{\nacknowledged: true,\ninsertedId: ObjectId(\"6429625454ee3326b240a06d\")\n}",
		"{\n_id: ObjectId(\"6429625454ee3326b240a06d\"),\nname: 'Ash',\nage: 10,\ncity: 'Pallet Town'}",
		"[\n{\n_id: ObjectId(\"642a8fb2267f49b99957ee13\"),\nname:'Ash',\nage: 10,\ncity: 'Pallet Town'\n}\n]\n",
	}

	for _, res := range testResponses {
		_, err := DecodeResponse(res)
		if err != nil {
			t.Error(err)
		}
	}
}
