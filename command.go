// Copyright (C) 2019 The go-mongo Authors. All rights reserved.
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

package mongo

import (
	"github.com/cybergarage/go-mongo/message"
)

// Command represents a query command of MongoDB database command.
type Command = message.Command

// Query represents a query of MongoDB database command.
type Query = message.Query

// CommandExecutor represents an interface for MongoDB query commands.
type CommandExecutor interface {
	message.CommandExecutor
	message.QueryCommandExecutor
}
