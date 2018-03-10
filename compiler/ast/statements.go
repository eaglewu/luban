package ast

import "bytes"

// ----------------IfStatement----------------

type IfStatement struct {
	*BaseNode
	Conditionals []*ConditionalExpression
	Alternative  *BlockStatement
}

func (st *IfStatement) stmtNode() {}

func (ie *IfStatement) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfStatement) String() string {
	var out bytes.Buffer
	for i, c := range ie.Conditionals {
		if i == 0 {
			out.WriteString("if ( ")
		} else {
			out.WriteString("} elsif ( ")
		}
		out.WriteString(c.String())
		out.WriteString(") {\n")
	}
	if ie.Alternative != nil {
		out.WriteString("\n")
		out.WriteString("} else {\n")
		out.WriteString(ie.Alternative.String())
	}
	out.WriteString("\n}")
	return out.String()
}

// ----------------ClassStatement----------------

type ClassStatement struct {
	*BaseNode
	Name           string
	Body           *BlockStatement
	SuperClass     Expression
	SuperClassName string
}

func (st *ClassStatement) stmtNode() {}

func (cs *ClassStatement) TokenLiteral() string {
	return cs.Token.Literal
}
func (cs *ClassStatement) String() string {
	var out bytes.Buffer

	out.WriteString("class ")
	out.WriteString(cs.Name)
	out.WriteString(" {\n")
	out.WriteString(cs.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type BlockStatement struct {
	*BaseNode
	Statements []Statement
}

func (st *BlockStatement) stmtNode() {}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}
	return out.String()
}

func (bs *BlockStatement) IsEmpty() bool {
	return len(bs.Statements) == 0
}

// ----------------ReturnStatement----------------

type ReturnStatement struct {
	*BaseNode
	ReturnValue Expression
}

func (st *ReturnStatement) stmtNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ----------------ExpressionStatement----------------

type ExpressionStatement struct {
	*BaseNode
	Expression Expression
}

func (st *ExpressionStatement) stmtNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// ----------------WhileStatement----------------

type WhileStatement struct {
	*BaseNode
	Condition Expression
	Body      *BlockStatement
}

func (st *WhileStatement) stmtNode() {}

func (ws *WhileStatement) TokenLiteral() string {
	return ws.Token.Literal
}

func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while( ")
	out.WriteString(ws.Condition.String())
	out.WriteString(") {\n")
	out.WriteString(ws.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// ----------------DoStatement----------------

type DoStatement struct {
	*BaseNode
	Condition Expression
	Body      *BlockStatement
}

func (st *DoStatement) stmtNode() {}

func (ws *DoStatement) TokenLiteral() string {
	return ws.Token.Literal
}

func (ws *DoStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while( ")
	out.WriteString(ws.Condition.String())
	out.WriteString(") {\n")
	out.WriteString(ws.Body.String())
	out.WriteString("\n}")
	return out.String()
}
