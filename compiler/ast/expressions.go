package ast

import "bytes"

// IntegerLiteral contains the node expression and its value
type IntegerLiteral struct {
	*BaseNode
	Value int
}

func (il *IntegerLiteral) exprNode() {}

// TokenLiteral gets the Integer type token
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// String gets the string format of the Integer type token
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// FloatLiteral contains the node expression and its value
type FloatLiteral struct {
	*BaseNode
	Value float64
}

func (e *FloatLiteral) exprNode() {}

// FloatLiteral.TokenLiteral gets the literal of the Float type token
func (il *FloatLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// FloatLiteral.String gets the string format of the Float type token
func (il *FloatLiteral) String() string {
	return il.Token.Literal
}

// StringLiteral contains the node expression and its value
type StringLiteral struct {
	*BaseNode
	Value string
}

func (e *StringLiteral) exprNode() {}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("\"")
	out.WriteString(sl.Token.Literal)
	out.WriteString("\"")
	return out.String()
}

type Identifier struct {
	*BaseNode
	Value string
}

func (e *Identifier) exprNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type BooleanExpression struct {
	*BaseNode
	Value bool
}

func (e *BooleanExpression) exprNode() {}

func (b *BooleanExpression) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BooleanExpression) String() string {
	return b.Token.Literal
}

type NullExpression struct {
	*BaseNode
	Value string
}

func (e *NullExpression) exprNode() {}

func (b *NullExpression) TokenLiteral() string {
	return b.Token.Literal
}

func (b *NullExpression) String() string {
	return b.Token.Literal
}

type ArrayExpression struct {
	*BaseNode
	Elements []Expression
}

func (e *ArrayExpression) exprNode() {}

func (ae *ArrayExpression) TokenLiteral() string {
	return ae.Token.Literal
}

func (ae *ArrayExpression) String() string {
	var out bytes.Buffer

	out.WriteString("[")

	if len(ae.Elements) == 0 {
		out.WriteString("]")
		return out.String()
	}

	out.WriteString(ae.Elements[0].String())

	for _, elem := range ae.Elements[1:] {
		out.WriteString(", ")
		out.WriteString(elem.String())
	}

	out.WriteString("]")
	return out.String()
}

// ConditionalExpression represents if or elsif expression
type ConditionalExpression struct {
	*BaseNode
	Condition   Expression
	Consequence *BlockStatement
}

func (e *ConditionalExpression) exprNode() {}

func (ce *ConditionalExpression) expressionNode() {}

// TokenLiteral returns `if` or `elsif`
func (ce *ConditionalExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *ConditionalExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ce.Condition.String())
	out.WriteString("\n")
	out.WriteString(ce.Consequence.String())

	return out.String()
}

type Constant struct {
	*BaseNode
	Value       string
	IsNamespace bool
}

func (e *Constant) exprNode() {}

func (c *Constant) TokenLiteral() string {
	return c.Token.Literal
}
func (c *Constant) String() string {
	return c.Value
}
