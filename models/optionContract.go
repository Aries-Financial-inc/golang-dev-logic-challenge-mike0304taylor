package models

import (
	"errors"
	"time"
)

// OptionContract represents the data structure of an options contract
type OptionContract struct {
	StrikePrice    float64   `json:"strike_price" example:"100"`
	Type           string    `json:"type" example:"Call"`
	Bid            float64   `json:"bid" example:"4"`
	Ask            float64   `json:"ask" example:"6"`
	LongShort      string    `json:"long_short" example:"long"`
	ExpirationDate time.Time `json:"expiration_date" example:"2025-12-31T00:00:00Z"`
}

// Option Contract Validation
func ValidateOptionsContract(contracts []OptionContract) error {
	if len(contracts) > 4 {
		return errors.New("only 4 option contracts can be proccessed in the same time")
	}

	for _, contract := range contracts {
		if contract.LongShort != "short" && contract.LongShort != "long" {
			return errors.New("the long_short field can be short or long")
		} else if contract.Type != "Call" && contract.Type != "Put" {
			return errors.New("the type field can be Call or Put")
		} else if contract.Ask <= contract.Bid {
			return errors.New("the ask field can not be lower than bid field")
		}
	}

	return nil
}
