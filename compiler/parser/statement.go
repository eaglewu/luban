package parser

import (
	"github.com/eaglewu/luban/compiler/ast"
	"github.com/eaglewu/luban/compiler/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {

	case token.LBrace: // '{' inner_statement_list '}'
		return p.parseInnerStatement()

	case token.Include:
		return p.parseIncludeStatement(false)
	case token.IncludeOnce:
		return p.parseIncludeStatement(true)
	case token.Eval:
		return p.parseEvalStatement()
	case token.Require:
		return p.parseIncludeStatement(false)
	case token.RequireOnce:
		return p.parseIncludeStatement(true)
	default:
		return nil
	}
}

func (p *Parser) parseInnerStatement() ast.Statement {
	return nil
}

func (p *Parser) parseIncludeStatement(once bool) ast.Statement {
	return nil
}

func (p *Parser) parseRequireStatement(once bool) ast.Statement {
	return nil
}

func (p *Parser) parseEvalStatement() ast.Statement {
	return nil
}
