package model

import "time"

type Price int

const (
	High Price = 1 + iota
	Low
)

type PowerData struct {
	TimeStamp            time.Time
	LowPriceConsumption  float64
	HighPriceConsumption float64
	LowPriceProduction   float64
	HighPriceProduction  float64

	CurrentPowerConsumption   float64
	CurrentPowerConsumptionP1 float64
	CurrentPowerConsumptionP2 float64
	CurrentPowerConsumptionP3 float64

	CurrentPowerProduction   float64
	CurrentPowerProductionP1 float64
	CurrentPowerProductionP2 float64
	CurrentPowerProductionP3 float64

	CurrentPrice Price
}
