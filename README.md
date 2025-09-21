# key-value-store

A simple REST API key-value store with a test client.

## Overview

The `kv_service` provides a REST API for a simple key-value store.
It is served locally at localhost:8080 by default.

The API base url is `/api/v1`

Endpoints:

| Name        | Method | Body               | Success Response Format     | Error Response Format      | Notes                                                          |
| ----------- | ------ | ------------------ | --------------------------- | -------------------------- | -------------------------------------------------------------- |
| /keys/<key> | GET    | N/A                | {"value": <value>}          | {"error": <error message>} | Requesting a key that does not exist returns a value of `null` |
| /keys/<key> | POST   | {"value": <value>} | {"message": "Key set."}     | {"error": <error message>} |                                                                |
| /keys/<key> | DELETE | N/A                | {"message": "Key deleted."} | {"error": <error message>} |                                                                |

## Setup

### Installation

The following dependencies are required:

1. Make
2. Docker
3. Docker compose

### Local development

To run the kv_service in dev mode with hot reloading and Gin debugging logs, cd to the project root and run `make dev-kvs-up`. To stop it, run `make dev-kvs-down`.

## Design notes

1. The kv store and service intentionally limit the 'error' cases by returning nil for keys not yet defined and no-oping if attempting to delete a key that does not exist. This reduces complexity by eliminating the need to check for and handle those errors within the calling code.
2. The kv service's endpoint structure of `/keys/:key` allows for extendibility if we want to have other operations across all keys, such as a `DELETE` request to `/keys` to clear all key/values at once.
