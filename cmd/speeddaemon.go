package cmd

import (
	"bytes"
	"context"
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

func (sd SpeedDeamon) Setup() context.Context { return context.TODO() }

func (sd SpeedDeamon) Handler(ctx context.Context, conn net.Conn) {

}

func Run() {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	s := SpeedDeamon{cfg}

	internal.RunTCP(s)
}
