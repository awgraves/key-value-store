# key-value-store

A simple REST API key-value store with a test client.

## Overview

The `kv_service` provides a REST API for a simple key-value store.

## Setup

### Installation

The following dependencies are required:

1. Make
2. Docker
3. Docker compose

### Local development

To run the kv_service in dev mode with hot reloading and Gin debugging logs, cd to the project root and run `make dev-kvs-up`. To stop it, run `make dev-kvs-down`.
