package tests

import (
	"bytes"
	"encoding/json"
	"golang-dev-logic-challenge/handlers"
	"golang-dev-logic-challenge/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	sampleOptionContract = []models.OptionContract{
		{
			StrikePrice:    100,
			Type:           "Call",
			Bid:            10.05,
			Ask:            12.04,
			LongShort:      "long",
			ExpirationDate: time.Date(2025, time.January, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			StrikePrice:    102.50,
			Type:           "Call",
			Bid:            12.10,
			Ask:            14,
			LongShort:      "long",
			ExpirationDate: time.Date(2025, time.January, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			StrikePrice:    103,
			Type:           "Put",
			Bid:            14,
			Ask:            15.50,
			LongShort:      "short",
			ExpirationDate: time.Date(2025, time.January, 17, 0, 0, 0, 0, time.UTC),
		},
		{
			StrikePrice:    105,
			Type:           "Put",
			Bid:            16,
			Ask:            18,
			LongShort:      "long",
			ExpirationDate: time.Date(2025, time.January, 17, 0, 0, 0, 0, time.UTC),
		},
	}
	sampleGraphData = []models.GraphPoint{
		{X: 50, Y: -12.04},
		{X: 60, Y: -12.04},
		{X: 70, Y: -12.04},
		{X: 80, Y: -12.04},
		{X: 90, Y: -12.04},
		{X: 100, Y: -12.04},
		{X: 110, Y: -2.04},
		{X: 120, Y: 7.96},
		{X: 130, Y: 17.96},
		{X: 140, Y: 27.96},
	}
	sampleMaxProfit      = "unlimited"
	sampleMaxLoss        = 12.04
	sampleBreakEvenPoint = 112.04
)

func TestOptionsContractModelValidation(t *testing.T) {
	assert.NoError(t, models.ValidateOptionsContract(sampleOptionContract))
}

func TestAnalysisEndpoint(t *testing.T) {
	assert.Equal(t, sampleGraphData, handlers.CalculateGraphData(sampleOptionContract[0]))
	assert.Equal(t, sampleMaxProfit, handlers.CalculateMaxProfit(sampleOptionContract[0]))
	assert.Equal(t, sampleMaxLoss, handlers.CalculateMaxLoss(sampleOptionContract[0]))
	assert.Equal(t, sampleBreakEvenPoint, handlers.CalculateBreakEvenPoints(sampleOptionContract[0]))
}

func TestIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/analyze", handlers.AnalysisHandler)

	jsonData, err := json.Marshal([]models.OptionContract{
		sampleOptionContract[0],
	})
	if !assert.Equal(t, nil, err) {
		return
	}

	// Create a new HTTP recorder to capture the response
	recorder := httptest.NewRecorder()

	// Create a new HTTP POST request to the /analyze endpoint with the JSON data as the request body
	request, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	// Check if the response status code is 200 (OK)
	if !assert.Equal(t, http.StatusOK, recorder.Code) {
		return
	}

	analysisResponses := []handlers.AnalysisResponse{}
	err = json.Unmarshal(recorder.Body.Bytes(), &analysisResponses)
	if !assert.Equal(t, nil, err) {
		return
	}

	// Assert that the analysis response matches the expected values
	assert.Equal(t, []handlers.AnalysisResponse{
		{
			GraphData:       sampleGraphData,
			MaxProfit:       sampleMaxProfit,
			MaxLoss:         sampleMaxLoss,
			BreakEvenPoints: sampleBreakEvenPoint,
		},
	}, analysisResponses)
}
