package cmd

import (
	"context"
	"net"

	"github.com/AVAniketh0905/protohackers/internal"
)

type BudgetChat struct{ *internal.Config }

func (BudgetChat) Setup() context.Context { return context.TODO() }

func (BudgetChat) Handler(_ context.Context, conn net.Conn) {
	// TODO
}

func Run() {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	s := BudgetChat{cfg}

	internal.RunTCP(s)
}
