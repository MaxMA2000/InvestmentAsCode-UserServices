package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)


func main() {
    router := gin.Default()

    router.GET("/fetch-assets", fetchAssets)

    router.Run("localhost:9090")
}


func fetchAssets(c *gin.Context) {
		// Make a GET request to the external API
		response, err := http.Get("http://localhost:3000/v1/asset")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assets"})
			return
		}
		defer response.Body.Close()

		// Check if the request was successful
		if response.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assets"})
			return
		}

		// Parse the response body as JSON
		var assets []map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&assets); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		// Return the assets as the response
		c.JSON(http.StatusOK, assets)
}

