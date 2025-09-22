package main

import (
	"os"

	"github.com/awgraves/key-value-store/test_client/client"
)

// getKVServiceAPIv1BaseURL returns the base URL for the KV service API v1
// Uses environment variable KV_SERVICE_API_V1_BASE_URL with fallback to localhost
func getKVServiceAPIv1BaseURL() string {
	if url := os.Getenv("KV_SERVICE_API_V1_BASE_URL"); url != "" {
		return url
	}
	return "http://localhost:8080/api/v1"
}

func main() {
	kvAPIv1BaseURL := getKVServiceAPIv1BaseURL()
	apiClient := client.NewHTTPClient(kvAPIv1BaseURL)

	r := setupRouter(apiClient, kvAPIv1BaseURL)
	r.Run(":8081")
}
