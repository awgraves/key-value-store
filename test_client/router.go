package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// getKVServiceAPIv1BaseURL returns the base URL for the KV service API v1
// Uses environment variable KV_SERVICE_API_V1_BASE_URL with fallback to localhost
func getKVServiceAPIv1BaseURL() string {
	if url := os.Getenv("KV_SERVICE_API_V1_BASE_URL"); url != "" {
		return url
	}
	return "http://localhost:8080/api/v1"
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	kvAPIv1BaseURL := getKVServiceAPIv1BaseURL()

	v1 := r.Group("/api/v1")
	{
		var client APIv1Client = NewAPIv1Client(kvAPIv1BaseURL)

		v1.GET("/test_deletion", func(c *gin.Context) {
			// Example: This would test deleting a key from the KV service
			testKey := "test-key"
			testValue := "test-value"

			// set the test key
			err := client.SetKey(testKey, testValue)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error setting test key",
					"error":   err.Error(),
				})
				return
			}
			// check the key was set
			value, err := client.GetKey(testKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error getting test key after setting",
					"error":   err.Error(),
				})
				return
			}
			if value != testValue {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error verifying test key after setting",
					"error":   fmt.Sprintf("Test key should be '%v'. Got value %v instead.", testValue, value),
				})
				return
			}
			// delete the key
			err = client.DeleteKey(testKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error deleting test key",
					"error":   err.Error(),
				})
				return
			}
			// check the key was deleted
			value, err = client.GetKey(testKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error getting test key after deletion",
					"error":   err.Error(),
				})
				return
			}
			if value != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error verifying test key after deletion",
					"error":   fmt.Sprintf("Test key should be nil after deletion. Got value %v instead.", value),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Test deletion successful",
			})
		})

		v1.GET("/test_overwrite", func(c *gin.Context) {
			// Example: This would test overwriting a key in the KV service
			testKey := "test-key"
			testOriginalValue := "test-value"
			testNewValue := "new-value"

			// set the test key
			err := client.SetKey(testKey, testOriginalValue)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error setting test key",
					"error":   err.Error(),
				})
				return
			}
			// check the key was set
			value, err := client.GetKey(testKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error getting test key after setting",
					"error":   err.Error(),
				})
				return
			}
			if value != testOriginalValue {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error verifying test key after setting",
					"error":   fmt.Sprintf("Test key should be 'test-value'. Got value %v instead.", value),
				})
				return
			}
			// set the test key again
			err = client.SetKey(testKey, testNewValue)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error setting test key again",
					"error":   err.Error(),
				})
				return
			}
			// check the key was overwritten
			value, err = client.GetKey(testKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error getting test key after overwriting",
					"error":   err.Error(),
				})
				return
			}
			if value != "new-value" {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error verifying test key after overwriting",
					"error":   fmt.Sprintf("Test key should be %v. Got value %v instead.", testNewValue, value),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Test overwrite successful",
			})
		})

		v1.GET("/config", func(c *gin.Context) {
			// Endpoint to check the current configuration
			c.JSON(http.StatusOK, gin.H{
				"kv_api_v1_base_url": kvAPIv1BaseURL,
			})
		})
	}

	return r
}
