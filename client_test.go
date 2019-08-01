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
package mongo

import (
	"testing"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func TestNewClient(t *testing.T) {
	server := NewServer()
	err := server.Start()
	if err != nil {
		t.Skip(err)
		return
	}
	defer server.Stop()

	// See : MongoDB Go Driver Tutorial
	// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial

	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}

	docs := []Trainer{
		ash,
		misty,
		brock,
	}

	client := NewClient()

	err = client.Connect()
	if err != nil {
		t.Skip(err)
		return
	}

	client.SetDatabase("test")
	client.SetCollection("trainers")

	for _, doc := range docs {
		err := client.InsertOne(doc)
		if err != nil {
			t.Skip(err)
			return
		}
	}

	err = client.Close()
	if err != nil {
		t.Error(err)
		return
	}
}
