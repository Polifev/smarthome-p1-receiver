package parser

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/polifev/smarthome-p1-receiver/model"
)

type Parser struct {
	pattern *regexp.Regexp
}

func NewParser() *Parser {
	return &Parser{
		pattern: regexp.MustCompile("^([0-9]-[0-9]:[0-9]+\\.[0-9]+\\.[0-9])\\((.*)\\)$"),
	}
}

func (p *Parser) ParsePayload(rawPayload []byte) model.PowerData {
	payload := string(rawPayload)
	payload = strings.TrimSpace(payload)

	data := model.PowerData{}
	for _, line := range strings.Split(payload, "\n") {
		line = strings.TrimSpace(line)
		groups := p.pattern.FindStringSubmatch(line)
		if groups == nil {
			continue
		}
		switch groups[1] {
		// Total accumulated power
		case "1-0:1.8.1":
			data.HighPriceConsumption = parseKw(groups[2])
		case "1-0:1.8.2":
			data.LowPriceConsumption = parseKw(groups[2])
		case "1-0:2.8.1":
			data.HighPriceProduction = parseKw(groups[2])
		case "1-0:2.8.2":
			data.LowPriceProduction = parseKw(groups[2])

		// Current power consumption
		case "1-0:1.7.0":
			data.CurrentPowerConsumption = parseKw(groups[2])
		case "1-0:21.7.0":
			data.CurrentPowerConsumptionP1 = parseKw(groups[2])
		case "1-0:41.7.0":
			data.CurrentPowerConsumptionP2 = parseKw(groups[2])
		case "1-0:61.7.0":
			data.CurrentPowerConsumptionP3 = parseKw(groups[2])

		// Current power production
		case "1-0:2.7.0":
			data.CurrentPowerProduction = parseKw(groups[2])
		case "1-0:22.7.0":
			data.CurrentPowerProductionP1 = parseKw(groups[2])
		case "1-0:42.7.0":
			data.CurrentPowerProductionP2 = parseKw(groups[2])
		case "1-0:62.7.0":
			data.CurrentPowerProductionP3 = parseKw(groups[2])

		// Time
		case "0-0:1.0.0":
			data.TimeStamp = parseTime(groups[2])

		// Tariff
		case "0-0:96.14.0":
			data.CurrentPrice = parsePrice(groups[2])
		}
	}
	return data
}

func parseKw(kwStr string) float64 {
	parts := strings.Split(kwStr, "*")
	val, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return math.NaN()
	}
	return val
}

func parseTime(str string) time.Time {
	const format = "060102150405S"
	location, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		panic(err)
	}

	t, err := time.ParseInLocation(format, str, location)
	if err != nil {
		return time.Time{}
	}
	return t
}

func parsePrice(str string) model.Price {
	val, err := strconv.ParseInt(str, 10, 8)
	if err != nil {
		return 0
	}
	return model.Price(val)
}
