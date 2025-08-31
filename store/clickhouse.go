package store

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/polifev/smarthome-p1-receiver/model"

	"gopkg.in/yaml.v3"
)

type ClickHouseStore struct {
	connection driver.Conn
}

type ClickHouseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func DefaultClickHouseConfig() ClickHouseConfig {
	return ClickHouseConfig{
		Host:     "localhost",
		Port:     "9000",
		Username: "default",
		Password: "",
		Database: "default",
	}
}

func LoadClickHouseConfig(filename string) (ClickHouseConfig, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	cfg := DefaultClickHouseConfig()

	if err != nil {
		return cfg, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func NewClickHouseStore(cfg ClickHouseConfig) (*ClickHouseStore, error) {
	ctx := context.Background()
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
	})

	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &ClickHouseStore{connection: conn}, nil
}

func (store *ClickHouseStore) PutData(data model.PowerData) error {
	ctx := context.Background()
	return store.connection.Exec(ctx, `INSERT INTO powers(
		TimeStamp,
		LowPriceConsumption,
		HighPriceConsumption,
		LowPriceProduction,
		HighPriceProduction,
		CurrentPowerConsumption,
		CurrentPowerConsumptionP1,
		CurrentPowerConsumptionP2,
		CurrentPowerConsumptionP3,
		CurrentPowerProduction,
		CurrentPowerProductionP1,
		CurrentPowerProductionP2,
		CurrentPowerProductionP3,
		CurrentPrice
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		data.TimeStamp,
		data.LowPriceConsumption,
		data.HighPriceConsumption,
		data.LowPriceProduction,
		data.HighPriceProduction,
		data.CurrentPowerConsumption,
		data.CurrentPowerConsumptionP1,
		data.CurrentPowerConsumptionP2,
		data.CurrentPowerConsumptionP3,
		data.CurrentPowerProduction,
		data.CurrentPowerProductionP1,
		data.CurrentPowerProductionP2,
		data.CurrentPowerProductionP3,
		data.CurrentPrice,
	)
}

func (store *ClickHouseStore) Close() error {
	return store.connection.Close()
}
