package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CryptoResponse struct {
	ID             int     `json:"id"`
	AssetID        int     `json:"assetId"`
	Symbol         string  `json:"symbol"`
	Date           string  `json:"date"`
	Open           float64 `json:"open"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Close          float64 `json:"close"`
	AdjClose       float64 `json:"adjClose"`
	Volume         int     `json:"volume"`
	UnadjustedVol  int     `json:"unadjustedVolume"`
	Change         float64 `json:"change"`
	ChangePercent  float64 `json:"changePercent"`
	VWAP           float64 `json:"vwap"`
	Label          string  `json:"label"`
	ChangeOverTime float64 `json:"changeOverTime"`
	AsOfDate       string  `json:"asOfDate"`
}

type StockResponse struct {
	ID             int     `json:"id"`
	AssetID        int     `json:"assetId"`
	Symbol         string  `json:"symbol"`
	Date           string  `json:"date"`
	Open           float64 `json:"open"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Close          float64 `json:"close"`
	AdjClose       float64 `json:"adjClose"`
	Volume         int     `json:"volume"`
	UnadjustedVol  int     `json:"unadjustedVolume"`
	Change         float64 `json:"change"`
	ChangePercent  float64 `json:"changePercent"`
	VWAP           float64 `json:"vwap"`
	Label          string  `json:"label"`
	ChangeOverTime float64 `json:"changeOverTime"`
	AsOfDate       string  `json:"asOfDate"`
}

type AssetPayload struct {
	Date   string            `json:"date"`
	Assets map[string][]Item `json:"assets"`
}

type Item struct {
	ID     int `json:"id"`
	Amount int `json:"amount"`
}

func main() {
	router := gin.Default()
	router.GET("/fetch-assets", fetchAssets)
	router.POST("/process-single-date-assets", processAssets)
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
func processAssets(c *gin.Context) {
	// Declaration and initialization of payload variable
	var payload AssetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Create combinedResponse map and set initial values
	combinedResponse := make(map[string]interface{})
	combinedResponse["date"] = payload.Date
	combinedResponse["responses"] = make(map[string]interface{})

	// Loop through assets in payload
	for assetType, items := range payload.Assets {
		// Create assetResponses slice
		assetResponses := make([]interface{}, 0)

		// Loop through items of each asset type
		for _, item := range items {
			// Declare variables for response and endpoint
			var response interface{}
			var endpoint string

			// Set endpoint and response based on assetType
			switch assetType {
			case "stock":
				endpoint = fmt.Sprintf("http://localhost:8081/v1/stock/%d/%s", item.ID, payload.Date)
				response = &StockResponse{}
			case "crypto":
				endpoint = fmt.Sprintf("http://localhost:5025/v1/crypto/%d/%s", item.ID, payload.Date)
				response = &CryptoResponse{}
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid asset type"})
				return
			}
			// Create HTTP request
			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
				return
			}

			// Send HTTP request and get response
			client := http.DefaultClient
			apiResponse, err := client.Do(req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
				return
			}
			defer apiResponse.Body.Close()

			// Check response status code
			if apiResponse.StatusCode != http.StatusOK {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch asset data"})
				return
			}

			// Read response body
			body, err := ioutil.ReadAll(apiResponse.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
				return
			}
			// Unmarshal response body into the respective response struct
			if err := json.Unmarshal(body, response); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
				return
			}
			// Append response to assetResponses slice
			assetResponses = append(assetResponses, response)
		}

		// Add assetResponses to combinedResponse
		combinedResponse["responses"].(map[string]interface{})[assetType] = assetResponses
	}

	// Return combinedResponse as JSON
	c.JSON(http.StatusOK, combinedResponse)
}
