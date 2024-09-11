package internal

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
)

type Configuration interface {
	Port() int
	Addr() string
	ParseFlags()
}

type Config struct{ port int }

func NewConfig(defaultPort int) *Config {
	return &Config{port: defaultPort}
}

func (cfg *Config) Port() int { return cfg.port }

func (cfg *Config) Addr() string { return fmt.Sprintf("0.0.0.0:%v", cfg.port) }

func (cfg *Config) ParseFlags() {
	flag.IntVar(&cfg.port, "port", cfg.port, "port to listen on")
	flag.Parse()
}

type Server interface {
	Configuration

	Setup() context.Context
	Handler(ctx context.Context, conn net.Conn)
}

func RunTCP(s Server) {
	addr := s.Addr()
	ctx := s.Setup()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("l: ", err)
	}
	// defer l.Close()
	log.Printf("listener started at, %v\n", addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("conn started at, ", conn.LocalAddr())

		go s.Handler(ctx, conn)
	}
}
