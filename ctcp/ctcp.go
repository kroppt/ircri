package ctcp

import "errors"

// Message contains text, the text's type, and the next part of the message, if there is one
type Message struct {
	value       string
	textType    TextType
	nextMessage *Message
}

// TextType is used to specify whether an input is literal text or CTCP text
type TextType int

const (
	// PLAIN represents literal text
	PLAIN TextType = iota
	// CTCP represents a CTCP request or response
	CTCP
)

// EncodeMessage turns the input message into a valid CTCP escaped message.
func EncodeMessage(message Message) (Message, error) {
	return Message{"", PLAIN, nil}, errors.New("not implemented")
}

// DecodeMessage turns the input message from a CTCP escaped message into plain text.
func DecodeMessage(message Message) (Message, error) {
	return Message{"", PLAIN, nil}, errors.New("not implemented")
}
