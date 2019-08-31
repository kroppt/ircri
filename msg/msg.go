package msg

import "errors"

// Encoder has an Encode function which allows it to be converted into a string.
type Encoder interface {
	Encode() (string, error)
}

// Valuer has a Value function which allows the underlying string value it to be unwrapped from its type.
type Valuer interface {
	Value() string
}

// Plain represents a plain text string.
type Plain string

// CTCP represents a CTCP string.
type CTCP string

// Encode turns the input Plain into a valid escaped string.
func (p Plain) Encode() (string, error) {
	return "", errors.New("not implemented")
}

// Encode turns the input CTCP into a valid CTCP escaped request or response string.
func (c CTCP) Encode() (string, error) {
	return "", errors.New("not implemented")
}

// Value unwraps the type to expose the underlying string.
func (p Plain) Value() string {
	return string(p)
}

// Value unwraps the type to expose the underlying string.
func (c CTCP) Value() string {
	return string(c)
}

// Decode turns the input message from a CTCP escaped message into plain text.
func Decode(string) ([]Encoder, error) {
	return nil, errors.New("not implemented")
}
