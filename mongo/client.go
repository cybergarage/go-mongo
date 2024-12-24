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
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is an instance for Graphite protocols.
type Client struct {
	host       string
	port       int
	database   string
	collection string
	timeout    time.Duration
	conn       *mongo.Client
}

// NewClient returns a client instance.
func NewClient() *Client {
	client := &Client{
		host:       DefaultHost,
		port:       DefaultPort,
		database:   "",
		collection: "",
		timeout:    (time.Second * DefaultTimeoutSecond),
		conn:       nil,
	}

	return client
}

// SetHost sets a destination host.
func (client *Client) SetHost(host string) {
	client.host = host
}

// Host returns a destination host.
func (client *Client) Host() string {
	return client.host
}

// SetPort sets a destination port.
func (client *Client) SetPort(port int) {
	client.port = port
}

// Port returns a destination port.
func (client *Client) Port() int {
	return client.port
}

// SetDatabase sets a destination database.
func (client *Client) SetDatabase(db string) {
	client.database = db
}

// Database returns a destination database.
func (client *Client) Database() string {
	return client.database
}

// SetCollection sets a destination collection.
func (client *Client) SetCollection(col string) {
	client.database = col
}

// Collection returns a destination collection.
func (client *Client) Collection() string {
	return client.collection
}

// SetTimeout sets a timeout for the request.
func (client *Client) SetTimeout(d time.Duration) {
	client.timeout = d
}

// Timeout return  the timeout for the request.
func (client *Client) Timeout() time.Duration {
	return client.timeout
}

// Connect connects the specified destination MongoDB server.
func (client *Client) Connect() error {
	var err error
	uri := fmt.Sprintf("mongodb://%s", net.JoinHostPort(client.host, strconv.Itoa(client.port)))
	clientOptions := options.Client().ApplyURI(uri)
	client.conn, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the current connection.
func (client *Client) Close() error {
	if client.conn == nil {
		return nil
	}
	err := client.conn.Disconnect(context.TODO())
	if err == nil {
		client.conn = nil
	}
	return err
}

// InsertOne posts the specified document to the destination collection.
func (client *Client) InsertOne(doc interface{}) error {
	if client.conn == nil {
		return fmt.Errorf(errorLostConnection, client.host, client.port)
	}

	col := client.conn.Database(client.database).Collection(client.collection)
	if col == nil {
		return fmt.Errorf(errorCollectionNotFound, client.database, client.collection)
	}

	_, err := col.InsertOne(context.TODO(), doc)

	return err
}
