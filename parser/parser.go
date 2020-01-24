package parser

import (
	"strings"
)

// Message represents a parsed IRC message.
type Message struct {
	Tags    []Tag
	Prefix  Prefix
	Command Command
	Params  []string
}

// Tag is a key setting which optionally prepends IRC messages.
//
// A tag may optionally have a value associated with the key.
// A tag may optionally have a vendor associated with the key.
type Tag struct {
	Vendor string
	Key    string
	Value  string
}

// Prefix is additional information about a message which optionally prepends
// IRC messages after the optional tags.
//
// A prefix requires either a server name or a nickname.
// Optionally, when specifying the nickname, one can specify the host and the
// username.
type Prefix struct {
	Server   string
	Nickname string
	Username string
	Host     string
}

// Command is required for any message.
//
// A command name is given for regular commands, typically sent by the
// client to the server.
//
// A numeric command (aka numeric reply) is typically sent by the server to the
// client as a reply to a previous command.
type Command struct {
	Name    string
	Numeric string
}

// Parser holds state information necessary for parsing IRC messages.
type Parser struct {
	// message currently being built
	msg    Message
	output chan<- Message
	input  []rune
	pos    int
}

// NewParser returns a parser with an output buffer of the given size.
//
// The parser is responsible for creating and closing its output channel.
func NewParser(size int) (*Parser, <-chan Message) {
	chout := make(chan Message, size)
	return &Parser{
		output: chout,
		pos:    0,
	}, chout
}

// Run begins the main execution loop for parsing.
//
// The output is passed through the Parser's output channel.
func (p *Parser) Run() {
	var state StateFn
	for state = beginState; state != nil; {
		state = state(p)
	}
}

// Next returns the next rune to be looked at and a boolean for if a rune exists.
func (p *Parser) Next() (rune, bool) {
	if p.pos > len(p.input) {
		return 0, false
	}
	r := p.input[p.pos]
	p.pos++
	return r, true
}

// Rewind moves the position index back one rune.
func (p *Parser) Rewind() {
	p.pos--
}

// Consume deletes the runes before position index.
func (p *Parser) Consume() string {
	out := string(p.input[:p.pos])
	p.input = p.input[p.pos:]
	p.pos = 0
	return out
}

// StateFn returns the next state function to run.
type StateFn func(*Parser) StateFn

// beginState is the entry point to the IRC message parsing state machine.
func beginState(p *Parser) StateFn {
	p.msg = Message{}
	r, ok := p.Next()
	if !ok {
		return nil
	}
	if r == '@' {
		p.Consume()
		return tagState
	}
	if r == ':' {
		p.Consume()
		return prefixState
	}
	if (r > 'a' && r < 'z') || (r > 'A' && r < 'Z') || (r > '0' && r < '9') {
		// the rune is part of a command
		p.Rewind()
		return commandState
	}
	p.Consume()
	return nil
}

func isHostnameRune(r rune) bool {
	switch true {
	case '0' <= r && r <= '9':
	case 'A' <= r && r <= 'Z':
	case 'a' <= r && r <= 'z':
	case r == '.':
	case r == '-':
	default:
		return false
	}
	return true
}

func isValueRune(r rune) bool {
	switch r {
	case '\x00': // NUL
	case '\x07': // BEL
	case '\x0D': // CR
	case '\x0A': // LF
	case ';':
	case ' ':
	default:
		return true
	}
	return false
}

// parseUntil parses runes until the given predicate fails for one of the runes.
func parseUntil(p *Parser, pred func(rune) bool) string {
	var sb strings.Builder
	r, ok := p.Next()
	if !ok {
		return sb.String()
	}
	for pred(r) {
		sb.WriteRune(r)
		r, ok = p.Next()
		if !ok {
			return sb.String()
		}
	}
	p.Rewind()
	return sb.String()
}

func tagState(p *Parser) StateFn {
	var newtag Tag
	var key string
	var vendor string
	// parse tag key
	// - parse vendor hostname optional string
	// - parse any number of alpha, digit, '.', and '-' runes
	key = parseUntil(p, isHostnameRune)
	r, ok := p.Next()
	if r == '/' {
		vendor = key
		key = parseUntil(p, isHostnameRune)
		if len(key) == 0 {
			// TODO handle error
			return nil
		}
	}
	if strings.ContainsRune(key, '.') {
		// TODO handle error
		return nil
	}
	if vendor[0] == '.' || vendor[0] == '-' {
		// TODO handle error
		return nil
	}
	if vendor[len(vendor)-1] == '.' {
		vendor = vendor[:len(vendor)-1]
	}
	if len(vendor) > 253 {
		// TODO handle error
		return nil
	}
	for _, lbl := range strings.Split(vendor, ".") {
		if len(lbl) < 0 || len(lbl) > 63 {
			// TODO handle error
			return nil
		}
	}
	newtag.Key = key
	newtag.Vendor = vendor

	// parse optional tag value
	var value string
	if r == '=' {
		// parse any number of runes except NUL, BELL, CR, LF, ';', ' '
		value = parseUntil(p, isValueRune)
		r, ok = p.Next()
		if !ok {
			// TODO handle error
			return nil
		}
	}
	newtag.Value = value

	// parse optional ';' rune
	if r == ';' {
		// consume runes and store in tag list
		p.msg.Tags = append(p.msg.Tags, newtag)
		p.Consume()
		return tagState
	}

	// ending rune for all tags
	if r == ' ' {
		// consume runes and store in tag list
		p.msg.Tags = append(p.msg.Tags, newtag)
		r, ok = p.Next()
		p.Consume()
		if !ok {
			// TODO handle error
			return nil
		}
		// transition to new state
		if r == ':' {
			return prefixState
		}
		return commandState
	}
	return nil
}

// TODO
func prefixState(p *Parser) StateFn {
	return nil
}

// TODO
func commandState(p *Parser) StateFn {
	return nil
}