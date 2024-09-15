package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/AVAniketh0905/protohackers/internal"
)

// helper functions for converting hex string to rwquired Binary format
type Binary string // hex string

type BinReq struct {
	Char string
	Num1 int32
	Num2 int32
}

func hextoi32(b string) (int32, error) {
	num, err := strconv.ParseUint(b, 16, 32)
	if err != nil {
		return 0, err
	}

	return int32(num), nil
}

func (b Binary) Unmarshall(req *BinReq) error {
	if len(b) > 18 {
		return fmt.Errorf("only expected 18 bytes but got, %v", len(b))
	}

	// char 1 byte
	char, err := strconv.ParseInt(string(b[:2]), 16, 32)
	if err != nil {
		return err
	}
	switch char {
	case 81:
		req.Char = "Q"
	case 73:
		req.Char = "I"
	default:
		return fmt.Errorf("only supports Q/I, got %v", char)
	}

	num1, err := hextoi32(string(b[2:10]))
	if err != nil {
		return err
	}
	req.Num1 = num1

	num2, err := hextoi32(string(b[10:18]))
	if err != nil {
		return err
	}
	req.Num2 = num2

	return nil
}

type MeansToEnd struct{ *internal.Config }

func (MeansToEnd) Setup() context.Context { return context.TODO() }

func (MeansToEnd) Handler(_ context.Context, conn net.Conn) {
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		var req BinReq
		hexStr := Binary(buf[:n])
		err = hexStr.Unmarshall(&req)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(req)
		_, err = conn.Write([]byte(req.Char))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Run() {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	s := MeansToEnd{cfg}

	internal.RunTCP(s)
}
