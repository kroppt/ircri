package ctcp

import "errors"

// Message contains text along with the text's type
type Message struct {
	value    string
	textType TextType
}

// TextType is used to specify whether an input is literal text or CTCP text
type TextType int

const (
	// PLAIN represents literal text
	PLAIN TextType = iota
	// CTCP represents a CTCP request or response
	CTCP
)

// EncodeMessage turns the input string into a valid CTCP escaped message.
func EncodeMessage(message Message) (Message, error) {
	return Message{"", PLAIN}, errors.New("not implemented")
}

// EncodeChainMessage turns the input string into a valid CTCP escaped message.
func EncodeChainMessage(preMessage Message, postMessage Message) (Message, error) {
	return Message{"", PLAIN}, errors.New("not implemented")
}

// DecodeMessage turns the input string from a CTCP escaped message into plain text.
func DecodeMessage(message Message) (Message, error) {
	return Message{"", PLAIN}, errors.New("not implemented")
}
