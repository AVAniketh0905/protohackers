package tests

import (
	"fmt"
	"net"
	"reflect"
	"strconv"
	"testing"

	"github.com/AVAniketh0905/protohackers/cmd"
	"github.com/AVAniketh0905/protohackers/internal"
)

var reqs []struct {
	hex       string
	actualReq cmd.BinReq
} = []struct {
	hex       string
	actualReq cmd.BinReq
}{
	{"51000003e8000186a0", cmd.BinReq{
		Char: "Q",
		Num1: 1000,
		Num2: 100000,
	}},
	{"490000303900000065", cmd.BinReq{
		Char: "I",
		Num1: 12345,
		Num2: 101,
	}},
	{"490000a00000000005", cmd.BinReq{
		Char: "I",
		Num1: 40960,
		Num2: 5,
	}},
	{"510000300000004000", cmd.BinReq{
		Char: "Q",
		Num1: 12288,
		Num2: 16384,
	}},
}

var testReqs []struct {
	hex       string
	actualReq cmd.BinReq
} = []struct {
	hex       string
	actualReq cmd.BinReq
}{
	{"490000303900000065", cmd.BinReq{
		Char: "I",
		Num1: 12345,
		Num2: 101,
	}},
	{"490000303a00000066", cmd.BinReq{
		Char: "I",
		Num1: 12346,
		Num2: 102,
	}},
	{"490000303b00000064", cmd.BinReq{
		Char: "I",
		Num1: 12347,
		Num2: 100,
	}},
	{"490000a00000000005", cmd.BinReq{
		Char: "I",
		Num1: 40960,
		Num2: 5,
	}},
	{"510000300000004000", cmd.BinReq{
		Char: "Q",
		Num1: 12288,
		Num2: 16384,
	}},
}

func TestBinaryRead(t *testing.T) {
	for _, req := range reqs {
		var binReq cmd.BinReq
		hexReq := cmd.Binary(req.hex)
		err := hexReq.Unmarshall(&binReq)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(binReq, req.actualReq) {
			t.Errorf("did not match, got: %v, actual %v", binReq, req.actualReq)
		}
	}
}

func TestTCPServer(t *testing.T) {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	client, err := net.Dial("tcp", cfg.Addr())
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println("client started...")

	for _, req := range testReqs {
		_, err := client.Write([]byte(req.hex))
		if err != nil {
			t.Error(err)
		}

		if req.actualReq.Char == "Q" {
			buf := make([]byte, 4096)
			n, err := client.Read(buf)
			if err != nil {
				t.Error(err)
			}

			ans, err := strconv.ParseInt(string(buf[:n]), 16, 32)
			if err != nil {
				t.Error(err)
			}

			if ans != 101 {
				t.Errorf("answer does not match, got: %d, actual: %d", ans, 101)
			}
		}
	}
}
