package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/test_deletion", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Test deletion."})
		})

		v1.GET("/test_overwrite", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Test overwrite."})
		})
	}

	return r
}
