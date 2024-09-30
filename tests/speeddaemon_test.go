package tests

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/AVAniketh0905/protohackers/cmd"
)

func TestStr(t *testing.T) {
	hexMsgs := []string{"00", "03666f6f", "08456C626572657468"}
	msg := []string{"", "foo", "Elbereth"}

	for i, h := range hexMsgs {
		bytes, err := hex.DecodeString(h)
		if err != nil {
			t.Error(i, h, err)
		}

		var dataStr cmd.Str
		dataStr.Unmarshall(bytes)

		if len(msg[i]) != int(dataStr.Len) {
			t.Error(i, h, "len does not match")
		}

		if !reflect.DeepEqual(msg[i], string(dataStr.Msg)) {
			t.Errorf("i: %d h: %q, msgs doesnt match, expected %q, got %q", i, h, msg[i], string(dataStr.Msg))
		}
	}

}
