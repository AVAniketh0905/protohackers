package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func receive(client net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		msg, err := bufio.NewReader(client).ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("Server closed the connection.")
			}
			break
		}
		log.Println("Message from server:", msg)
	}
}

func write(client net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	reader := bufio.NewReader(os.Stdin)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Fatal("Error reading input:", err)
			}
			break
		}

		_, err = client.Write([]byte(msg))
		if err != nil {
			log.Fatal("Error sending message:", err)
		}
	}
}

func Run() {
	client, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	buf := make([]byte, 1024)
	n, err := client.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
		return
	}

	log.Println(string(buf[:n]))

	var username string
	fmt.Scan(&username)

	username = fmt.Sprintf("%v\n", username)
	_, err = client.Write([]byte(username))
	if err != nil {
		log.Fatal(err)
	}

	buf = make([]byte, 1024)
	n, err = client.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
		return
	}

	log.Println(string(buf[:n]))

	var wg sync.WaitGroup

	wg.Add(1)
	go receive(client, &wg)

	wg.Add(1)
	go write(client, &wg)

	wg.Wait()
}
