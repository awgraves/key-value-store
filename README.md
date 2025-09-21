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
2. The kv service's endpoint structure of `/keys/:key` allows for extendibility if we want to have other operations across all keys, such as a `GET` or `DELETE` request to `/keys` to view all or clear all key value pairs at once, respectively.
