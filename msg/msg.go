package msg

const (
	mQuote = '\020'
	xQuote = '\134'
	nul    = '\000'
	nl     = '\n'
	cr     = '\r'
	xDelim = '\001'
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
	encoded := ""
	for _, char := range val {
		encoded += escapeXChar(byte(char))
	}
	// escape lower characters
	val = encoded
	encoded = ""
	for _, char := range val {
		encoded += escapeLowChar(byte(char))
	}
	return encoded, nil
}

// Encode turns the input CTCP into a valid CTCP escaped request or response string.
func (c CTCP) Encode() (string, error) {
	// escape CTCP characters
	val := c.Value()
	encoded := ""
	for _, char := range val {
		encoded += escapeXChar(byte(char))
	}
	// escape lower characters
	val = encoded
	encoded = ""
	for _, char := range val {
		encoded += escapeLowChar(byte(char))
	}
	return string(xDelim) + encoded + string(xDelim), nil
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
func Decode(message string) ([]Valuer, error) {
	decoded := ""
	results := []Valuer{}
	// low level dequoting
	for i, char := range message {
		if char == mQuote {
			if i != len(message)-1 {
				decoded += unescapeLowChar(message[i+1])
			}
		} else {
			decoded += string(char)
		}
	}
	decodedPart := ""
	low := true
	// CTCP level dequoting
	for i, char := range message {
		if char == xDelim {
			if low == true && len(decodedPart) != 0 {
				// split off first part into Plain
				results = append(results, Plain(decodedPart))
				low = false
			} else if low == true {
				// starts off with CTCP part
				low = false
			} else {
				// ends CTCP part
				results = append(results, CTCP(decodedPart))
				low = true
			}
			decodedPart = ""
		} else if char == xQuote {
			if i != len(message)-1 {
				decodedPart += unescapeXChar(message[i+1])
			}
		} else {
			decodedPart += string(char)
		}
	}
	if len(decodedPart) != 0 {
		results = append(results, Plain(decodedPart))
	}
	return results, nil
}

func escapeLowChar(char byte) string {
	switch char {
	case mQuote:
		return string(mQuote) + string(mQuote)
	case nul:
		return string(mQuote) + string('0')
	case nl:
		return string(mQuote) + string('n')
	case cr:
		return string(mQuote) + string('r')
	default:
		return string(char)
	}
}

func escapeXChar(char byte) string {
	switch char {
	case xDelim:
		return string(xQuote) + string('a')
	case xQuote:
		return string(xQuote) + string(xQuote)
	default:
		return string(char)
	}
}

func unescapeLowChar(char byte) string {
	switch char {
	case mQuote:
		return string(mQuote)
	case '0':
		return string(nul)
	case 'n':
		return string(nl)
	case 'r':
		return string(cr)
	default:
		return ""
	}
}

func unescapeXChar(char byte) string {
	switch char {
	case 'a':
		return string(xDelim)
	case xQuote:
		return string(xQuote)
	default:
		return ""
	}
}
