package token

import "strings"

type Type int

const (
	End Type = iota
	Include
	IncludeOnce
	Eval
	Require
	RequireOnce
	LogicalOr
	LogicalXor
	LogicalAnd
	Print
	Yield
	DoubleArrow
	YieldFrom
	PlusEqual
	MinusEqual
	MulEqual
	DivEqual
	ConcatEqual
	ModEqual
	AndEqual
	OrEqual
	XorEqual
	SlEqual
	SrEqual
	PowEqual
	Coalesce
	BooleanOr
	BooleanAnd
	IsEqual
	IsNotEqual
	IsIdentical
	IsNotIdentical
	Spaceship
	IsSmallerOrEqual
	IsGreaterOrEqual
	Sl
	Sr
	Instanceof
	Inc
	Dec
	IntCast
	DoubleCast
	StringCast
	ArrayCast
	ObjectCast
	BoolCast
	UnsetCast
	Pow
	New
	Clone
	Noelse
	Elseif
	Else
	Endif
	Static
	Abstract
	Final
	Private
	Protected
	Public
	Lnumber
	Dnumber
	String
	Variable
	InlineHtml
	EncapsedAndWhitespace
	ConstantEncapsedString
	StringVarname
	NumString
	Exit
	If
	Echo
	Do
	While
	Endwhile
	For
	Endfor
	Foreach
	Endforeach
	Declare
	Enddeclare
	As
	Switch
	Endswitch
	Case
	Default
	Break
	Continue
	Goto
	Function
	Const
	Return
	Try
	Catch
	Finally
	Throw
	Use
	Insteadof
	Global
	Var
	Unset
	Isset
	Empty
	HaltCompiler
	Class
	Trait
	Interface
	Extends
	Implements
	ObjectOperator
	List
	Array
	Callable
	Line
	File
	Dir
	ClassC
	TraitC
	MethodC
	FuncC
	Comment
	DocComment
	OpenTag
	OpenTagWithEcho
	CloseTag
	Whitespace
	StartHeredoc
	EndHeredoc
	DollarOpenCurlyBraces
	CurlyOpen
	PaamayimNekudotayim
	Namespace
	NsC
	NsSeparator
	Ellipsis
	Error

	// Single character
	Semicolon    // ';'
	Colon        // ':'
	Comma        // ','
	Dot          // '.'
	LBracket     // '['
	RBracket     // ']'
	LParen       // '('
	RParen       // ')'
	Bar          // '|'
	Caret        // '^'
	Ampersand    // '&'
	Plus         // '+'
	Minus        // '-'
	Asterisk     // '*'
	Slash        // '/'
	Assign       // '='
	Modulo       // '%'
	Bang         // '!'
	Tilde        // '~'
	Dollar       // '$'
	Lt           // '<'
	Gt           // '>'
	QuestionMark // '?'
	At           // '@'
	DoubleQuotes // '"'
	LBrace       // '{'
	RBrace       // '}'
	Backquote    // '`'
)

type Token struct {
	Line    int
	Type    Type
	Literal string
}

var tokenName = map[Type]string{
	End:                    "End",
	Include:                "Include",
	IncludeOnce:            "IncludeOnce",
	Eval:                   "Eval",
	Require:                "Require",
	RequireOnce:            "RequireOnce",
	LogicalOr:              "LogicalOr",
	LogicalXor:             "LogicalXor",
	LogicalAnd:             "LogicalAnd",
	Print:                  "Print",
	Yield:                  "Yield",
	DoubleArrow:            "DoubleArrow",
	YieldFrom:              "YieldFrom",
	PlusEqual:              "PlusEqual",
	MinusEqual:             "MinusEqual",
	MulEqual:               "MulEqual",
	DivEqual:               "DivEqual",
	ConcatEqual:            "ConcatEqual",
	ModEqual:               "ModEqual",
	AndEqual:               "AndEqual",
	OrEqual:                "OrEqual",
	XorEqual:               "XorEqual",
	SlEqual:                "SlEqual",
	SrEqual:                "SrEqual",
	PowEqual:               "PowEqual",
	Coalesce:               "Coalesce",
	BooleanOr:              "BooleanOr",
	BooleanAnd:             "BooleanAnd",
	IsEqual:                "IsEqual",
	IsNotEqual:             "IsNotEqual",
	IsIdentical:            "IsIdentical",
	IsNotIdentical:         "IsNotIdentical",
	Spaceship:              "Spaceship",
	IsSmallerOrEqual:       "IsSmallerOrEqual",
	IsGreaterOrEqual:       "IsGreaterOrEqual",
	Sl:                     "Sl",
	Sr:                     "Sr",
	Instanceof:             "Instanceof",
	Inc:                    "Inc",
	Dec:                    "Dec",
	IntCast:                "IntCast",
	DoubleCast:             "DoubleCast",
	StringCast:             "StringCast",
	ArrayCast:              "ArrayCast",
	ObjectCast:             "ObjectCast",
	BoolCast:               "BoolCast",
	UnsetCast:              "UnsetCast",
	Pow:                    "Pow",
	New:                    "New",
	Clone:                  "Clone",
	Noelse:                 "Noelse",
	Elseif:                 "Elseif",
	Else:                   "Else",
	Endif:                  "Endif",
	Static:                 "Static",
	Abstract:               "Abstract",
	Final:                  "Final",
	Private:                "Private",
	Protected:              "Protected",
	Public:                 "Public",
	Lnumber:                "Lnumber",
	Dnumber:                "Dnumber",
	String:                 "String",
	Variable:               "Variable",
	InlineHtml:             "InlineHtml",
	EncapsedAndWhitespace:  "EncapsedAndWhitespace",
	ConstantEncapsedString: "ConstantEncapsedString",
	StringVarname:          "StringVarname",
	NumString:              "NumString",
	Exit:                   "Exit",
	If:                     "If",
	Echo:                   "Echo",
	Do:                     "Do",
	While:                  "While",
	Endwhile:               "Endwhile",
	For:                    "For",
	Endfor:                 "Endfor",
	Foreach:                "Foreach",
	Endforeach:             "Endforeach",
	Declare:                "Declare",
	Enddeclare:             "Enddeclare",
	As:                     "As",
	Switch:                 "Switch",
	Endswitch:              "Endswitch",
	Case:                   "Case",
	Default:                "Default",
	Break:                  "Break",
	Continue:               "Continue",
	Goto:                   "Goto",
	Function:               "Function",
	Const:                  "Const",
	Return:                 "Return",
	Try:                    "Try",
	Catch:                  "Catch",
	Finally:                "Finally",
	Throw:                  "Throw",
	Use:                    "Use",
	Insteadof:              "Insteadof",
	Global:                 "Global",
	Var:                    "Var",
	Unset:                  "Unset",
	Isset:                  "Isset",
	Empty:                  "Empty",
	HaltCompiler:           "HaltCompiler",
	Class:                  "Class",
	Trait:                  "Trait",
	Interface:              "Interface",
	Extends:                "Extends",
	Implements:             "Implements",
	ObjectOperator:         "ObjectOperator",
	List:                   "List",
	Array:                  "Array",
	Callable:               "Callable",
	Line:                   "__LINE__",
	File:                   "__FILE__",
	Dir:                    "__DIR__",
	ClassC:                 "__CLASS__",
	TraitC:                 "__TRAIT__",
	MethodC:                "__METHOD__",
	FuncC:                  "__FUNCTION__",
	Comment:                "Comment",
	DocComment:             "DocComment",
	OpenTag:                "OpenTag",
	OpenTagWithEcho:        "OpenTagWithEcho",
	CloseTag:               "CloseTag",
	Whitespace:             "Whitespace",
	StartHeredoc:           "StartHeredoc",
	EndHeredoc:             "EndHeredoc",
	DollarOpenCurlyBraces:  "DollarOpenCurlyBraces",
	CurlyOpen:              "CurlyOpen",
	PaamayimNekudotayim:    "PaamayimNekudotayim",
	Namespace:              "Namespace",
	NsC:                    "__NAMESPACE__",
	NsSeparator:            "NsSeparator",
	Ellipsis:               "Ellipsis",
	Error:                  "Error",

	// Single character
	Semicolon:    "Semicolon",
	Colon:        "Colon",
	Comma:        "Comma",
	Dot:          "Dot",
	LBracket:     "LBracket",
	RBracket:     "RBracket",
	LParen:       "LParen",
	RParen:       "RParen",
	Bar:          "Bar",
	Caret:        "Caret",
	Ampersand:    "Ampersand",
	Plus:         "Plus",
	Minus:        "Minus",
	Asterisk:     "Asterisk",
	Slash:        "Slash",
	Assign:       "Assign",
	Modulo:       "Modulo",
	Bang:         "Bang",
	Tilde:        "Tilde",
	Dollar:       "Dollar",
	Lt:           "Lt",
	Gt:           "Gt",
	QuestionMark: "QuestionMark",
	At:           "At",
	DoubleQuotes: "DoubleQuotes",
	LBrace:       "LBrace",
	RBrace:       "RBrace",
	Backquote:    "Backquote",
}

func (t Type) String() string {
	if n, ok := tokenName[t]; ok {
		return n
	}
	return "Unknown"
}

var keywords = map[string]Type{
	"abstract":     Abstract,
	"and":          BooleanAnd,
	"array":        Array,
	"as":           As,
	"break":        Break,
	"callable":     Callable,
	"case":         Case,
	"catch":        Catch,
	"class":        Class,
	"clone":        Clone,
	"const":        Const,
	"continue":     Continue,
	"declare":      Declare,
	"default":      Default,
	"die":          Exit,
	"do":           Do,
	"echo":         Echo,
	"else":         Else,
	"elseif":       Elseif,
	"empty":        Empty,
	"enddeclare":   Enddeclare,
	"endfor":       Endfor,
	"endforeach":   Endforeach,
	"endif":        Endif,
	"endswitch":    Endswitch,
	"endwhile":     Endwhile,
	"eval":         Eval,
	"exit":         Exit,
	"extends":      Extends,
	"final":        Final,
	"finally":      Finally,
	"for":          For,
	"foreach":      Foreach,
	"function":     Function,
	"global":       Global,
	"goto":         Goto,
	"if":           If,
	"implements":   Implements,
	"include":      Include,
	"include_once": IncludeOnce,
	"instanceof":   Instanceof,
	"insteadof":    Insteadof,
	"interface":    Interface,
	"isset":        Isset,
	"list":         List,
	"namespace":    Namespace,
	"new":          New,
	"or":           BooleanOr,
	"print":        Print,
	"private":      Private,
	"protected":    Protected,
	"public":       Public,
	"require":      Require,
	"require_once": RequireOnce,
	"return":       Return,
	"static":       Static,
	"switch":       Switch,
	"throw":        Throw,
	"trait":        Trait,
	"try":          Try,
	"unset":        Unset,
	"use":          Use,
	"var":          Var,
	"while":        While,
}

var identifiers = map[string]Type{
	"exit":       Exit,
	"die":        Exit,
	"function":   Function,
	"const":      Const,
	"return":     Return,
	"yield":      Yield,
	"try":        Try,
	"catch":      Catch,
	"finally":    Finally,
	"throw":      Throw,
	"if":         If,
	"elseif":     Elseif,
	"endif":      Endif,
	"else":       Else,
	"while":      While,
	"endwhile":   Endwhile,
	"do":         Do,
	"for":        For,
	"endfor":     Endfor,
	"foreach":    Foreach,
	"endforeach": Endforeach,
	"declare":    Enddeclare,
	"instanceof": Instanceof,
	"as":         As,
	"switch":     Switch,
	"endswitch":  Endswitch,
	"case":       Case,
	"default":    Default,
	"break":      Break,
	"continue":   Continue,
	"goto":       Goto,
	"echo":       Echo,
	"print":      Print,
	"class":      Class,
	"interface":  Interface,
	"trait":      Trait,
	"extends":    Extends,
	"implements": Implements,

	"new":             New,
	"clone":           Clone,
	"var":             Var,
	"eval":            Eval,
	"include":         Include,
	"include_once":    IncludeOnce,
	"require":         Require,
	"require_once":    RequireOnce,
	"namespace":       Namespace,
	"use":             Use,
	"insteadof":       Insteadof,
	"global":          Global,
	"isset":           Isset,
	"empty":           Empty,
	"__halt_compiler": HaltCompiler,
	"static":          Static,
	"abstract":        Abstract,
	"final":           Final,
	"private":         Private,
	"protected":       Protected,
	"public":          Public,
	"unset":           Unset,
	"list":            List,
	"array":           Array,
	"callable":        Callable,

	"__class__":     ClassC,
	"__trait__":     TraitC,
	"__function__":  FuncC,
	"__method__":    MethodC,
	"__line__":      Line,
	"__file__":      File,
	"__dir__":       Dir,
	"__namespace__": NsC,

	"or":  LogicalOr,
	"and": LogicalAnd,
	"xor": LogicalXor,
}

func LookupIdent(ident string) Type {
	if t, ok := identifiers[strings.ToLower(ident)]; ok {
		return t
	}
	return String
}

func NewToken(t Type, literal string, line int) Token {
	return Token{Type: t, Literal: literal, Line: line}
}
