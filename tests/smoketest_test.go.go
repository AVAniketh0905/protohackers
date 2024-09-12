package tests

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/AVAniketh0905/protohackers/internal"
)

func TestSmokeTest(t *testing.T) {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	// Start server in a seperate terminal

	client, err := net.Dial("tcp", cfg.Addr())
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println("client started...")

	msgs := []string{"Hello1", "Hello2", "Hello3", "Hello4", "Hello5"}

	for i := range 5 {
		_, err = client.Write([]byte(msgs[i]))
		if err != nil {
			t.Error(err)
		}

		buf := make([]byte, 1024)
		n, err := client.Read(buf)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(string(buf[:n]), msgs[i]) {
			t.Errorf("msg doesnt match, actual %v, got %v", msgs[i], string(buf[:n]))
		}
	}
}
