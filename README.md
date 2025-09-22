# key-value-store

A simple REST API key-value store with a test client.

## Overview

The `kv_service` provides a REST API for a simple key-value store.

The API base url is `/api/v1`

| Endpoint   | Method | Description      | Request Body     | Success Response Format | Error Response Format | Notes                                                 |
| ---------- | ------ | ---------------- | ---------------- | ----------------------- | --------------------- | ----------------------------------------------------- |
| /keys/:key | GET    | Retrieve a value | N/A              | {"value": value}        | {"error": msg}        | Returns a `null` value response for keys not found    |
| /keys/:key | POST   | Set a value      | {"value": value} | {"message": msg}        | {"error": msg}        |                                                       |
| /keys/:key | DELETE | Delete a key     | N/A              | {"message": msg}        | {"error": msg}        | Returns a success response even for non-existent keys |

The `test_client` is another service that connects to the `kv_service` and verifies its functionality.

The API base url is also `/api/v1` (mirrors the kv_service).

GET requests to the following endpoints will return a 200 response if the test scenario passes, otherwise it will return a 500.
| Endpoint | Description | Success Response Format | Error Response Format |
| --------------- | -------------------------------------------------------------- | -------------------------------- | --------------------------------- |
| /test_deletion | Verifies a key can be set and then deleted | {"message": msg} | {"message": msg, "error": err} |
| /test_overwrite | Verifies a key can be set and then overwritten with a new value | {"message": msg} | {"message": msg, "error": err} |
| /config | Returns the current service config values | {"kv_api_v1_base_url": value} | N/A |

## Setup

### Installation

The following dependencies are required:

1. Make
2. Docker
3. Docker compose

### Local development

For development ease, this project uses `make`. For a list of available targets, run the `make` command without any args.

To run both `kv_service` and `test_client` together in "dev mode" (hot reloading + Gin debugging logs), execute `make dev-up` from the project root. They can be stopped with `make dev-down`.

To start/stop only the `kv_service` in dev mode without starting the test_client, run `make dev-kvs-up` / `make dev-kvs-down`.

## Testing

To run all go unit tests across both services, execute `make test`.

To test the `kv_service` with the `test_client`:

1. Execute `make dev-up` from the project root.
2. Visit http://localhost:8081/api/v1/:endpoint_name (see endpoints in the above table in the overview section.)
3. You should receive success responses for passing tests, and error messages to pinpoint where something failed.

## Design notes

1. The kv store and service intentionally limit the 'error' cases by returning nil for keys not yet defined and no-oping if attempting to delete a key that does not exist. This reduces complexity by eliminating the need to check for and handle those errors within the calling code.
2. The kv service's endpoint structure of `/keys/:key` allows for extendibility if we want to have other operations across all keys, such as a `GET` or `DELETE` request to `/keys` to view all or clear all key value pairs at once, respectively.
3. Both services define their handler logic within their respective `router.go` files. At this stage its simpler to keep these together, though would certainly split those out into separate `handlers` modules if the number of endpoints grew.
