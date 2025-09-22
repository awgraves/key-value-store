package main

import (
	"net/http"

	"github.com/awgraves/key-value-store/kv_service/store"
	"github.com/gin-gonic/gin"
)

func setupRouter(kvStore store.Store) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		keys := v1.Group("/keys")
		{
			keys.GET("/:key", func(c *gin.Context) {
				key := c.Param("key")
				value := kvStore.Get(key)
				c.JSON(http.StatusOK, gin.H{"value": value})
			})

			keys.POST("/:key", func(c *gin.Context) {
				key := c.Param("key")
				var request struct {
					Value any `json:"value" binding:"required"`
				}
				if err := c.ShouldBindJSON(&request); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				kvStore.Set(key, request.Value)
				c.JSON(http.StatusOK, gin.H{"message": "Key set."})
			})

			keys.DELETE("/:key", func(c *gin.Context) {
				key := c.Param("key")
				kvStore.Delete(key)
				c.JSON(http.StatusOK, gin.H{"message": "Key deleted."})
			})
		}
	}

	return r
}
