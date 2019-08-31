package ctcp

import "testing"

func TestEncodeMessage(t *testing.T) {
	t.Run("Test encode client version command", func(t *testing.T) {
		message, err := EncodeMessage(Message{"VERSION", CTCP, nil})
		if err != nil {
			t.Error(err)
		}
		expected := Message{"\001VERSION\001", PLAIN, nil}
		if message != expected {
			t.Fail()
		}
	})
	t.Run("Test encode newline escaping", func(t *testing.T) {
		message, err := EncodeMessage(Message{"Hi there!\nHow are you?", PLAIN, nil})
		if err != nil {
			t.Error(err)
		}
		expected := Message{"Hi there!\020nHow are you?", PLAIN, nil}
		if message != expected {
			t.Fail()
		}
	})
	t.Run("Test encode PLAIN followed by CTCP", func(t *testing.T) {
		message, err := EncodeMessage(Message{"Say hi to Ron\n\t/actor", PLAIN, &Message{"USERINFO", CTCP, nil}})
		if err != nil {
			t.Error(err)
		}
		expected := Message{"Say hi to Ron\020n\t/actor\001USERINFO\001", PLAIN, nil}
		if message != expected {
			t.Fail()
		}
	})
}

func TestDecodeMessage(t *testing.T) {
	t.Run("Test decode client version command", func(t *testing.T) {
		message, err := DecodeMessage(Message{"\001VERSION\001", PLAIN, nil})
		if err != nil {
			t.Error(err)
		}
		expected := Message{"VERSION", CTCP, nil}
		if message != expected {
			t.Fail()
		}
	})
}
