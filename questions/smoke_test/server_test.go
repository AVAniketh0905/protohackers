package smoketest

import (
	"fmt"
	"net"
	"testing"
	"time"

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

	for i := range 7 {
		time.Sleep(time.Duration((i + 1)) * time.Second)
		msg := []byte(fmt.Sprintf("Hello, %d ", i))
		_, err = client.Write(msg)
		if err != nil {
			t.Error(err)
		}
	}
}
