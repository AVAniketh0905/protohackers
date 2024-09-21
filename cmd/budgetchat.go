package cmd

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"regexp"
	"sync"

	"github.com/AVAniketh0905/protohackers/internal"
	"github.com/google/uuid"
)

type BudgetChat struct{ *internal.Config }

func (BudgetChat) Setup() context.Context { return context.TODO() }

func isAlphaNum(s string) bool {
	return regexp.MustCompile("^[a-zA-Z0-9_]*$").MatchString(s)
}

var connMap *sync.Map = &sync.Map{}

func (BudgetChat) Handler(_ context.Context, conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	var username string
	id := uuid.New().String()
	connMap.Store(id, conn)

	_, err := conn.Write([]byte("Welcome to budgetchat! What shall I call you?"))
	if err != nil {
		log.Fatal(err)
	}

	for {
		if username == "" {
			username, err = bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Fatal(err)
				}
				break
			}

			log.Printf("user %q connected at, %v\n.", username, conn.LocalAddr())
		} else {
			userInput, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Fatal(err)
				}
				break
			}

			log.Printf("user at %v sent %q\n.", conn.LocalAddr(), userInput)

			connMap.Range(func(key, value interface{}) bool {
				conn, ok := value.(net.Conn)
				if !ok {
					log.Fatal("imporper value for, ", key)
				}

				log.Printf("Received from %v: %s", conn.RemoteAddr(), userInput)

				msg := "[" + username[:len(username)-1] + "] " + userInput
				_, err := conn.Write([]byte(msg))
				if err != nil {
					log.Fatalf("Error %v, for id %v", err, key)
				}

				return true
			})
		}
	}
}

func Run() {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	s := BudgetChat{cfg}

	internal.RunTCP(s)
}
