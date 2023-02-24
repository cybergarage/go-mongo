# Examples

## go-mongod

The go-mongod is a simple MongoDB-compatible implementation using the go-mongo. The sample implementation is a in-memory MongoDB-compatible server.
```
	go-mongod is an example of implementing a compatible MongoDB server using go-mongo.
	NAME
		go-mongod

	SYNOPSIS
		go-mongod [OPTIONS]

	OPTIONS
	-v      : Enable verbose output.
	-p      : Enable profiling.

	RETURN VALUE
		Return EXIT_SUCCESS or EXIT_FAILURE
```

To install the binary, use the following command. The install command installs the utility programs into `GO_PATH/bin`.

```
make install
```

The profile option enables pprof serves of Go which has the HTTP interface to observe go-mongod profile data.

- [The Go Programming Language - Package pprof](https://golang.org/pkg/net/http/pprof/)
- [The Go Blog - Profiling Go Programs](https://blog.golang.org/profiling-go-programs)
