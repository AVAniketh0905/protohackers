package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
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

type data struct {
	timestamp int32
	price     int32
}

type DB []data

func query(db DB, minTime, maxTime int32) int32 {
	res, size := int32(0), 0
	for _, row := range db {
		if row.timestamp >= minTime && row.timestamp <= maxTime {
			// log.Println(minTime, maxTime, row.timestamp)
			res += row.price
			size += 1
		}
	}

	// log.Println(res, size, res/int32(size))
	return res / int32(size)
}

func (MeansToEnd) Setup() context.Context { return context.TODO() }

func (MeansToEnd) Handler(_ context.Context, conn net.Conn) {
	store := DB{}

	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		for cnt := 18; cnt <= n; cnt += 18 {
			var req BinReq
			hexStr := Binary(buf[cnt-18 : cnt])
			err = hexStr.Unmarshall(&req)
			if err != nil {
				log.Fatal(err)
			}

			// log.Println(req)

			if req.Char == "I" {
				store = append(store, data{req.Num1, req.Num2})
				sort.Slice(store, func(i, j int) bool {
					return store[i].timestamp < store[j].timestamp
				})
			} else if req.Char == "Q" {
				ans := query(store, req.Num1, req.Num2)
				hexAns := fmt.Sprintf("%x", ans)

				_, err = conn.Write([]byte(hexAns))
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// func Run() {
// 	cfg := internal.NewConfig(internal.PORT)
// 	cfg.ParseFlags()

// 	s := MeansToEnd{cfg}

// 	internal.RunTCP(s)
// }
