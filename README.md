# key-value-store

A simple REST API key-value store with a test client.

## Overview

### KV Service

The `kv_service` provides a REST API for a simple key-value store.

The API is exposed by default at: `http://localhost:8080/api/v1`.

| Endpoint   | Method | Description      | Request Body     | Success Response Format | Error Response Format | Notes                                                 |
| ---------- | ------ | ---------------- | ---------------- | ----------------------- | --------------------- | ----------------------------------------------------- |
| /keys/:key | GET    | Retrieve a value | N/A              | {"value": value}        | {"error": msg}        | Returns a `null` value response for keys not found    |
| /keys/:key | POST   | Set a value      | {"value": value} | {"message": msg}        | {"error": msg}        |                                                       |
| /keys/:key | DELETE | Delete a key     | N/A              | {"message": msg}        | {"error": msg}        | Returns a success response even for non-existent keys |

### Test Client

The `test_client` is a separate service that provides its own REST API that connects to the `kv_service` and verifies its functionality.

The test client API is exposed by default at: `http://localhost:8081/api/v1`.

GET requests to the following endpoints will return a 200 response with a success message if the test scenario passes, otherwise they will return a status code 500 with additional error details.
| Endpoint | Description | Success Response Format | Error Response Format |
| --------------- | -------------------------------------------------------------- | -------------------------------- | --------------------------------- |
| /test_deletion | Verifies a key can be set and then deleted | {"message": msg} | {"message": msg, "error": err} |
| /test_overwrite | Verifies a key can be set and then overwritten with a new value | {"message": msg} | {"message": msg, "error": err} |
| /config | Returns the current service config values | {"kv_api_v1_base_url": value} | N/A |

## Setup

### Installation

The following dependencies are required:

1. [Make](https://www.gnu.org/software/make/)
2. [Docker](https://www.docker.com/)
3. [Docker compose](https://docs.docker.com/compose/install/)

### Make

For development ease, this project uses `make`. For a list of available targets, run the `make` command without any args.

### Local development

For local development, use the make commands targeting the development Docker configuration:

- `make dev-up` - Start both services in development mode
- `make dev-down` - Stop both services in development mode
- `make dev-kvs-up` - Start only the key-value service in development mode
- `make dev-kvs-down` - Stop only the key-value service in development mode
- `make dev-logs` - View logs from development services

The development setup includes:

- Hot reloading for code changes
- Gin debug logging

### Production deployment

For production deployment, use make commands targeting the production Docker configuration:

- `make prod-up` - Start both services in production mode (optimized builds, no hot reloading)
- `make prod-down` - Stop both services in production mode
- `make prod-build` - Build production images without starting services
- `make prod-logs` - View logs from production services

The production setup includes:

- Multi-stage Docker builds for minimal image sizes
- Optimized Go binaries without development tools
- Release mode for Gin (no debug logging)

## Testing

### Unit tests

To run all go unit tests across both services, execute `make test`.
To run unit tests for a specific service, execute `make test-kvs` or `make test-client`.

To update the test image (after adding/removing dependencies, for example), execute: `make build-test-image`

### Test client -> KV service

To test the `kv_service` with the `test_client`:

1. Execute either `make dev-up` or `make prod-up` from the project root.
2. Visit `http://localhost:8081/api/v1/:endpoint_name` (see endpoints in the above table in the overview section.)
3. You should receive success responses for passing tests, otherwise you should receive error messages to pinpoint where something may have failed.

## Design notes

1. The kv store implementation and service intentionally limit the 'error' cases by returning nil for keys not yet defined and no-oping if attempting to delete a key that does not exist. This reduces complexity by eliminating the need to check for and handle those errors within the calling code.
2. The kv service's endpoint structure of `/keys/:key` allows for extendibility if we want to have other operations across all keys, such as a `GET` or `DELETE` request to `/keys` to view all or clear all key value pairs at once, respectively.
3. Both services define their handler logic within their respective `router.go` files. At this stage I think its simpler to keep these together, though would certainly split those out into separate `handlers` modules if the number of endpoints grew.
4. `store` and `client` are separate modules with their own generic interfaces and currently 1 implmentation each. These offer flexibility to write other implementations in the future within these modules (ie a distributed KV store, or a gRPC client).
