package tests

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/AVAniketh0905/protohackers/internal"
)

func TestUnusualDB(t *testing.T) {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	client, err := net.Dial("udp", cfg.Addr())
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	keys := []string{"hello", "golang", "network", "udp"}
	values := []string{"world", "gopher", "prog", "tcp"}

	for i, k := range keys {
		req := fmt.Sprintf("%v=%v", k, values[i])
		_, err := client.Write([]byte(req))
		if err != nil {
			t.Error(err)
		}
	}

	for i, k := range keys {
		_, err := client.Write([]byte(k))
		if err != nil {
			t.Error(err)
		}

		buf := make([]byte, 1024)
		n, err := client.Read(buf)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(fmt.Sprintf("%v=%v", k, values[i]), string(buf[:n])) {
			t.Errorf("the response does not match, expected %v, got %v", fmt.Sprintf("%v=%v", k, values[i]), string(buf[:n]))
		}
	}
}
