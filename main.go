package main

import (
	"log"

	"github.com/polifev/smarthome-p1-receiver/parser"
	"github.com/polifev/smarthome-p1-receiver/reader"
	"github.com/polifev/smarthome-p1-receiver/store"
)

func main() {
	log.Println("[INIT] connecting to database...")
	chConfig, err := store.LoadClickHouseConfig("clickhouse.yaml")
	if err != nil {
		log.Fatalf("unable to read clickhouse config: %v", err)
	}

	s, err := store.NewClickHouseStore(chConfig)
	if err != nil {
		log.Fatalf("unable to create store: %v", err)
	}
	defer s.Close()
	log.Println("[INIT] connected to database !")

	p := parser.NewParser()

	// Connect MQTT
	//log.Println("[INIT] connecting to MQTT broker...")
	//mqttConfig, err := reader.LoadMqttConfig("mqtt.yaml")
	//if err != nil {
	//	log.Fatalf("unable to read mqtt config: %v", err)
	//}
	//
	//r, err := reader.NewMqttReader(mqttConfig)
	//if err != nil {
	//	log.Fatalf("unable to create mqtt reader: %v", err)
	//}

	// Connect to ESP8266
	tcpConfig, err := reader.LoadTcpConfig("tcp.yaml")
	if err != nil {
		log.Fatalf("unable to read tcp config: %v", err)
	}

	r, err := reader.NewTcpReader(tcpConfig)
	if err != nil {
		log.Fatalf("unable to create tcp reader: %v", err)
	}

	defer r.Close()

	// Read messages
	for msg := range r.GetInputChan() {
		powerData := p.ParsePayload(msg)
		err := s.PutData(powerData)
		if err != nil {
			log.Printf("error inserting data: %v", err)
		}
	}
}
