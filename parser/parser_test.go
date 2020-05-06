package parser

import (
	"reflect"
	"testing"
	"time"
)

type basicExpect struct {
	name   string
	input  string
	expect Message
}

func testParserExpect(t *testing.T, tests []basicExpect) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, out := NewParser(1)
			p.input = []rune(test.input)
			go p.Run()
			select {
			case msg, ok := <-out:
				if !ok {
					t.Error("message channel closed unexpectedly\n")
				} else if !reflect.DeepEqual(msg, test.expect) {
					t.Errorf("expected %v to equal %v\n", msg, test.expect)
				}
			case <-time.After(1 * time.Second):
				t.Error("timed out after 1 second\n")
			}
		})
	}
}

func TestParserBasic(t *testing.T) {
	tests := []basicExpect{
		{"numeric command", "132\r\n", Message{Command: "132"}},
		{"string command", "TESTING\r\n", Message{Command: "TESTING"}},
	}
	testParserExpect(t, tests)
}

func TestParserExamples(t *testing.T) {
	paramEx1Msg := Message{
		Prefix:  Prefix{Name: "irc.example.com"},
		Command: "CAP",
		Params:  []string{"*", "LIST", ""},
	}
	paramEx2Msg := Message{
		Command: "CAP",
		Params:  []string{"*", "LS", "multi-prefix sasl"},
	}
	paramEx3Msg := Message{
		Command: "CAP",
		Params:  []string{"REQ", "sasl message-tags foo"},
	}
	paramEx4Msg := Message{
		Prefix: Prefix{
			Name:     "dan",
			Username: "d",
			Host:     "localhost",
		},
		Command: "PRIVMSG",
		Params:  []string{"#chan", "Hey!"},
	}
	completeEx1Msg := Message{
		Prefix:  Prefix{Name: "irc.example.com"},
		Command: "CAP",
		Params:  []string{"LS", "*", "multi-prefix extended-join sasl"},
	}
	completeEx2Msg := Message{
		Tags: []Tag{{
			Key:   "id",
			Value: "234AB",
		}},
		Prefix: Prefix{
			Name:     "dan",
			Username: "d",
			Host:     "localhost",
		},
		Command: "PRIVMSG",
		Params:  []string{"#chan", "Hey what's up!"},
	}
	completeEx3Msg := Message{
		Command: "CAP",
		Params:  []string{"REQ", "sasl"},
	}
	tests := []basicExpect{
		{"param example 1", ":irc.example.com CAP * LIST :\r\n", paramEx1Msg},
		{"param example 2", "CAP * LS :multi-prefix sasl\r\n", paramEx2Msg},
		{"param example 3", "CAP REQ :sasl message-tags foo\r\n", paramEx3Msg},
		{"param example 4", ":dan!d@localhost PRIVMSG #chan :Hey!\r\n", paramEx4Msg},
		{"param example 5", ":dan!d@localhost PRIVMSG #chan Hey!\r\n", paramEx4Msg},
		{"complete example 1", ":irc.example.com CAP LS * :multi-prefix extended-join sasl\r\n", completeEx1Msg},
		{"complete example 2", "@id=234AB :dan!d@localhost PRIVMSG #chan :Hey what's up!\r\n", completeEx2Msg},
		{"complete example 3", "CAP REQ :sasl\r\n", completeEx3Msg},
	}
	testParserExpect(t, tests)
}
