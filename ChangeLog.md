# Changelog

## v1.3.x (2023-xx-xx)
- Updated executor interfaces for major MongoDB commands
- Added createIndex interface

## v1.2.2 (2024-12-28)
- Supported certificate authentication for TLS connection
- Fixed gosec warnings

## v1.2.1 (2024-09-18)
- Supported SCRAM-SHA-1 and SCRAM-SHA-256 authentication
- Supported helloOk protocol negotiation
- Added a wire protocol utility (wirehexdump)
- Fixed gosec warnings

## v1.2.0 (2024-08-22)
- Updated protocol.Message and message.Query interfaces
- Updated SASL authentication interfaces

## v1.1.4 (2024-06-29)
- Added connection manager

## v1.1.3 (2024-05-25)
- Added auth manager

## v1.1.2 (2024-05-22)
- Updated mongo shell client

## v1.1.1 (2024-05-21)
- Updated embedded MongoDB tests

## v1.1.0 (2024-05-20)
- Supported TLS connection
- Updated go-tracing package

## v1.0.2 (2023-05-04)
- Updated Conn to embed tracing context
- Updated tracing spans

## v1.0.1 (2023-05-04)
- Updated Conn interface
- Updated tracing spans

## v1.0.0 (2023-05-04)
- Fixed executor interfaces for basic MongoDB commands
- Updated logger functions to output more detailed messages
- Added tracing interface

## v0.9.5 (2023-04-23)
- Updated embed test interface

## v0.9.4 (2023-04-08)
- Added scenario test framework
- Added mongosh-based client for testing
- Updated wire protocol parser to support the protocols used by mongosh

## v0.9.3 (2023-04-02)
- Added connection logs
- Added Dockerfile

## v0.9.2 (2023-03-28)
- Added Conn parameter to executor functions
- Updated mongo-driver from v1.10.0 to v1.11.2

## v0.9.1 (2023-02-23)
- Upgraded to go 1.20
- Fixed compiler warnings

## v0.9.0 (2019-08-15)
- Initial public release
