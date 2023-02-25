# Getting Started

This section describes how to implement your MongoDB-compatible serverusing go-mongo and see [Examples](examples.md) about the sample implementation.

## STEP1: Inheriting Server

The go-mysql offers a core server, [mongo.Server](../mongo/server.go), and so inherit the core server in your instance as the following.

```
import (
	"github.com/cybergarage/go-mongo/mongo"
)

type MyServer struct {
	*mongo.Server
	documents []bson.Document
}

func NewMyServer() *MyServer {
	return &MyServer{
		Server: mongo.NewServer(),
	}
}
```

## STEP2: Preparing Query Handler

To handle queries to the your server, prepare a query handler according to [mongo.UserCommandExecutor](../mongo/command.go) interface.

```
func NewMyServer() *MyServer {
	myserver := &MyServer{
		Server: mysql.NewServer(),
	}
    Myserver.SetQueryExecutor(myserver)
    return myserver
}

func (server *MyServer) Insert(*Query) (int32, error) {
    .....
}

func (server *MyServer) Update(*Query) (int32, error) {
    .....
}

func (server *MyServer) Find(*Query) ([]bson.Document, error) {
    .....
}

func (server *MyServer) Delete(*Query) (int32, error) {
    .....
}
```

Since the go-mongo handles all system commands automatically, developers can easily implement their MongoDB-compatible server only by simply handling the query commands.

## STEP3: Starting Server 

After implementing the query handler, start your server using  [mongo.Server::Start()](../mongo/server.go).

```
server := NewMyServer()

err := server.Start()
if err != nil {
	t.Error(err)
	return
}
defer server.Stop()

.... 
```
