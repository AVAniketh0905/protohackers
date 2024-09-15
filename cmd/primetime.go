package cmd

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"math"
	"net"
	"strings"

	"github.com/AVAniketh0905/protohackers/internal"
)

type PrimeTime struct{ *internal.Config }

type Req struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

type Resp struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func isPrime(num int) bool {
	for i := 2; i <= int(math.Sqrt(float64(num))); i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func (PrimeTime) Setup() context.Context { return context.TODO() }

func (PrimeTime) Handler(_ context.Context, conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 4096)

		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatalf("failed to read from conn, %v", err)
			}
			break
		}

		content := string(buf[:n])
		reqStr := strings.Split(content, "\n")[0]

		var req Req
		err = json.Unmarshal([]byte(reqStr), &req)
		if err != nil {
			log.Fatalf("failed to unmarshal json req, %v", err)
		}

		var res Resp
		res.Method = "isPrime"
		res.Prime = isPrime(req.Number)

		respBytes, err := json.Marshal(res)
		if err != nil {
			log.Fatalf("failed to marshal json resp, %v", err)
		}
		respBytes = append(respBytes, '\n')

		_, err = conn.Write(respBytes)
		if err != nil {
			log.Fatalf("failed to write to conn, %v", err)
		}
	}
}

func Run() {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	s := PrimeTime{cfg}

	internal.RunTCP(s)
}
