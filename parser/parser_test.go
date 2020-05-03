package parser

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
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
