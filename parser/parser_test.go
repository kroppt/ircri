package parser

import (
	"fmt"
	"testing"
)

func TestParserBasic(t *testing.T) {
	t.Run("numeric command", func(t *testing.T) {
		p, out := NewParser(1)
		p.input = []rune("132\r\n")
		go p.Run()
		msg := <-out
		fmt.Println(msg)
	})
	t.Run("string command", func(t *testing.T) {
		p, out := NewParser(1)
		p.input = []rune("TESTING\r\n")
		go p.Run()
		msg := <-out
		fmt.Println(msg)
	})
}

func TestParserExamples(t *testing.T) {
	t.Run("param example 1", func(t *testing.T) {
		p, out := NewParser(1)
		p.input = []rune(":irc.example.com CAP * LIST :\r\n")
		go p.Run()
		msg := <-out
		fmt.Println(msg)
	})
}
