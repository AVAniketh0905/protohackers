package cmd

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
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
	ErrorType         MsgType = 16  //0x10
	PlateType         MsgType = 32  //0x20
	TicketType        MsgType = 33  //0x21
	WantHeartBeatType MsgType = 64  //0x40
	HeartBeatType     MsgType = 65  //0x41
	IAmCameraType     MsgType = 128 //0x80
	IAmDispatcherType MsgType = 129 //0x81
)

type Message interface {
	Type() MsgType
	Unmarshall(data []byte) error
	Marshall() ([]byte, error)
}

type Error struct {
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
		return fmt.Errorf("mismatch Msg Type")
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
		return fmt.Errorf("mismatch Msg Type")
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
		return fmt.Errorf("mismatch Msg Type")
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
	data = append(data, byte(t.Type()))
	plateStr, err := t.Plate.Marshall()
	if err != nil {
		return nil, err
	}
	data = append(data, plateStr...)
	data = binary.BigEndian.AppendUint16(data, t.Road)
	data = binary.BigEndian.AppendUint16(data, t.Mile1)
	data = binary.BigEndian.AppendUint32(data, t.Timestamp1)
	data = binary.BigEndian.AppendUint16(data, t.Mile2)
	data = binary.BigEndian.AppendUint32(data, t.Timestamp2)
	data = binary.BigEndian.AppendUint16(data, t.Speed)
	return data, nil
}

type WantHeartBeat struct {
	Interval uint32
}

func (hb WantHeartBeat) Type() MsgType {
	return WantHeartBeatType
}

func (whb *WantHeartBeat) Unmarshall(data []byte) error {
	r := bytes.NewReader(data)
	msgTypebyte, err := r.ReadByte()
	if err != nil {
		return err
	}

	msgType := MsgType(msgTypebyte)
	if msgType != whb.Type() {
		return fmt.Errorf("mismatch Msg Type")
	}

	rmgBytes := make([]byte, 1024)
	n, err := r.Read(rmgBytes)
	if err != nil {
		return err
	}

	whb.Interval = binary.BigEndian.Uint32(rmgBytes[:n])
	return nil
}

func (whb WantHeartBeat) Marshall() (data []byte, err error) {
	data = append(data, byte(whb.Type()))
	data = binary.BigEndian.AppendUint32(data, whb.Interval)
	return data, nil
}

type HeartBeat struct {
}

func (hb HeartBeat) Type() MsgType {
	return HeartBeatType
}

type IAmCamera struct {
	Road  uint16
	Mile  uint16
	Limit uint16 // miles per hour
}

func (c IAmCamera) Type() MsgType {
	return IAmCameraType
}

func (c *IAmCamera) Unmarshall(data []byte) error {
	r := bytes.NewReader(data)
	msgTypebyte, err := r.ReadByte()
	if err != nil {
		return err
	}

	msgType := MsgType(msgTypebyte)
	if msgType != c.Type() {
		return fmt.Errorf("mismatch Msg Type")
	}

	rmgBytes := make([]byte, 1024)
	n, err := r.Read(rmgBytes)
	if err != nil {
		return err
	}

	m := 0

	c.Road = binary.BigEndian.Uint16(rmgBytes[m : m+2])
	m += 2

	c.Mile = binary.BigEndian.Uint16(rmgBytes[m : m+2])
	m += 2

	c.Limit = binary.BigEndian.Uint16(rmgBytes[m:n])

	return nil
}

func (c IAmCamera) Marshall() (data []byte, err error) {
	data = append(data, byte(c.Type()))
	data = binary.BigEndian.AppendUint16(data, c.Road)
	data = binary.BigEndian.AppendUint16(data, c.Mile)
	data = binary.BigEndian.AppendUint16(data, c.Limit)
	return data, nil
}

type IAmDispatcher struct {
	NumRoads uint8
	Roads    []uint16
}

func (d IAmDispatcher) Type() MsgType {
	return IAmDispatcherType
}

func (d *IAmDispatcher) Unmarshall(data []byte) error {
	r := bytes.NewReader(data)
	msgTypebyte, err := r.ReadByte()
	if err != nil {
		return err
	}

	msgType := MsgType(msgTypebyte)
	if msgType != d.Type() {
		return fmt.Errorf("mismtach Msg Type")
	}

	rmgBytes := make([]byte, 1024)
	_, err = r.Read(rmgBytes)
	if err != nil {
		return err
	}
	m := 0

	d.NumRoads = rmgBytes[0]
	m += 1

	for range d.NumRoads {
		road := binary.BigEndian.Uint16(rmgBytes[m : m+2])
		d.Roads = append(d.Roads, road)
		m += 2
	}
	return nil
}

func (d IAmDispatcher) Marshall() (data []byte, err error) {
	data = append(data, byte(d.Type()))
	data = append(data, byte(d.NumRoads))
	for i := range d.NumRoads {
		data = binary.BigEndian.AppendUint16(data, d.Roads[i])
	}
	return data, nil
}

func (sd SpeedDeamon) Setup() context.Context { return context.TODO() }

func (sd SpeedDeamon) Handler(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		var data []byte = buf[:n]
		log.Println("data: ", data)

		switch MsgType(data[0]) {
		case PlateType:
			var plateMsg Plate
			err := plateMsg.Unmarshall(data)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("plate: ", plateMsg)
		case IAmCameraType:
			var camMsg IAmCamera
			err := camMsg.Unmarshall(data)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("camera: ", camMsg)
		case IAmDispatcherType:
			var dispatcherMsg IAmDispatcher
			err := dispatcherMsg.Unmarshall(data)
			if err != nil {
				log.Println(err)
			}

			log.Println("ticket dispatcher: ", dispatcherMsg)

			var ticket Ticket // TODO: sample msg
			ticket.Plate = Str{Len: 4, Msg: []byte("UN1X")}
			ticket.Road = 123
			ticket.Mile1 = 8
			ticket.Timestamp1 = 0
			ticket.Mile2 = 9
			ticket.Timestamp1 = 45
			ticket.Speed = ((ticket.Mile2 - ticket.Mile1) / (uint16(ticket.Timestamp2) - uint16(ticket.Timestamp1)) * 3600)

			tick, err := ticket.Marshall()
			if err != nil {
				log.Fatal(err)
			}

			_, err = conn.Write(tick)
			if err != nil {
				log.Fatal(err)
			}
		default: // ErrorType
			msg := []byte("illegal msg")
			errMsg := Error{
				err: Str{
					Len: uint8(len(msg)),
					Msg: msg,
				},
			}

			data, err := errMsg.Marshall()
			if err != nil {
				log.Fatal(err)
			}

			_, err = conn.Write(data)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("terminating exsisting connection...")
			return
		}
	}
}

func Run() {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	s := SpeedDeamon{cfg}

	internal.RunTCP(s)
}
