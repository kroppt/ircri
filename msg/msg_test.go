package msg

import "testing"

func TestEncode(t *testing.T) {
	t.Run("Test encode client version command", func(t *testing.T) {
		message := CTCP("VERSION")
		expected := "\001VERSION\001"
		value, err := message.Encode()
		if err != nil {
			t.Error(err)
			return
		}
		if value != expected {
			t.Error("For", message, "expected", expected, "got", value)
		}
	})
	t.Run("Test encode newline escaping", func(t *testing.T) {
		message := Plain("Hi there!\nHow are you?")
		expected := "Hi there!\020nHow are you?"
		value, err := message.Encode()
		if err != nil {
			t.Error(err)
			return
		}
		if value != expected {
			t.Error("For", message, "expected", expected, "got", value)
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
				return
			}
			value += val
		}
		if value != expected {
			t.Error("For", messages, "expected", expected, "got", value)
		}
	})
	t.Run("Test encode SED encrypted text", func(t *testing.T) {
		message := CTCP("SED \n\t\big\020\001\000\\:")
		expected := "\001SED \020n\t\big\020\020\\a\0200\\\\:\001"
		value, err := message.Encode()
		if err != nil {
			t.Error(err)
			return
		}
		if value != expected {
			t.Error("For", []rune(message), "expected", []rune(expected), "got", []rune(value))
		}
	})
}

func TestDecode(t *testing.T) {
	t.Run("Test decode client version command", func(t *testing.T) {
		message := "\001VERSION\001"
		expected := []Valuer{CTCP("VERSION")}
		values, err := Decode(message)
		if err != nil {
			t.Error(err)
			return
		}
		if len(values) != len(expected) {
			t.Error("For", message, "expected", expected, "got", values)
			return
		}
		for i, exp := range expected {
			if i >= len(values) || values[i].Value() != exp.Value() {
				t.Error("For", message, "expected", expected, "got", values)
			}
		}
	})
}
