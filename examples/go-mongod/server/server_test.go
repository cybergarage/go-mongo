// Copyright (C) 2022 The go-mongo Authors All rights reserved.
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

package server

import (
	"context"
	"fmt"
	"log"
	"testing"

	gomongo_log "github.com/cybergarage/go-mongo/mongo/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func TestServer(t *testing.T) {
	gomongo_log.SetStdoutDebugEnbled(true)

	server := NewServer()
	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Stop()

	// MongoDB Go Driver
	// https://github.com/mongodb/mongo-go-driver
	//
	// MongoDB Go Driver Tutorial
	// https://www.mongodb.com/blog/post/mongodb-go-driver-tutorial

	// Test documents

	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}

	// Connect to MongoDB using the Go Driver

	url := "mongodb://localhost:27017"

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Use BSON Objects in Go

	collection := client.Database("test").Collection("trainers")

	// Insert documents

	trainers := []Trainer{
		ash,
		misty,
		brock,
	}

	for _, trainer := range trainers {
		insertResult, err := collection.InsertOne(context.TODO(), trainer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	}

	// Find all inserted documents

	filters := []bson.D{
		{{Key: "name", Value: "Ash"}},
		{{Key: "name", Value: "Misty"}},
		{{Key: "name", Value: "Brock"}},
	}

	trainers = []Trainer{
		ash,
		misty,
		brock,
	}

	for n, filter := range filters {
		var result Trainer

		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Found a single document: %+v\n", result)

		if result != trainers[n] {
			t.Errorf("Found result is not matched : (%v != %v)", result, filter)
			return
		}
	}

	// Update documents

	filter := bson.D{{Key: "name", Value: "Ash"}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "age", Value: 11},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Find all inserted documents

	ash.Age = 11

	filters = []bson.D{
		{{Key: "name", Value: "Brock"}},
		{{Key: "name", Value: "Misty"}},
		{{Key: "name", Value: "Ash"}},
	}

	trainers = []Trainer{
		brock,
		misty,
		ash,
	}

	for n, filter := range filters {
		var result Trainer

		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Found a single document: %+v\n", result)

		if result != trainers[n] {
			t.Errorf("Found result is not matched : (%v != %v)", result, filter)
			return
		}
	}

	// Delete all documents

	filters = []bson.D{
		{{Key: "name", Value: "Ash"}},
		{{Key: "name", Value: "Misty"}},
		{{Key: "name", Value: "Brock"}},
	}

	trainers = []Trainer{
		ash,
		misty,
		brock,
	}

	for i, filter := range filters {
		var result Trainer

		deleteResult, err := collection.DeleteMany(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

		// Check deleted document

		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		if err == nil {
			t.Errorf("Found the delete document: %+v\n", filter)
			return
		}

		// Find other documents

		for j := (i + 1); j < len(filters); j++ {
			var result Trainer
			err = collection.FindOne(context.TODO(), filters[j]).Decode(&result)
			if err != nil {
				log.Fatal(err)
				return
			}

			fmt.Printf("Found a single document: %+v\n", result)

			if result != trainers[j] {
				t.Errorf("Found result is not matched : (%v != %v)", result, trainers[j])
				return
			}

		}
	}

}
