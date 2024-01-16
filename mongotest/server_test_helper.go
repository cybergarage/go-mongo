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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testDBURL = "mongodb://localhost:27017"

type Trainer struct {
	Name string
	Age  int
	City string
}

func RunClientTest(t *testing.T, server *Server) {
	t.Helper()

	t.Run("Tutorial", func(t *testing.T) {
		t.Run("test.trainers", func(t *testing.T) {
			TestTutorialCRUDOperations(t)
		})
	})

	t.Run("DBAuth", func(t *testing.T) {
		TestDBAuth(t, server)
	})
}

func TestTutorialCRUDOperations(t *testing.T) {
	t.Helper()

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

	clientOptions := options.Client().ApplyURI(testDBURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		t.Error(err)
		return
	}

	// Use BSON Objects in Go

	collection := client.Database("test").Collection("trainers")

	// Insert documents

	t.Run("InsertOne", func(t *testing.T) {
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

	t.Run("FindOne", func(t *testing.T) {
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

	t.Run("UpdateOne", func(t *testing.T) {
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

	t.Run("FindOne", func(t *testing.T) {
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

	t.Run("DeleteMany", func(t *testing.T) {
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

func TestDBAuth(t *testing.T, server *Server) {
	t.Helper()

	server.SetAuthrization(true)
	defer server.SetAuthrization(false)

	// Authentication Mechanisms — Go Driver
	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/auth/

	credential := options.Credential{ // nolint: exhaustruct
		AuthSource: "admin",
		Username:   "test",
		Password:   "test",
	}

	clientOptions := options.Client().ApplyURI(testDBURL).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		t.Skip(err)
		return
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		t.Skip(err)
		return
	}
}
