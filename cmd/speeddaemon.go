package cmd

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/AVAniketh0905/protohackers/internal"
)

type SpeedDeamon struct{ *internal.Config }

type Str struct {
	Len uint8
	Msg []uint8
}

func (s *Str) Unmarshall(data []byte) error {
	r := bytes.NewBuffer(data)

	len, err := r.ReadByte()
	if err != nil {
		return err
	}
	s.Len = len

	var msg []uint8
	for range len {
		b, _ := r.ReadByte()
		msg = append(msg, b)
	}
	s.Msg = msg

	return nil
}

func (s Str) Marshall() (b []byte, err error) {
	len := byte(s.Len)
	b = append(b, len)

	b = append(b, s.Msg...)
	return b, err
}

type MsgType uint8

const (
	ErrorType  MsgType = 16 //0x10
	PlateType  MsgType = 32 //0x20
	TicketType MsgType = 33 //0x21
)

type Message interface {
	Type() MsgType
	Unmarshall(data []byte) error
	Marshall() ([]byte, error)
}

type Error struct {
	Message
	err Str
}

func (e Error) Type() MsgType {
	return ErrorType
}

func (e Error) ErrMsg() string {
	return string(e.err.Msg)
}

func (e *Error) Unmarshall(data []byte) error {
	r := bytes.NewReader(data)
	msgTypebyte, err := r.ReadByte()
	if err != nil {
		return err
	}

	msgType := MsgType(msgTypebyte)
	if msgType != e.Type() {
		return fmt.Errorf("mismtach Msg Type")
	}

	rmgBytes := make([]byte, 1024)
	n, err := r.Read(rmgBytes)
	if err != nil {
		return err
	}

	var msg Str
	if err := msg.Unmarshall(rmgBytes[:n]); err != nil {
		return err
	}

	e.err = msg
	return nil
}

func (e Error) Marshall() (p []byte, err error) {
	p = append(p, byte(e.Type())) // type
	msgStr, err := e.err.Marshall()
	if err != nil {
		return nil, err
	}
	p = append(p, msgStr...)
	return p, nil
}

type Plate struct {
	Message
	Plate     Str
	Timestamp uint32
}

func (p Plate) Type() MsgType {
	return PlateType
}

func (p Plate) PlateMsg() string {
	return string(p.Plate.Msg)
}

func (p *Plate) Unmarshall(data []byte) error {
	r := bytes.NewReader(data)
	msgTypebyte, err := r.ReadByte()
	if err != nil {
		return err
	}

	msgType := MsgType(msgTypebyte)
	if msgType != p.Type() {
		return fmt.Errorf("mismtach Msg Type")
	}

	rmgBytes := make([]byte, 1024)
	n, err := r.Read(rmgBytes)
	if err != nil {
		return err
	}

	var plate Str
	if err := plate.Unmarshall(rmgBytes[:n]); err != nil {
		return err
	}

	timestamp := binary.BigEndian.Uint32(rmgBytes[plate.Len+1 : n])

	p.Plate = plate
	p.Timestamp = timestamp
	return nil
}

func (p Plate) Marshall() (data []byte, err error) {
	data = append(data, byte(p.Type()))
	plateStr, err := p.Plate.Marshall()
	if err != nil {
		return nil, err
	}
	data = append(data, plateStr...)
	binary.BigEndian.AppendUint32(data, p.Timestamp)

	return data, nil
}

type Ticket struct {
	Message
	Plate      Str
	Road       uint16
	Mile1      uint16
	Timestamp1 uint32
	Mile2      uint16
	Timestamp2 uint32
	Speed      uint16 // 100x miles per hour
}

func (t Ticket) Type() MsgType {
	return TicketType
}

func (t Ticket) PlateMsg() string {
	return string(t.Plate.Msg)
}

func (t *Ticket) Unmarshall(data []byte) error {
	r := bytes.NewReader(data)
	msgTypebyte, err := r.ReadByte()
	if err != nil {
		return err
	}

	msgType := MsgType(msgTypebyte)
	if msgType != t.Type() {
		return fmt.Errorf("mismtach Msg Type")
	}

	rmgBytes := make([]byte, 1024)
	n, err := r.Read(rmgBytes)
	if err != nil {
		return err
	}
	m := 0

	var plate Str
	if err := plate.Unmarshall(rmgBytes[:n]); err != nil {
		return err
	}
	t.Plate = plate
	m += int(plate.Len) + 1

	t.Road = binary.BigEndian.Uint16(rmgBytes[m : m+2])
	m += 2

	t.Mile1 = binary.BigEndian.Uint16(rmgBytes[m : m+2])
	m += 2

	t.Timestamp1 = binary.BigEndian.Uint32(rmgBytes[m : m+4])
	m += 4

	t.Mile2 = binary.BigEndian.Uint16(rmgBytes[m : m+2])
	m += 2

	t.Timestamp2 = binary.BigEndian.Uint32(rmgBytes[m : m+4])
	m += 4

	t.Speed = binary.BigEndian.Uint16(rmgBytes[m:n])

	return nil
}

func (t Ticket) Marshall() (data []byte, err error) {
	return data, nil
}

func (sd SpeedDeamon) Setup() context.Context { return context.TODO() }

func (sd SpeedDeamon) Handler(ctx context.Context, conn net.Conn) {

}

func Run() {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	s := SpeedDeamon{cfg}

	internal.RunTCP(s)
}
