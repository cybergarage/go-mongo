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
	Host       string
	Port       int
	Database   string
	Collection string
	Timeout    time.Duration
	conn       *mongo.Client
}

// NewClient returns a client instance.
func NewClient() *Client {
	client := &Client{
		Host:       DefaultHost,
		Port:       DefaultPort,
		Database:   "",
		Collection: "",
		Timeout:    (time.Second * DefaultTimeoutSecond),
		conn:       nil,
	}

	return client
}

// SetHost sets a destination host.
func (client *Client) SetHost(host string) {
	client.Host = host
}

// GetHost returns a destination host.
func (client *Client) GetHost() string {
	return client.Host
}

// SetPort sets a destination port.
func (client *Client) SetPort(port int) {
	client.Port = port
}

// GetPort returns a destination port.
func (client *Client) GetPort() int {
	return client.Port
}

// SetDatabase sets a destination database.
func (client *Client) SetDatabase(db string) {
	client.Database = db
}

// GetDatabase returns a destination database.
func (client *Client) GetDatabase() string {
	return client.Database
}

// SetCollection sets a destination collection.
func (client *Client) SetCollection(col string) {
	client.Database = col
}

// GetCollection returns a destination collection.
func (client *Client) GetCollection() string {
	return client.Collection
}

// SetTimeout sets a timeout for the request.
func (client *Client) SetTimeout(d time.Duration) {
	client.Timeout = d
}

// GetTimeout return  the timeout for the request.
func (client *Client) GetTimeout() time.Duration {
	return client.Timeout
}

// Connect connects the specified destination MongoDB server.
func (client *Client) Connect() error {
	var err error
	uri := fmt.Sprintf("mongodb://%s", net.JoinHostPort(client.Host, strconv.Itoa(client.Port)))
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
		return fmt.Errorf(errorLostConnection, client.Host, client.Port)
	}

	col := client.conn.Database(client.Database).Collection(client.Collection)
	if col == nil {
		return fmt.Errorf(errorCollectionNotFound, client.Database, client.Collection)
	}

	_, err := col.InsertOne(context.TODO(), doc)

	return err
}
