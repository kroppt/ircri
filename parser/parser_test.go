package parser

import (
	"reflect"
	"testing"
	"time"
)

// when debugging, set this timeout to like an hour or something
const timeout = 1 * time.Second

type basicExpect struct {
	name   string
	input  string
	expect Message
}

func testParserExpect(t *testing.T, tests []basicExpect) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in, out := make(chan []rune, 1), make(chan Message, 1)
			p := NewParser(in, out)
			in <- []rune(test.input)
			cancel := make(chan struct{})
			done := make(chan struct{})
			go func() {
				p.Run(cancel)
				done <- struct{}{}
			}()
			defer func() {
				close(cancel)
				close(done)
				close(in)
				close(out)
			}()
			select {
			case <-done:
			case <-time.After(timeout):
				t.Error("timed out after 1 second\n")
				return
			}
			select {
			case msg, ok := <-out:
				if !ok {
					t.Error("message channel closed unexpectedly\n")
				} else if !reflect.DeepEqual(msg, test.expect) {
					t.Errorf("expected %v to equal %v\n", msg, test.expect)
				}
			default:
				t.Error("expected parsed message but got failure\n")
			}
		})
	}
}

type failExpect struct {
	name  string
	input string
}

func testParserFails(t *testing.T, tests []failExpect) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in, out := make(chan []rune, 1), make(chan Message, 1)
			p := NewParser(in, out)
			in <- []rune(test.input)
			cancel := make(chan struct{})
			done := make(chan struct{})
			go func() {
				p.Run(cancel)
				done <- struct{}{}
			}()
			defer func() {
				close(cancel)
				close(done)
				close(in)
				close(out)
			}()
			select {
			case <-done:
			case <-time.After(timeout):
				t.Error("timed out after 1 second\n")
				return
			}
			select {
			case msg, ok := <-out:
				if ok {
					t.Errorf("expected parse failure but got %v\n", msg)
				}
			default:
			}
		})
	}
}

func TestParserBasic(t *testing.T) {
	singleTagMsg := Message{
		Tags: []Tag{
			{Key: "verb"},
		},
		Command: "TEST",
	}
	manySpacesParamsMsg := Message{
		Command: "TEST",
		Params:  []string{"abc"},
	}
	manySpacesTagsMsg := Message{
		Tags: []Tag{
			{Key: "verb"},
		},
		Command: "TEST",
	}
	manySpacesPrefixMsg := Message{
		Tags: []Tag{
			{Key: "verb"},
		},
		Prefix:  Prefix{Name: "abc.com"},
		Command: "TEST",
	}
	tests := []basicExpect{
		{"numeric command", "132\r\n", Message{Command: "132"}},
		{"string command", "TESTING\r\n", Message{Command: "TESTING"}},
		{"single tag", "@verb TEST\r\n", singleTagMsg},
		{"many spaces params", "TEST   abc\r\n", manySpacesParamsMsg},
		{"many spaces tags", "@verb   TEST\r\n", manySpacesTagsMsg},
		{"many spaces prefix", "@verb  :abc.com  TEST\r\n", manySpacesPrefixMsg},
	}
	testParserExpect(t, tests)
}

func TestParserExamples(t *testing.T) {
	tagEx1Msg := Message{
		Tags: []Tag{
			{Key: "id", Value: "123AB"},
			{Key: "rose"},
		},
		Command: "CAP",
	}
	tagEx2Msg := Message{
		Tags: []Tag{
			{Key: "url"},
			{Key: "netsplit", Value: "tur,ty"},
		},
		Command: "CAP",
	}
	tagEx3Msg := Message{
		Tags: []Tag{
			{Vendor: "localhost", Key: "verb"},
		},
		Command: "CAP",
	}
	tagEx4Msg := Message{
		Tags: []Tag{
			{Vendor: "localhost", Key: "id", Value: "123AB"},
		},
		Command: "CAP",
	}
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
		Prefix:  Prefix{Name: "dan", Username: "d", Host: "localhost"},
		Command: "PRIVMSG",
		Params:  []string{"#chan", "Hey!"},
	}
	completeEx1Msg := Message{
		Prefix:  Prefix{Name: "irc.example.com"},
		Command: "CAP",
		Params:  []string{"LS", "*", "multi-prefix extended-join sasl"},
	}
	completeEx2Msg := Message{
		Tags: []Tag{
			{Key: "id", Value: "234AB"},
		},
		Prefix:  Prefix{Name: "dan", Username: "d", Host: "localhost"},
		Command: "PRIVMSG",
		Params:  []string{"#chan", "Hey what's up!"},
	}
	completeEx3Msg := Message{
		Command: "CAP",
		Params:  []string{"REQ", "sasl"},
	}
	completeEx4Msg := Message{
		Tags: []Tag{
			{Vendor: "address1", Key: "k1", Value: "v1"},
			{Vendor: "address2", Key: "k2", Value: "v2"},
			{Key: "k3", Value: "v3"},
			{Key: "k4"},
			{Key: "k5"},
		},
		Prefix:  Prefix{Name: "full", Username: "nick", Host: "address"},
		Command: "CMD",
		Params:  []string{"param1", "param2", "spaced param"},
	}
	tests := []basicExpect{
		{"tag example 1", "@id=123AB;rose CAP\r\n", tagEx1Msg},
		{"tag example 2", "@url=;netsplit=tur,ty CAP\r\n", tagEx2Msg},
		{"tag example 3", "@localhost/verb CAP\r\n", tagEx3Msg},
		{"tag example 4", "@localhost/id=123AB CAP\r\n", tagEx4Msg},
		{"param example 1", ":irc.example.com CAP * LIST :\r\n", paramEx1Msg},
		{"param example 2", "CAP * LS :multi-prefix sasl\r\n", paramEx2Msg},
		{"param example 3", "CAP REQ :sasl message-tags foo\r\n", paramEx3Msg},
		{"param example 4", ":dan!d@localhost PRIVMSG #chan :Hey!\r\n", paramEx4Msg},
		{"param example 5", ":dan!d@localhost PRIVMSG #chan Hey!\r\n", paramEx4Msg},
		{"complete example 1", ":irc.example.com CAP LS * :multi-prefix extended-join sasl\r\n", completeEx1Msg},
		{"complete example 2", "@id=234AB :dan!d@localhost PRIVMSG #chan :Hey what's up!\r\n", completeEx2Msg},
		{"complete example 3", "CAP REQ :sasl\r\n", completeEx3Msg},
		{"complete example 4", "@address1/k1=v1;address2/k2=v2;k3=v3;k4=;k5 :full!nick@address CMD param1 param2 :spaced param\r\n", completeEx4Msg},
	}
	testParserExpect(t, tests)
}

func TestParserFailures(t *testing.T) {
	tests := []failExpect{
		{"empty input", ""},
		{"empty message", "\r\n"},
		{"short numeric command", "12\r\n"},
		{"long numeric command", "1234\r\n"},
		{"number-letter command", "12A\r\n"},
		{"extra tag delim", "@id=123AB; CAP\r\n"},
		{"end after tags", "@id=123AB\r\n"},
		{"end after prefix", ":irc.example.com\r\n"},
		{"trailing space", "CAP \r\n"},
	}
	testParserFails(t, tests)
}
