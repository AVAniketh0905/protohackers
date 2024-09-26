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
	"github.com/google/uuid"
)

type BudgetChat struct{ *internal.Config }

func (BudgetChat) Setup() context.Context { return context.TODO() }

func isAlphaNum(s string) bool {
	return regexp.MustCompile("^[a-zA-Z0-9_]*$").MatchString(s)
}

var connMap *sync.Map = &sync.Map{}
var userList []string

type User struct {
	conn net.Conn
	name string
}

type EndConnErr error

func (BudgetChat) Handler(_ context.Context, conn net.Conn) {
	var username string
	id := uuid.New().String()

	defer func() {
		conn.Close()
	}()

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
			if isAlphaNum(username) {
				log.Fatal("incorrect username")
			}
			userList = append(userList, username)
			user := User{conn: conn, name: username}
			connMap.Store(id, user)

			log.Printf("user %q connected at, %v\n.", username, conn.LocalAddr())

			connMap.Range(func(key, value interface{}) bool {
				user, ok := value.(User)
				if !ok {
					log.Fatal("improper user for, ", key)
				}

				if user.name == username {
					msg := fmt.Sprintf("* The room contains: %v\n", userList)
					_, err := user.conn.Write([]byte(msg))
					if err != nil {
						log.Fatalf("Error %v, for id %v", err, key)
					}
				} else {
					msg := fmt.Sprintf("* %v has entered the room\n", strings.Split(username, "\n")[0])
					_, err := user.conn.Write([]byte(msg))
					if err != nil {
						log.Fatalf("Error %v, for id %v", err, key)
					}
				}
				return true
			})
		} else {
			userInput, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Println("error:", err)
				}

				connMap.Delete(id)
				userList = []string{}
				connMap.Range(func(key, val any) bool {
					user, ok := val.(User)
					if !ok {
						log.Fatal("value from map is not User")
						return false
					}
					userList = append(userList, user.name)

					return true
				})
				return
			}

			if userInput == "\n" {
				log.Println("empty string")
				continue
			}

			log.Printf("user at %v sent %q\n.", conn.LocalAddr(), userInput)

			connMap.Range(func(key, value interface{}) bool {
				user, ok := value.(User)
				if !ok {
					log.Fatal("improper user for, ", key)
				}

				log.Printf("Received from %v: %s", username, userInput)

				msg := fmt.Sprintf("[%v] %v", username[:len(username)-1], userInput)
				_, err := user.conn.Write([]byte(msg))
				if err != nil {
					log.Fatalf("Error %v, for id %v", err, key)
				}

				return true
			})
		}
	}
}

// func Run() {
// 	cfg := internal.NewConfig(internal.PORT)
// 	cfg.ParseFlags()
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)

// 	s := BudgetChat{cfg}

// 	internal.RunTCP(s)
// }
