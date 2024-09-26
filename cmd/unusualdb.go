package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/AVAniketh0905/protohackers/internal"
)

type UnusualDB struct{ *internal.Config }

var dataMap *sync.Map = &sync.Map{}

func (u UnusualDB) Setup() context.Context { return context.TODO() }

func (u UnusualDB) Handler(_ context.Context, conn net.PacketConn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFrom(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(clientAddr, err)
			}
			break
		}

		strBuf := string(buf[:n])
		log.Println(strBuf)

		if strings.Contains(strBuf, "=") {
			//insert
			key, value, _ := strings.Cut(strBuf, "=")
			dataMap.Store(key, value)
		} else {
			// retrieves
			var res string
			val, ok := dataMap.Load(strBuf)
			if !ok {
				res = fmt.Sprintf("%v=", strBuf)
			} else {
				res = fmt.Sprintf("%v=%v", strBuf, val)
			}

			log.Println("response", res)

			_, err := conn.WriteTo([]byte(res), clientAddr)
			if err != nil {
				log.Fatal(clientAddr, err)
			}
		}
	}
}

func Run() {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	s := UnusualDB{cfg}

	internal.RunUDP(s)
}
