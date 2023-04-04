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

package shell

import (
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
)

func ResponseToJSONString(res string) string {
	reps := []struct {
		from string
		to   string
	}{
		{from: "(.*):", to: "\"$1\":"},
		{from: "ObjectId\\((.*)\\)", to: "$1"},
		{from: "\\'(.*)\\'", to: "\"$1\""},
	}
	jsonStr := res
	for _, rep := range reps {
		re := regexp.MustCompile(rep.from)
		jsonStr = re.ReplaceAllString(jsonStr, rep.to)
	}
	return jsonStr
}

func UnmarshalResponse(res string, to any) (any, error) {
	// Extended JSON
	// https://github.com/mongodb/specifications/blob/master/source/extended-json.rst
	// bson package - go.mongodb.org/mongo-driver/bson - Go Packages
	jsonStr := ResponseToJSONString(res)
	err := bson.UnmarshalExtJSON([]byte(jsonStr), true, &to)
	if err != nil {
		return nil, err
	}
	return to, nil
}

func DecodeResponse(res string) (any, error) {
	var v any
	return UnmarshalResponse(res, v)
}
