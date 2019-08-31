package msg

import "testing"

func TestEncode(t *testing.T) {
	t.Run("Test encode client version command", func(t *testing.T) {
		message := CTCP("VERSION")
		expected := "\001VERSION\001"
		value, err := message.Encode()
		if err != nil {
			t.Error(err)
		}
		if value != expected {
			t.Fail()
		}
	})
	t.Run("Test encode newline escaping", func(t *testing.T) {
		message := Plain("Hi there!\nHow are you?")
		expected := "Hi there!\020nHow are you?"
		value, err := message.Encode()
		if err != nil {
			t.Error(err)
		}
		if value != expected {
			t.Fail()
		}
	})
	t.Run("Test encode PLAIN followed by CTCP", func(t *testing.T) {
		messages := [2]Encoder{Plain("Say hi to Ron\n\t/actor"), CTCP("USERINFO")}
		expected := "Say hi to Ron\020n\t/actor\001USERINFO\001"
		value := ""
		for _, msg := range messages {
			val, err := msg.Encode()
			if err != nil {
				t.Error(err)
			}
			value += val
		}
		if value != expected {
			t.Fail()
		}
	})
}

func TestDecode(t *testing.T) {
	t.Run("Test decode client version command", func(t *testing.T) {
		message := "\001VERSION\001"
		expected := []Encoder{CTCP("VERSION")}
		values, err := Decode(message)
		if err != nil {
			t.Error(err)
		}
		if len(values) != len(expected) {
			t.Fail()
		}
		for i, exp := range expected {
			if values[i] != exp {
				t.Fail()
			}
		}
	})
}
