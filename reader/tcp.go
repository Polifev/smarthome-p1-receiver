package reader

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

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
	d := net.Dialer{
		Timeout: 1 * time.Second,
	}
	conn, err := d.Dial("tcp", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)))
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
			log.Printf("[TCP] reading...")
			err := conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			if err != nil {
				log.Println("[TCP] set deadline error:", err)
				tcpReader.Close()
				return
			}
			message, err := reader.ReadString('!')

			if err != nil {
				log.Println("[TCP] read error:", err)
				tcpReader.Close()
				return
			}

			data := []byte(message)
			tcpReader.input <- data
			log.Printf("[TCP] %d bytes read !", len(data))
		}
	}()

	return tcpReader, nil
}

func (r *TcpReader) GetInputChan() chan []byte {
	return r.input
}

func (r *TcpReader) Close() error {
	r.stopped = true
	err := r.conn.Close()
	close(r.input)
	return err
}
