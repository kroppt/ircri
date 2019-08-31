package ctcp

import "testing"

func TestEncodeMessage(t *testing.T) {
	t.Run("Test encode client version command", func(t *testing.T) {
		message, err := EncodeMessage("VERSION")
		if err != nil {
			t.Error(err)
		}
		if message != "\x01VERSION\x01" {
			t.Fail()
		}
	})
}

func TestDecodeMessage(t *testing.T) {
	t.Run("Test decode client version command", func(t *testing.T) {
		message, err := DecodeMessage("\x01VERSION\x01")
		if err != nil {
			t.Error(err)
		}
		if message != "VERSION" {
			t.Fail()
		}
	})
}
