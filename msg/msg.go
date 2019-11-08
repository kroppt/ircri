package msg

import (
	"bytes"
	"errors"
)

const (
	mQuote byte = '\020'
	xQuote byte = '\134'
	nul    byte = '\000'
	nl     byte = '\n'
	cr     byte = '\r'
	xDelim byte = '\001'
	space  byte = '\040'
)

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
	// escape CTCP characters
	val := p.Value()
	encoded := new(bytes.Buffer)
	for _, char := range []byte(val) {
		encoded.Write(escapeXChar(char))
	}
	// escape lower characters
	val = encoded.String()
	encoded.Reset()
	for _, char := range []byte(val) {
		encoded.Write(escapeLowChar(char))
	}
	return encoded.String(), nil
}

// Encode turns the input CTCP into a valid CTCP escaped request or response string.
func (c CTCP) Encode() (string, error) {
	// escape CTCP characters
	val := c.Value()
	if len(val) > 0 && isWhitespace(val[0]) {
		return "", errors.New("expected valid character at position 0")
	}
	encoded := new(bytes.Buffer)
	for _, char := range []byte(val) {
		encoded.Write(escapeXChar(char))
	}
	// escape lower characters
	val = encoded.String()
	encoded.Reset()
	for _, char := range []byte(val) {
		encoded.Write(escapeLowChar(char))
	}
	return string(xDelim) + encoded.String() + string(xDelim), nil
}

// Value unwraps the type to expose the underlying string.
func (p Plain) Value() string {
	return string(p)
}

// Value unwraps the type to expose the underlying string.
func (c CTCP) Value() string {
	return string(c)
}

// Decode turns a CTCP escaped string message into a slice of the components of the message.
func Decode(message string) ([]Valuer, error) {
	decoded := new(bytes.Buffer)
	results := []Valuer{}
	// low level dequoting
	for i, char := range []byte(message) {
		if char == mQuote {
			if i != len(message)-1 {
				val, ok := unescapeLowChar(message[i+1])
				if ok {
					decoded.WriteByte(val)
				}
			}
		} else {
			decoded.WriteByte(char)
		}
	}
	val := decoded.String()
	decoded.Reset()
	low := true
	// CTCP level dequoting
	for i, char := range []byte(val) {
		if char == xDelim {
			if low == true && decoded.Len() != 0 {
				// split off first part into Plain
				results = append(results, Plain(decoded.String()))
				low = false
			} else if low == true {
				// starts off with CTCP part
				low = false
			} else {
				// ends CTCP part
				results = append(results, CTCP(decoded.String()))
				low = true
			}
			decoded.Reset()
		} else if char == xQuote {
			if i != len(message)-1 {
				val, ok := unescapeXChar(message[i+1])
				if ok {
					decoded.WriteByte(val)
				}
			}
		} else {
			decoded.WriteByte(char)
		}
	}
	if decoded.Len() != 0 {
		results = append(results, Plain(decoded.String()))
	}
	return results, nil
}

func escapeLowChar(char byte) []byte {
	switch char {
	case mQuote:
		return []byte{mQuote, mQuote}
	case nul:
		return []byte{mQuote, '0'}
	case nl:
		return []byte{mQuote, 'n'}
	case cr:
		return []byte{mQuote, 'r'}
	default:
		return []byte{char}
	}
}

func escapeXChar(char byte) []byte {
	switch char {
	case xDelim:
		return []byte{xQuote, 'a'}
	case xQuote:
		return []byte{xQuote, xQuote}
	default:
		return []byte{char}
	}
}

func unescapeLowChar(char byte) (byte, bool) {
	switch char {
	case mQuote:
		return mQuote, true
	case '0':
		return nul, true
	case 'n':
		return nl, true
	case 'r':
		return cr, true
	default:
		return 0, false
	}
}

func unescapeXChar(char byte) (byte, bool) {
	switch char {
	case 'a':
		return xDelim, true
	case xQuote:
		return xQuote, true
	default:
		return 0, false
	}
}

func isWhitespace(char byte) bool {
	switch char {
	case xDelim:
		return true
	case space:
		return true
	default:
		return false
	}
}
