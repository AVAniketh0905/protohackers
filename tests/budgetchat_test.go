package tests

import (
	"bufio"
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/AVAniketh0905/protohackers/internal"
)

func TestConnStart(t *testing.T) {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	client, err := net.Dial("tcp", cfg.Addr())
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println("client started...")

	buf := []byte("hello im good\n")
	_, err = client.Write(buf)
	if err != nil {
		t.Error(err)
	}

	buf = make([]byte, 1024)
	n, err := client.Read(buf)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(string(buf[:n]), "[bob] hello im good\n") {
		t.Errorf("messages did not match, expected, %q but got, %q", "[bob] hello im good\n", string(buf[:n]))
	}
}

func simulateClient(t *testing.T, wg *sync.WaitGroup, cfg *internal.Config, id int, messages []string) {
	defer wg.Done()

	var conn net.Conn
	var err error
	for attempts := 0; attempts < 3; attempts++ {
		conn, err = net.Dial("tcp", cfg.Addr())
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond) // Retry delay
	}
	if err != nil {
		t.Errorf("Client %d failed to connect after retries: %v", id, err)
		return
	}
	defer conn.Close()

	for _, msg := range messages {
		_, err := fmt.Fprintf(conn, msg+"\n")
		if err != nil {
			t.Error(err)
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.Errorf("Client %d failed to read: %v", id, err)
			return
		}
		t.Logf("Client %d sent: %s, received: %s", id, strings.TrimSpace(msg), strings.TrimSpace(response))
	}
}

func TestMultipleConns(t *testing.T) {
	var wg sync.WaitGroup

	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	numClients := 2
	messages := [][]string{
		{"Hello from Client 1", "How are you? 1", "Bye! 1"},
		{"Hi from Client 2", "What's up? 2", "Hello! 2"},
	}

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go simulateClient(t, &wg, cfg, i+1, messages[i])
	}

	wg.Wait()
}
