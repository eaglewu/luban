package parser

import (
	"fmt"

	"github.com/eaglewu/luban/compiler/ast"
	"github.com/eaglewu/luban/compiler/lexer"
	"github.com/eaglewu/luban/compiler/token"
)

const (
	_ int = iota
	// EndOfFileError represents normal EOF error
	EndOfFileError
	// WrongTokenError means that token is not what we expected
	WrongTokenError
	// UnexpectedTokenError means that token is not expected to appear in current condition
	UnexpectedTokenError
	// UnexpectedEndError means we get unexpected "end" keyword (this is mainly created for REPL)
	UnexpectedEndError
	// MethodDefinitionError means there's an error on method definition's method name
	MethodDefinitionError
	// InvalidAssignmentError means user assigns value to wrong type of expressions
	InvalidAssignmentError
	// SyntaxError means there's a grammatical in the source code
	SyntaxError
	// ArgumentError means there's a method parameter's definition error
	ArgumentError
)

// Error represents parser's parsing error
type Error struct {
	// Message contains the readable message of error
	Message string
	errType int
}

// The Parser structure holds the parser's internal state.
type Parser struct {
	Lexer *lexer.Lexer
	error *Error

	curToken  token.Token
	peekToken token.Token
}

// New parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		Lexer: l,
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
again:
	tok := p.Lexer.NextToken()

	switch tok.Type {
	case token.Comment:
		fallthrough
	case token.DocComment:
		fallthrough
	case token.OpenTag:
		fallthrough
	case token.Whitespace:
		goto again

	case token.CloseTag:
		tok.Type = token.Semicolon
		break
	}

	p.peekToken = tok
}

// ParseProgram parse source input to structure of ast.Program
func (p *Parser) ParseProgram() (*ast.Program, *Error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for typ := p.curToken.Type; typ != token.End; p.nextToken() {
		if typ == token.Error {
			return nil, &Error{Message: p.curToken.Literal}
		}
		stmt := p.parseStatement()
		if p.error != nil {
			return nil, p.error
		}
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}
	return program, nil
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf(
		"unexpected '%s', expecting '%s' in php shell code on line %d",
		p.peekToken.Type, t, p.peekToken.Line,
	)
	p.error = &Error{Message: msg, errType: UnexpectedTokenError}
}
