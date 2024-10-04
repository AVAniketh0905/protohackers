package tests

import (
	"encoding/hex"
	"log"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/AVAniketh0905/protohackers/cmd"
	"github.com/AVAniketh0905/protohackers/internal"
)

func TestStr(t *testing.T) {
	hexMsgs := []string{"00", "03666f6f", "08456C626572657468"}
	msg := []string{"", "foo", "Elbereth"}

	for i, h := range hexMsgs {
		bytes, err := hex.DecodeString(h)
		if err != nil {
			t.Error(i, h, err)
		}

		var dataStr cmd.Str
		dataStr.Unmarshall(bytes)

		if len(msg[i]) != int(dataStr.Len) {
			t.Error(i, h, "len does not match")
		}

		if !reflect.DeepEqual(msg[i], string(dataStr.Msg)) {
			t.Errorf("i: %d h: %q, msgs doesnt match, expected %q, got %q", i, h, msg[i], string(dataStr.Msg))
		}
	}
}

func TestMsgTypes(t *testing.T) {
	hexMsgs := []string{"1003626164", "100b696c6c6567616c206d7367"}
	errMsgs := []string{"bad", "illegal msg"}

	for i, hm := range hexMsgs {
		hexBytes, err := hex.DecodeString(hm)
		if err != nil {
			t.Error(err)
		}
		var errMsg cmd.Error
		err = errMsg.Unmarshall(hexBytes)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(errMsg.ErrMsg(), errMsgs[i]) {
			t.Errorf("msgs did not match, expected %q, got %q", errMsgs[i], errMsg.ErrMsg())
		}
	}

	hexMsgs = []string{"2004554e3158000003e8", "200752453035424b470001e240"}
	plateMsgs := []struct {
		plate     string
		timestamp uint32
	}{
		{
			plate:     "UN1X",
			timestamp: 1000,
		},
		{
			plate:     "RE05BKG",
			timestamp: 123456,
		},
	}

	for i, hm := range hexMsgs {
		hexBytes, err := hex.DecodeString(hm)
		if err != nil {
			t.Error(err)
		}

		var plateMsg cmd.Plate
		err = plateMsg.Unmarshall(hexBytes)
		if err != nil {
			t.Error(err)
		}

		if plateMsg.Type() != cmd.PlateType {
			t.Error("mismatch Msg Type")
		}

		if !reflect.DeepEqual(plateMsg.PlateMsg(), plateMsgs[i].plate) {
			t.Errorf("msgs did not match, expected %q, got %q", plateMsgs[i].plate, plateMsg.PlateMsg())
		}

		if !reflect.DeepEqual(plateMsg.Timestamp, plateMsgs[i].timestamp) {
			t.Errorf("msgs did not match, expected %v, got %v", plateMsgs[i].timestamp, plateMsg.Timestamp)
		}
	}

	hexMsgs = []string{"2104554e3158007b00080000000000090000002d1f40"}
	ticketMsgs := []struct {
		plate      string
		road       uint16
		mile1      uint16
		timestamp1 uint32
		mile2      uint16
		timestamp2 uint32
		speed      uint16
	}{
		{plate: "UN1X", road: 123, mile1: 8, timestamp1: 0, mile2: 9, timestamp2: 45, speed: 8000},
	}

	for i, hm := range hexMsgs {
		hexBytes, err := hex.DecodeString(hm)
		if err != nil {
			t.Error(err)
		}

		var ticket cmd.Ticket
		if err := ticket.Unmarshall(hexBytes); err != nil {
			t.Error(err)
		}

		if ticket.Type() != cmd.TicketType {
			t.Error("mismatch Msg Type")
		}

		if !reflect.DeepEqual(ticketMsgs[i].plate, ticket.PlateMsg()) {
			t.Errorf("msgs did not match, expected %v, got %v", ticketMsgs[i].plate, ticket.PlateMsg())
		}

		if !reflect.DeepEqual(ticketMsgs[i].speed, ticket.Speed) {
			t.Errorf("msgs did not match, expected %v, got %v", ticketMsgs[i].speed, ticket.Speed)
		}
	}
}

func TestCameraType(t *testing.T) {
	hexMsgs := []string{"80007b0008003c"}
	camMsgs := []struct {
		road  uint16
		mile  uint16
		limit uint16
	}{
		{road: 123, mile: 6, limit: 60},
	}

	for i, hm := range hexMsgs {
		hexBytes, err := hex.DecodeString(hm)
		if err != nil {
			t.Error(err)
		}
		var camMsg cmd.IAmCamera
		err = camMsg.Unmarshall(hexBytes)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(camMsg.Road, camMsgs[i].road) {
			t.Errorf("msgs did not match, expected %d, got %d", camMsgs[i], camMsg)
		}
	}
}

func TestServer(t *testing.T) {
	// TODO
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	client1, err := net.Dial("tcp", cfg.Addr())
	if err != nil {
		t.Error(err)
	}

	c1Msgs := []string{
		"80007b0008003c",
		"2004554e315800000000",
		"80007b0009003c",
		"2004554e31580000002d",
		"8101007b",
	}

	for _, msg := range c1Msgs {
		data, err := hex.DecodeString(msg)
		if err != nil {
			t.Error(err)
		}
		log.Println("Msg: ", msg, data)
		_, err = client1.Write(data)
		if err != nil {
			t.Error(err)
		}
		time.Sleep(time.Second)
	}
	time.Sleep(time.Second)

	// get ticket
	buf := make([]byte, 1024)
	n, err := client1.Read(buf)
	if err != nil {
		t.Error(err)
	}

	var ticket cmd.Ticket
	err = ticket.Unmarshall(buf[:n])
	if err != nil {
		t.Error(err)
	}

	log.Println(ticket)
}
