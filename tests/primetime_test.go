package tests

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/AVAniketh0905/protohackers/cmd"
	"github.com/AVAniketh0905/protohackers/internal"
)

func TestPrimeTime(t *testing.T) {
	cfg := internal.NewConfig(internal.PORT)
	cfg.ParseFlags()

	client, err := net.Dial("tcp", cfg.Addr())
	if err != nil {
		t.Error(err)
	}
	defer client.Close()
	fmt.Println("client started...")

	json_reqs := []string{"{\"method\":\"isPrime\",\"number\":12}", "{\"method\":\"isPrime\",\"number\":546}", "{\"method\":\"isPrime\",\"number\":23}"}
	json_answers := []bool{false, false, true}

	for i, req := range json_reqs {
		fmt.Println("sending req...", req)

		_, err = client.Write([]byte(req))
		if err != nil {
			t.Error(err)
		}

		fmt.Println("finished writing to client...")

		buf := make([]byte, 1024)
		n, err := client.Read(buf)
		if err != nil {
			t.Error(err)
		}
		fmt.Println("received from server...", string(buf[:n]))

		var resp cmd.Resp
		err = json.Unmarshal(buf[:n], &resp)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(resp.Prime, json_answers[i]) {
			t.Errorf("answers did not match, actual %v, got %v", json_answers[i], resp.Prime)
		}
	}
}
