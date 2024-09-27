package cmd

import (
	"context"
	"net"

	"github.com/AVAniketh0905/protohackers/internal"
)

type BogusCoin struct{ *internal.Config }

func (b BogusCoin) Setup() context.Context { return context.TODO() }

func (b BogusCoin) Handler(_ context.Context, conn net.Conn) {}

func Run() {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	s := BogusCoin{cfg}

	internal.RunTCP(s)
}
