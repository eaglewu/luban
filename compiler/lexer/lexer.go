package lexer

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/eaglewu/luban/compiler/token"
)

type mode int

const (
	modeInitial            mode = iota // 0
	modeInScript                       // 1
	modeLookingForProperty             // 2
	modeBackquote                      // 3
	modeDoubleQuotes                   // 4
	modeHeredoc                        // 5
	modeLookingForVarname              // 6
	modeVarOffset                      // 7
	modeNowdoc                         // 8
)

var modeEntries = map[mode]stateFn{
	modeInitial:            lexInlineHtml,
	modeInScript:           lexInScript,
	modeDoubleQuotes:       lexDoubleQuotes,
	modeVarOffset:          lexVarOffset,
	modeLookingForProperty: lexLookingForProperty,
	modeLookingForVarname:  lexLookingForVarname,
	modeBackquote:          lexBackquote,
	modeHeredoc:            lexHeredoc,
	modeNowdoc:             lexNowdoc,
}

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// ErrEmptyStack error is for pop an empty stack
var ErrEmptyStack = errors.New("Stack is empty")

// Lexer is used for tokenizing programs
type Lexer struct {
	name      string           // the name of the input; used only for error reports
	input     string           // the string being scanned
	pos       int              // current position in the input
	start     int              // start position of this item
	width     int              // width of last rune read from input
	line      int              // 1+number of newlines seen
	tokens    chan token.Token // channel of scanned tokens
	mode      mode
	modeStack []mode
	abort     bool
	docLabel  string
}

// New initializes a new lexer with input string
func New(input string) *Lexer {
	l := &Lexer{
		input:     input,
		line:      1,
		tokens:    make(chan token.Token),
		mode:      modeInitial,
		modeStack: make([]mode, 0),
	}
	return l
}

// Run runs the state machine for the lexer.
func (l *Lexer) Run() {
	for {
		for state := modeEntries[l.mode]; state != nil; {
			state = state(l)
		}
		if l.abort {
			break
		}
	}
	close(l.tokens)
}

// NextToken makes lexer tokenize next character(s)
func (l *Lexer) NextToken() token.Token {
	tok := <-l.tokens
	if tok.Type == token.End || tok.Type == token.Error {
		l.abort = true
	}
	return tok
}

func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *Lexer) backup() *Lexer {
	l.pos -= l.width
	return l
}

func (l *Lexer) advanceEmit(s string, t token.Type, accpet ...string) {
	l.pos += len(s)
	if len(accpet) > 0 {
		l.acceptRun(accpet[0])
	}
	l.emit(t)
}

func (l *Lexer) peek() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

func (l *Lexer) peekN(n int) rune {
	if int(l.pos+n) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos+n:])
	return r
}

func (l *Lexer) emit(t token.Type) *Lexer {
	l.tokens <- token.Token{Line: l.line, Type: t, Literal: l.input[l.start:l.pos]}
	l.line += strings.Count(l.input[l.start:l.pos], "\n")
	l.start = l.pos
	return l
}

// accept consumes the next rune if it's from the valid set.
func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// acceptRunFn consumes a run of runes from the valid set function.
func (l *Lexer) acceptRunFn(validFn func(r rune) bool) {
	for validFn(l.next()) {
	}
	l.backup()
}

var labelAccpetFn = func(r rune) bool {
	if isLabel(r) {
		return true
	}
	return false
}

var spaceAccpetFn = func(r rune) bool {
	if isSpace(r) {
		return true
	}
	return false
}

// acceptRunLabel consumes a run of runes from the valid label.
func (l *Lexer) acceptRunLabel() {
	l.acceptRunFn(labelAccpetFn)
}

// acceptRunLabel consumes a run of runes from the valid label.
func (l *Lexer) acceptRunSpace() {
	l.acceptRunFn(spaceAccpetFn)
}

func (l *Lexer) hasPrefix(prefix string) bool {
	return strings.HasPrefix(l.input[l.pos:], prefix)
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token.Token{Line: l.line, Type: token.Error, Literal: fmt.Sprintf(format, args...)}
	return nil
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *Lexer) drain() {
	for range l.tokens {
	}
}

func (l *Lexer) pop() (mode, error) {
	n := len(l.modeStack) - 1
	if n < 0 {
		return 0, ErrEmptyStack
	}

	v := l.modeStack[n]
	l.modeStack = l.modeStack[0:n]
	//fmt.Printf("pop mode changed from %v to %#v Stack: %+v\n", l.mode, v, l.modeStack)
	l.mode = v
	return v, nil
}

func (l *Lexer) begin(m mode) {
	//fmt.Printf("begin mode changed from %v to %v Stack: %+v\n", l.mode, m, l.modeStack)
	l.mode = m
}

func (l *Lexer) push(m mode) {
	l.modeStack = append(l.modeStack, l.mode)
	//fmt.Printf("push mode changed from %v to %v Stack: %+v\n", l.mode, m, l.modeStack)
	l.mode = m
}

func (l *Lexer) more() bool {
	return l.pos < len(l.input)
}

func (l *Lexer) readNumber() bool {
	if !isDigit(l.peek()) {
		return false
	}
	if l.next() == '0' { // not decimal
		c1, c2 := l.peek(), l.peekN(1)
		if c1 == 'b' && (c2 == '0' || c2 == '1') { // binary
			l.pos++
			l.acceptRun("01")
			return true
		}
		if (c1 == 'x' || c1 == 'X') && isHex(c2) { // 0x Or 0X  hex
			l.pos++
			l.acceptRun("0123456789abcdefABCDEF")
			return true
		}
		l.acceptRun("01234567") // octal
		return true
	}
	l.acceptRun(digits) // decimal
	return true
}

func (l *Lexer) String() string {
	return fmt.Sprintf("=====Lexer======\nname: %s\ninputLen: %d\nstart: %d\npos: %d\nline: %d\nmode: %d\nmodeStack: %+v\nabort: %v\nlastWidth: %d\n==========",
		l.name, len(l.input), l.start, l.pos, l.line, l.mode, l.modeStack, l.abort, l.width)
}

func lex(input string) *Lexer {
	l := New(input)
	go l.Run()
	return l
}

func isDigit(ch rune) bool {
	return ('0' <= ch) && (ch <= '9')
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isLabelStart(ch rune) bool {
	return isLetter(ch) || ch == '_' || 0x80 <= ch && ch <= utf8.MaxRune
}

func isLabel(ch rune) bool {
	return isLabelStart(ch) || '0' <= ch && ch <= '9'
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

func isNewline(r rune) bool {
	return r == '\r' || r == '\n'
}

func isHex(r rune) bool {
	return '0' <= r && r <= '9' || 'a' <= r && r <= 'f' || 'A' <= r && r <= 'F'
}
