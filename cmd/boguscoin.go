package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/AVAniketh0905/protohackers/internal"
)

type BogusCoin struct{ *internal.Config }

type chatServer struct {
	server net.Conn
}

type chatStr string

const bogusTonyAddr string = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"

func NewChatServer(addr string, port int32) *chatServer {
	addr = fmt.Sprintf("%v:%v", addr, port)
	server, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &chatServer{
		server: server,
	}
}

func (cs *chatServer) Read(p []byte) (n int, err error) {
	n, err = cs.server.Read(p)
	return
}

func (cs *chatServer) Write(p []byte) (n int, err error) {
	n, err = cs.server.Write(p)
	return
}

func (b BogusCoin) Setup() context.Context {
	key := chatStr("srv")
	server := NewChatServer("chat.protohackers.com", 16963)
	ctx := context.WithValue(context.TODO(), key, server)
	return ctx
}

func (b BogusCoin) Handler(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	parentAny := ctx.Value(chatStr("srv"))
	parentSrv, ok := parentAny.(chatServer)
	if !ok {
		log.Fatal("incorrect values in parent server")
	}

	for {
		reqStr, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// handle malicious logic
		reqStr = strings.ReplaceAll(reqStr, "\n", "")
		strSplit := strings.Split(reqStr, " ")

		for i, v := range strSplit {
			if string(v[0]) != "7" {
				continue
			}

			// TODO: other checks

			strSplit[i] = bogusTonyAddr
		}

		reqStr = strings.Join(strSplit, " ") + "\n"
		_, err = parentSrv.Write([]byte(reqStr))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Run() {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	s := BogusCoin{cfg}

	internal.RunTCP(s)
}
