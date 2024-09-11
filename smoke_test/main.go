package smoketest

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/AVAniketh0905/protohackers/internal"
)

type SmokeTest struct{ *internal.Config }

func (SmokeTest) Setup() context.Context { return context.TODO() }

func (SmokeTest) Handler(_ context.Context, conn net.Conn) {
	defer conn.Close()

	// if _, err := io.Copy(conn, conn); err != nil {
	// 	log.Println(err)
	// }
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println(err)
				return
			}
			log.Fatal(err)
		}

		log.Println("read data successfully...", buf[:n])

		_, err = conn.Write(buf)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("write data successfully...")
	}
}

func Run() {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	s := SmokeTest{cfg}

	internal.RunTCP(s)
}
