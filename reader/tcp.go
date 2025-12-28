package reader

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

type TcpReader struct {
	input   chan []byte
	conn    *net.TCPConn
	stopped bool
}

type TcpConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func DefaultTcpConfig() TcpConfig {
	return TcpConfig{
		Host: "10.140.40.198",
		Port: 6969,
	}
}

func LoadTcpConfig(filename string) (TcpConfig, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	cfg := DefaultTcpConfig()

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

func NewTcpReader(cfg TcpConfig) (*TcpReader, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}

	tcpReader := &TcpReader{
		input:   make(chan []byte),
		conn:    conn.(*net.TCPConn),
		stopped: false,
	}

	go func() {
		reader := bufio.NewReader(conn)
		for !tcpReader.stopped {
			message, err := reader.ReadString('!')
			if err != nil {
				log.Println("read error:", err)
				continue
			}
			tcpReader.input <- []byte(message)
		}
	}()

	return tcpReader, nil
}

func (r *TcpReader) GetInputChan() chan []byte {
	return r.input
}

func (r *TcpReader) Close() error {
	r.stopped = true
	return r.conn.Close()
}
