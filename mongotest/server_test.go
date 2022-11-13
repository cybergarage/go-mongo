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

package mongotest

import (
	"context"
	"fmt"
	"testing"

	"github.com/cybergarage/go-logger/log"
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
	log.SetStdoutDebugEnbled(true)

	server := NewServer()
	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Tutorial", func(t *testing.T) {
		t.Run("CRUD Operations", func(t *testing.T) {
			testTutorialCRUDOperations(t)
		})
	})

	// YCSB

	workloads := []string{"workloada", "workloadb"}
	t.Run("YCSB", func(t *testing.T) {
		for _, workload := range workloads {
			t.Run(workload, func(t *testing.T) {
				ExecYCSBWorkload(t, workload)
			})
		}
	})

	err = server.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}

func testTutorialCRUDOperations(t *testing.T) {
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
		t.Error(err)
		return
	}

	// Use BSON Objects in Go

	collection := client.Database("test").Collection("trainers")

	// Insert documents

	t.Run("Insert documents", func(t *testing.T) {
		trainers := []Trainer{
			ash,
			misty,
			brock,
		}

		for _, trainer := range trainers {
			t.Run(fmt.Sprintf("%s %d %s", trainer.Name, trainer.Age, trainer.City), func(t *testing.T) {
				insertResult, err := collection.InsertOne(context.TODO(), trainer)
				if err != nil {
					t.Error(err)
				}
				t.Log("Inserted a single document: ", insertResult.InsertedID)
			})
		}
	})

	// Find all inserted documents

	t.Run("Find inserted documents", func(t *testing.T) {
		filters := []bson.D{
			{{Key: "name", Value: "Ash"}},
			{{Key: "name", Value: "Misty"}},
			{{Key: "name", Value: "Brock"}},
		}

		trainers := []Trainer{
			ash,
			misty,
			brock,
		}

		for n, filter := range filters {
			t.Run(fmt.Sprintf("%s", filter), func(t *testing.T) {
				var result Trainer

				err = collection.FindOne(context.TODO(), filter).Decode(&result)
				if err != nil {
					t.Error(err)
				}

				t.Logf("Found a single document: %+v\n", result)

				if result != trainers[n] {
					t.Errorf("Found result is not matched : (%v != %v)", result, filter)
					return
				}
			})
		}
	})

	// Update documents

	t.Run("Update documents", func(t *testing.T) {
		filter := bson.D{{Key: "name", Value: "Ash"}}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "age", Value: 11},
			}},
		}

		t.Run(fmt.Sprintf("%s", update), func(t *testing.T) {
			updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				t.Error(err)
			}
			t.Logf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
		})
	})

	// Find all inserted documents

	t.Run("Find updated documents", func(t *testing.T) {
		ash.Age = 11

		filters := []bson.D{
			{{Key: "name", Value: "Brock"}},
			{{Key: "name", Value: "Misty"}},
			{{Key: "name", Value: "Ash"}},
		}

		trainers := []Trainer{
			brock,
			misty,
			ash,
		}

		for n, filter := range filters {
			t.Run(fmt.Sprintf("%s", filter), func(t *testing.T) {
				var result Trainer

				err = collection.FindOne(context.TODO(), filter).Decode(&result)
				if err != nil {
					t.Error(err)
				}

				t.Logf("Found a single document: %+v\n", result)

				if result != trainers[n] {
					t.Errorf("Found result is not matched : (%v != %v)", result, filter)
					return
				}
			})
		}
	})

	// Delete all documents

	t.Run("Delete documents", func(t *testing.T) {
		filters := []bson.D{
			{{Key: "name", Value: "Ash"}},
			{{Key: "name", Value: "Misty"}},
			{{Key: "name", Value: "Brock"}},
		}

		trainers := []Trainer{
			ash,
			misty,
			brock,
		}

		for i, filter := range filters {
			t.Run(fmt.Sprintf("%s", filter), func(t *testing.T) {
				var result Trainer

				deleteResult, err := collection.DeleteMany(context.TODO(), filter)
				if err != nil {
					t.Error(err)
				}

				t.Logf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

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
						t.Error(err)
						return
					}

					t.Logf("Found a single document: %+v\n", result)

					if result != trainers[j] {
						t.Errorf("Found result is not matched : (%v != %v)", result, trainers[j])
						return
					}
				}
			})
		}
	})
}