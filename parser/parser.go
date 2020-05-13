package parser

import (
	"strings"
	"unicode"
)

// Message represents a parsed IRC message.
type Message struct {
	Tags    []Tag
	Prefix  Prefix
	Command string
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
// A prefix requires a name, a servername or nickname.
// Optionally, when specifying the nickname, one can specify the host and the
// username.
type Prefix struct {
	Name     string
	Username string
	Host     string
}

// Parser holds state information necessary for parsing IRC messages.
type Parser struct {
	// message currently being built
	msg   Message
	state StateFn
	cin   <-chan []rune
	cout  chan<- Message
	input []rune
	pos   int
}

// NewParser returns a parser with an output buffer of the given size.
//
// The parser is responsible for creating and closing its output channel.
func NewParser(in <-chan []rune, out chan<- Message) *Parser {
	return &Parser{
		state: beginState,
		cin:   in,
		cout:  out,
	}
}

// Run begins the main execution loop for parsing.
//
// The output is passed through the Parser's output channel.
func (p *Parser) Run(cancel <-chan struct{}) {
	for {
		select {
		case rs, ok := <-p.cin:
			if !ok {
				panic("parser: input channel closed prematurely")
			}
			p.input = append(p.input, rs...)
		case <-cancel:
			return
		default:
			if p.state == nil {
				return
			}
			p.state = p.state(p)
		}
	}
}

// Next returns the next rune to be looked at and a boolean for if a rune exists.
func (p *Parser) Next() (rune, bool) {
	if p.pos >= len(p.input) {
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
	p.Consume()
	p.msg = Message{}
	r, ok := p.Next()
	if !ok {
		// TODO handle error
		return nil
	}
	if r == '@' {
		return tagState
	}
	if r == ':' {
		return prefixState
	}
	if (r > 'a' && r < 'z') || (r > 'A' && r < 'Z') || (r > '0' && r < '9') {
		// the rune is part of a command
		p.Rewind()
		return commandState
	}
	return nil
}

func tagState(p *Parser) StateFn {
	// parse tag key
	// - parse vendor hostname optional string
	// - parse any number of alpha, digit, '.', and '-' runes

	var newtag Tag
	var key string
	var vendor string

	p.Consume()
	key = parseUntil(p, isHostnameRune)
	r, ok := p.Next()
	if !ok {
		// TODO handle error
		return nil
	}
	if r == '/' {
		vendor = key
		p.Consume()
		key = parseUntil(p, isHostnameRune)
		if len(key) == 0 {
			// TODO handle error
			return nil
		}
		r, ok = p.Next()
		if !ok {
			// TODO handle error
			return nil
		}
	}
	if strings.ContainsRune(key, '.') {
		// TODO handle error
		return nil
	}
	if vendor != "" {
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
	}
	if key == "" {
		// TODO handle error
		return nil
	}
	newtag.Key = key
	newtag.Vendor = vendor

	// parse optional tag value
	var value string
	if r == '=' {
		p.Consume()
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
		return tagState
	}

	// ending rune for all tags
	if r != ' ' {
		// TODO handle error
		return nil
	}
	if !skipSpaces(p) {
		// TODO handle error
		return nil
	}

	// consume runes and store in tag list
	p.msg.Tags = append(p.msg.Tags, newtag)
	r, ok = p.Next()
	if !ok {
		// TODO handle error
		return nil
	}
	// transition to new state
	if r == ':' {
		return prefixState
	}
	p.Rewind()
	return commandState
}

func prefixState(p *Parser) StateFn {
	p.Consume()
	p.msg.Prefix.Name = parseUntil(p, isPrefixNameRune)

	r, ok := p.Next()
	if !ok {
		// TODO handle error
		return nil
	}

	if r == '!' || r == '@' {
		if r == '!' {
			p.Consume()
			p.msg.Prefix.Username = parseUntil(p, func(r rune) bool {
				return r != '@'
			})
		}

		r, ok = p.Next()
		if !ok || r != '@' {
			// TODO handle error
			return nil
		}
		p.Consume()

		p.msg.Prefix.Host = parseUntil(p, func(r rune) bool {
			return r != ' '
		})
		r, ok = p.Next()
		if !ok {
			// TODO handle error
			return nil
		}
	}

	if r != ' ' {
		// TODO handle error
		return nil
	}
	if !skipSpaces(p) {
		// TODO handle error
		return nil
	}

	return commandState
}

func commandState(p *Parser) StateFn {
	p.Consume()
	cmd := parseUntil(p, isCommandRuneFunc())
	r, ok := p.Next()
	if !ok {
		// TODO handle error
		return nil
	}
	if r == '\n' {
		// remove CR
		cmd = cmd[:len(cmd)-1]
	}
	first := []rune(cmd)[0]
	// verify length
	if unicode.IsNumber(first) {
		if len(cmd) != 3 {
			// TODO handle error
			return nil
		}
		for _, r := range cmd {
			if !unicode.IsNumber(r) {
				// TODO handle error
				return nil
			}
		}
	} else {
		for _, r := range cmd {
			if !unicode.IsLetter(r) {
				// TODO handle error
				return nil
			}
		}
	}
	// verify contents
	p.msg.Command = cmd
	if r == ' ' {
		if !skipSpaces(p) {
			// TODO handle error
			return nil
		}
		return paramState
	} else if r == '\n' {
		return endState
	}
	// TODO handle error
	return nil
}

func paramState(p *Parser) StateFn {
	p.Consume()
	r, ok := p.Next()
	if !ok {
		// TODO handle error
		return nil
	}
	if r == ':' {
		p.Rewind()
		return trailState
	}
	if !isParamRune(r) {
		// TODO handle error
		return nil
	}
	p.msg.Params = append(p.msg.Params, parseUntil(p, isParamMiddleRune))
	r, ok = p.Next()
	if !ok {
		// TODO handle error
		return nil
	}
	if r == ' ' {
		if !skipSpaces(p) {
			// TODO handle error
			return nil
		}
		return paramState
	}
	if r != '\x0D' { // CR
		return trailState
	}
	return endState
}

func trailState(p *Parser) StateFn {
	r, ok := p.Next()
	if !ok {
		// TODO handle error
		return nil
	}
	if r != ':' {
		// TODO handle error
		return nil
	}
	p.Consume()
	p.msg.Params = append(p.msg.Params, parseUntil(p, isTrailingParamRune))
	r, ok = p.Next()
	if !ok {
		// TODO handle error
		return nil
	}
	if r != '\x0D' { // CR
		// TODO handle error
		return nil
	}
	r, ok = p.Next()
	if !ok {
		// TODO handle error
		return nil
	}
	if r != '\x0A' { // LF
		// TODO handle error
		return nil
	}
	return endState
}

func endState(p *Parser) StateFn {
	p.cout <- p.msg
	return beginState
}

// parseUntil parses runes until the given predicate fails for one of the runes.
func parseUntil(p *Parser, pred func(rune) bool) string {
	r, ok := p.Next()
	if !ok {
		return ""
	}
	for pred(r) {
		r, ok = p.Next()
		if !ok {
			return p.Consume()
		}
	}
	p.Rewind()
	return p.Consume()
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

func isPrefixNameRune(r rune) bool {
	switch r {
	case '!':
	case '@':
	case ' ':
	default:
		return true
	}
	return false
}

func isCommandRuneFunc() func(r rune) bool {
	var lastRuneCR bool
	return func(r rune) bool {
		switch {
		case r == ' ':
		case lastRuneCR && r == '\n':
		case r == '\r':
			lastRuneCR = true
			return true
		case unicode.IsLetter(r):
			fallthrough
		case unicode.IsNumber(r):
			lastRuneCR = false
			return true
		}
		return false
	}
}

func isParamRune(r rune) bool {
	switch r {
	case '\x00': // NUL
	case '\x0D': // CR
	case '\x0A': // LF
	case ':':
	case ' ':
	default:
		return true
	}
	return false
}

func isParamMiddleRune(r rune) bool {
	switch r {
	case '\x00': // NUL
	case '\x0D': // CR
	case '\x0A': // LF
	case ' ':
	default:
		return true
	}
	return false
}

func isTrailingParamRune(r rune) bool {
	switch r {
	case '\x00': // NUL
	case '\x0D': // CR
	case '\x0A': // LF
	default:
		return true
	}
	return false
}

func skipSpaces(p *Parser) bool {
	r, ok := p.Next()
	if !ok {
		return false
	}
	for r == ' ' {
		r, ok = p.Next()
		if !ok {
			return false
		}
	}
	p.Rewind()
	return true
}
