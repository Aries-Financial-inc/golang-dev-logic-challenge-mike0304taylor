package handlers

import (
	"golang-dev-logic-challenge/models"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

func round(value float64) float64 {
	return math.Round(value*100) / 100
}

func CalculateGraphData(contract models.OptionContract) []models.GraphPoint {
	graphData := []models.GraphPoint{}

	for i := 50.0; i < 150.0; i += 10.0 {
		marketPrice := contract.StrikePrice * i / 100
		reward := 0.0
		if contract.Type == "Call" && contract.LongShort == "long" {
			reward = max(marketPrice-contract.StrikePrice-contract.Ask, -contract.Ask)
		} else if contract.Type == "Call" && contract.LongShort == "short" {
			reward = min(contract.StrikePrice+contract.Bid-marketPrice, contract.Bid)
		} else if contract.Type == "Put" && contract.LongShort == "long" {
			reward = max(contract.StrikePrice-marketPrice-contract.Ask, -contract.Ask)
		} else if contract.Type == "Put" && contract.LongShort == "short" {
			reward = min(contract.Bid+marketPrice-contract.StrikePrice, contract.Bid)
		}
		graphData = append(graphData, models.GraphPoint{
			X: round(marketPrice),
			Y: round(reward),
		})
	}

	return graphData
}

func CalculateMaxProfit(contract models.OptionContract) any {
	maxProfit := 0.0
	if contract.Type == "Call" && contract.LongShort == "long" {
		return "unlimited"
	} else if contract.Type == "Call" && contract.LongShort == "short" {
		maxProfit = contract.Bid
	} else if contract.Type == "Put" && contract.LongShort == "long" {
		return "unlimited"
	} else if contract.Type == "Put" && contract.LongShort == "short" {
		maxProfit = contract.Bid
	}
	return round(maxProfit)
}

func CalculateMaxLoss(contract models.OptionContract) any {
	maxLoss := 0.0
	if contract.Type == "Call" && contract.LongShort == "long" {
		maxLoss = contract.Ask
	} else if contract.Type == "Call" && contract.LongShort == "short" {
		return "unlimited"
	} else if contract.Type == "Put" && contract.LongShort == "long" {
		maxLoss = contract.Ask
	} else if contract.Type == "Put" && contract.LongShort == "short" {
		return "unlimited"
	}
	return round(maxLoss)
}

func CalculateBreakEvenPoints(contract models.OptionContract) float64 {
	breakEvenPoint := 0.0
	if contract.Type == "Call" && contract.LongShort == "long" {
		breakEvenPoint = contract.StrikePrice + contract.Ask
	} else if contract.Type == "Call" && contract.LongShort == "short" {
		breakEvenPoint = contract.StrikePrice + contract.Bid
	} else if contract.Type == "Put" && contract.LongShort == "long" {
		breakEvenPoint = contract.StrikePrice - contract.Ask
	} else if contract.Type == "Put" && contract.LongShort == "short" {
		breakEvenPoint = contract.StrikePrice - contract.Bid
	}
	return round(breakEvenPoint)
}

// AnalysisResponse represents the data structure of the analysis result
type AnalysisResponse struct {
	GraphData       []models.GraphPoint `json:"graph_data"`
	MaxProfit       any                 `json:"max_profit"`
	MaxLoss         any                 `json:"max_loss"`
	BreakEvenPoints float64             `json:"break_even_points"`
}

// PingExample godoc
// @Summary Analyze option contracts
// @Schemes
// @Description Get graph, maximum profit and loss and break even point
// @Tags Analysis
// @Accept json
// @Produce json
// @Param params body []models.OptionContract true "Option contracts"
// @Success 200 {object} []AnalysisResponse
// @Router /analyze [post]
func AnalysisHandler(context *gin.Context) {
	contracts := []models.OptionContract{}

	err := context.ShouldBindJSON(&contracts)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, "The option contracts you have attempted to access are invalid.")
		return
	}

	err = models.ValidateOptionsContract(contracts)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, "Validation Error: "+err.Error()+".")
		return
	}

	analysisResponses := []AnalysisResponse{}
	for _, contract := range contracts {
		analysisResponses = append(analysisResponses, AnalysisResponse{
			GraphData:       CalculateGraphData(contract),
			MaxProfit:       CalculateMaxProfit(contract),
			MaxLoss:         CalculateMaxLoss(contract),
			BreakEvenPoints: CalculateBreakEvenPoints(contract),
		})
	}

	context.IndentedJSON(http.StatusOK, analysisResponses)
}
