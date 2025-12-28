package main

import (
	"log"
	"time"

	"github.com/polifev/smarthome-p1-receiver/parser"
	"github.com/polifev/smarthome-p1-receiver/reader"
	"github.com/polifev/smarthome-p1-receiver/store"
)

func main() {
	for {
		err := retryLoop()
		if err != nil {
			log.Printf("[ERR] error caused reboot: %v", err)
			log.Printf("[REBOOT] waiting for reboot")
			time.Sleep(5 * time.Second)
		}
	}
}

func retryLoop() error {
	log.Println("[INIT] connecting to database...")
	chConfig, err := store.LoadClickHouseConfig("clickhouse.yaml")
	if err != nil {
		log.Printf("[ERR] unable to read clickhouse config: %v", err)
		return err
	}

	s, err := store.NewClickHouseStore(chConfig)
	if err != nil {
		log.Printf("unable to create store: %v", err)
		return err
	}
	defer s.Close()
	log.Println("[INIT] connected to database !")

	p := parser.NewParser()

	// Connect to ESP8266
	tcpConfig, err := reader.LoadTcpConfig("tcp.yaml")
	if err != nil {
		log.Printf("[ERR] unable to read tcp config: %v", err)
		return err
	}

	r, err := reader.NewTcpReader(tcpConfig)
	if err != nil {
		log.Printf("[ERR] unable to create tcp reader: %v", err)
		return err
	}
	defer r.Close()

	// Read messages
	for msg := range r.GetInputChan() {
		powerData := p.ParsePayload(msg)
		err := s.PutData(powerData)
		if err != nil {
			log.Printf("[ERR] error inserting data: %v", err)
			return err
		}
	}
	return nil
}
