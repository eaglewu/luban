package ast

import (
	"bytes"

	"github.com/eaglewu/luban/compiler/token"
)

// const (
// 	SpecialShift     = 6
// 	IsListShift      = 7
// 	NumChildrenShift = 8
// )
// const (
// // special nodes
// Zval Kind = 1 << SpecialShift
// Znode

// // declaration nodes
// FuncDecl
// Closure
// Method
// Class

// // list nodes
// ArgList = 1 << IsListShift
// Array
// EncapsList
// ExprList
// StmtList
// If
// SwitchList
// CatchList
// ParamList
// ClosureUses
// PropDecl
// ConstDecl
// ClassConstDecl
// NameList
// TraitAdaptations
// Use

// // 0 child nodes
// MagicConst = 0 << NumChildrenShift
// Type

// // 1 child node
// Var = 1 << NumChildrenShift
// Const
// Unpack
// UnaryPlus
// UnaryMinus
// C
// Empty
// Isset
// Silence
// ShellExec
// Clone
// Exit
// Print
// IncludeOrEval
// UnaryOp
// PreInc
// PreDec
// PostInc
// PostDec
// YieldFrom

// Global
// Unset
// Return
// Label
// Ref
// HaltCompiler
// Echo
// Throw
// Goto
// Break
// Continue

// // 2 child nodes
// Dim = 2 << NumChildrenShift
// Prop
// StaticProp
// Call
// ClassConst
// Assign
// AssignRef
// AssignOp
// BinaryOp
// Greater
// GreaterEqual
// And
// Or
// ArrayElem
// New
// Instanceof
// Yield
// Coalesce

// Static
// While
// DoWhile
// IfElem
// Switch
// SwitchCase
// Declare
// UseTrait
// TraitPrecedence
// MethodReference
// Namespace
// UseElem
// TraitAlias
// GroupUse

// // 3 child nodes
// MethodCall = 3 << NumChildrenShift
// StaticCall
// Conditional

// Try
// Catch
// Param
// PropElem
// ConstElem

// // 4 child nodes
// For = 4 << NumChildrenShift
// Foreach
// )

type BaseNode struct {
	Token  token.Token
	isStmt bool
}

func (b *BaseNode) Line() int {
	return b.Token.Line
}

func (b *BaseNode) IsExp() bool {
	return !b.isStmt
}

func (b *BaseNode) IsStmt() bool {
	return b.isStmt
}

func (b *BaseNode) MarkAsStmt() {
	b.isStmt = true
}

func (b *BaseNode) MarkAsExp() {
	b.isStmt = false
}

// Node all node types implement the Node interface.
type Node interface {
	TokenLiteral() string
	String() string
	Line() int
	IsExp() bool
	IsStmt() bool

	MarkAsStmt()
	MarkAsExp()
}

// Expr all expression nodes implement the Expr interface.
type Expression interface {
	Node
	exprNode()
}

// Stmt all statement nodes implement the Stmt interface.
type Statement interface {
	Node
	stmtNode()
}

// Program is the root node of entire AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
