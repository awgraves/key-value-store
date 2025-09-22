package main

import (
	"github.com/awgraves/key-value-store/kv_service/store"
)

func main() {
	kvStore := store.NewInMemoryStore()
	r := setupRouter(kvStore)

	r.Run(":8080")
}
