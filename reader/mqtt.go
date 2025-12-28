package reader

import (
	"fmt"
	"io"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v3"
)

type MqttConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	ClientId string `yaml:"client-id"`
}

type MqttReader struct {
	Input chan []byte
}

func DefaultMqttConfig() MqttConfig {
	return MqttConfig{
		Host:     "localhost",
		Port:     "1883",
		ClientId: "p1-receiver",
	}
}

func LoadMqttConfig(filename string) (MqttConfig, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	cfg := DefaultMqttConfig()

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

func NewMqttReader(cfg MqttConfig) (*MqttReader, error) {
	// Create options
	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(fmt.Sprintf("tcp://%s:%s", cfg.Host, cfg.Port))
	mqttOpts.SetClientID(cfg.ClientId)
	mqttOpts.OnConnect = func(client mqtt.Client) {
		log.Println("[INIT] connected to MQTT broker !")
	}

	// Connect
	mqttClient := mqtt.NewClient(mqttOpts)
	t := mqttClient.Connect()
	t.Wait()
	if t.Error() != nil {
		return nil, t.Error()
	}

	inputChan := make(chan []byte)

	// Subscribe to power topic
	t = mqttClient.Subscribe("power", 0, func(client mqtt.Client, msg mqtt.Message) { inputChan <- msg.Payload() })
	t.Wait()
	if t.Error() != nil {
		return nil, t.Error()
	}

	return &MqttReader{Input: inputChan}, nil
}

func (reader *MqttReader) GetInputChan() chan []byte {
	return reader.Input
}
