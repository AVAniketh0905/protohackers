package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"

	"github.com/AVAniketh0905/protohackers/internal"
)

type BogusCoin struct{ *internal.Config }

const bogusTonyAddr string = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"

var boguscoin = regexp.MustCompile(`^7[a-zA-Z0-9]{25,34}$`)

func (b BogusCoin) Setup() context.Context {
	return context.TODO()
}

func (b BogusCoin) Handler(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	addr := fmt.Sprintf("%v:%v", "chat.protohackers.com", 16963)
	server, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	once := sync.Once{}
	proxy := func(dst io.WriteCloser, src io.ReadCloser) {
		defer once.Do(func() {
			dst.Close()
			src.Close()
			log.Printf("closed connection: %v", addr)
		})

		for r := bufio.NewReader(src); ; {
			msg, err := r.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Fatal(err)
				}
				return
			}

			tokens := make([]string, 0, 8)
			for _, raw := range strings.Split(msg[:len(msg)-1], " ") {
				t := boguscoin.ReplaceAllString(raw, bogusTonyAddr)
				tokens = append(tokens, t)
			}

			newMsg := strings.Join(tokens, " ") + "\n"
			if _, err := dst.Write([]byte(newMsg)); err != nil {
				log.Fatal(err)
			}
		}
	}

	go proxy(conn, server)
	proxy(server, conn)
}

// func Run() {
// 	cfg := internal.NewConfig(internal.PORT)
// 	cfg.ParseFlags()

// 	s := BogusCoin{cfg}

// 	internal.RunTCP(s)
// }
