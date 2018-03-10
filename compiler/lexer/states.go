package lexer

import (
	"github.com/eaglewu/luban/compiler/token"
)

const (
	whiteSpace = " \t\r\n"
	digits     = "0123456789"
)

// longestLiteral hard coded of longest literal avoid sort.
// Warning: Must change here while length of literal is
// greater than the "__halt_compiler".
var longestLiteral = len("__halt_compiler") // 15

// simple literal tokens, no mode change.
var literalType = map[string]token.Type{}

func lexInlineHtml(l *Lexer) stateFn {
	for {
		if l.pos > l.start && l.hasPrefix("<?") {
			l.emit(token.InlineHtml)
		}
		if l.hasPrefix("<?php") {
			l.pos += len("<?php")
			if c := l.peek(); isSpace(c) {
				l.pos++
				if c == '\r' && l.peek() == '\n' {
					l.pos++
				}
			}
			l.emit(token.OpenTag)
			l.begin(modeInScript)
			return nil
		}
		if l.hasPrefix("<?=") {
			l.advanceEmit("<?=", token.OpenTagWithEcho)
			l.begin(modeInScript)
			return nil
		}
		if l.hasPrefix("<?") {
			l.advanceEmit("<?", token.OpenTag)
			l.begin(modeInScript)
			return nil
		}
		if l.next() == eof {
			break
		}
	}
	if l.pos > l.start {
		l.emit(token.InlineHtml)
	}
	l.emit(token.End)
	return nil
}

func lexInScript(l *Lexer) stateFn {
	if l.pos >= len(l.input) {
		l.emit(token.End).pop()
		return nil
	}

	// TODO:
	if isSpace(l.peek()) {
		l.acceptRun(whiteSpace)
		//		l.emit(token.Whitespace)
		return lexInScript
	}

	if l.hasPrefix("->") {
		l.pos += len("->")
		l.push(modeLookingForProperty)
		l.emit(token.ObjectOperator)
		return nil
	}

	if l.hasPrefix("#") || l.hasPrefix("//") {
		return lexComment
	}

	if l.hasPrefix("/*") {
		return lexDocComment
	}

	if l.hasPrefix("<<<") { // heredoc
		ori := l.pos
		l.pos += len("<<<")
		l.acceptRun(" \t")
		p1 := l.pos

		c := l.peek()
		if c == '"' || c == '\'' {
			l.pos++
			p1 = l.pos
		}

		if isLabelStart(l.peek()) {
			l.acceptRunLabel()
			if c == '\'' { // nowdoc
				l.acceptRunLabel()
				if l.peek() == '\'' {
					p2 := l.pos
					l.pos++
					if c := l.peek(); isNewline(c) {
						l.pos++
						if c == '\r' && l.peek() == '\n' {
							l.pos++
						}
						l.docLabel = l.input[p1:p2] // save label name
						//	fmt.Printf("label: '%s'\n", l.docLabel)
						l.emit(token.StartHeredoc).begin(modeNowdoc)
						return nil
					}
				}
			} else { // heredoc
				l.acceptRunLabel()
				p2 := l.pos
				if l.peek() == '"' {
					l.pos++
				}
				if c := l.peek(); isNewline(c) {
					l.pos++
					if c == '\r' && l.peek() == '\n' {
						l.pos++
					}
					l.docLabel = l.input[p1:p2] // save label name
					//fmt.Printf("label: '%s'\n", l.docLabel)
					l.emit(token.StartHeredoc).begin(modeHeredoc)
					return nil
				}
			}
		}
		l.pos = ori
	}

	if l.readNumber() {
		if l.accept(".") { // '123.'  -> dnumber
			l.acceptRun(digits)
			if p, ok := isExpNumSufix(l); ok {
				l.pos += p
				l.acceptRun(digits)
			}
			l.emit(token.Dnumber)
		} else {
			if p, ok := isExpNumSufix(l); ok {
				l.pos += p
				l.acceptRun(digits)
				l.emit(token.Dnumber)
			} else {
				l.emit(token.Lnumber)
			}
		}
		return lexInScript
	}

	switch cur := l.peek(); cur {
	case '"':
		l.pos++
		p := l.pos
		for c := l.next(); c != eof; c = l.next() {
			switch c {
			case '"':
				l.emit(token.ConstantEncapsedString)
				return nil
			case '$':
				if r2 := l.peek(); isLabelStart(r2) || r2 == '{' {
					break
				}
				continue
			case '{':
				if l.peek() == '$' {
					break
				}
				continue
			case '\\':
				if l.pos < len(l.input) {
					l.next()
				}
				fallthrough
			default:
				continue
			}
			break
		}
		l.pos = p

		l.emit(token.DoubleQuotes)
		l.begin(modeDoubleQuotes)
		return nil
	case '\'':
		l.pos++
		for c := l.next(); c != eof; c = l.next() {
			switch c {
			case '\'':
				l.emit(token.ConstantEncapsedString)
				return lexInScript
			case '\\':
				l.next()
			}
		}
		l.backup().emit(token.ConstantEncapsedString)
		return lexInScript
	case '`':
		l.pos++
		l.emit(token.Backquote).begin(modeBackquote)
		return nil
	case '{':
		l.pos++
		l.emit(token.LBrace).push(modeInScript)
		return nil
	case '$':
		l.pos++
		if isLabelStart(l.peek()) {
			l.acceptRunLabel()
			l.emit(token.Variable)
		} else {
			l.emit(token.Dollar)
		}
		return lexInScript
	case '}':
		l.pos++
		if len(l.modeStack) > 0 {
			l.pop()
		}
		l.emit(token.RBrace)
		return nil
	case '.':
		l.pos++
		c1, c2 := l.peek(), l.peekN(1)
		if c1 == '.' && c2 == '.' {
			l.pos += 2
			l.emit(token.Ellipsis)
			return lexInScript
		}
		if c1 == '=' {
			l.pos++
			l.emit(token.ConcatEqual)
			return lexInScript
		}
		if isDigit(c1) { // '.3'  ->  number
			l.acceptRun(digits)
			l.emit(token.Dnumber)
			return lexInScript
		}
		l.emit(token.Dot)
		return lexInScript
	case '\\':
		l.pos++
		l.emit(token.NsSeparator)
		return lexInScript
	case ';':
		l.pos++
		l.emit(token.Semicolon)
		return lexInScript
	case ':':
		l.pos++
		if l.peek() == ':' {
			l.pos++
			l.emit(token.PaamayimNekudotayim)
		} else {
			l.emit(token.Colon)
		}
		return lexInScript
	case ',':
		l.pos++
		l.emit(token.Comma)
		return lexInScript
	case '[':
		l.pos++
		l.emit(token.LBracket)
		return lexInScript
	case ']':
		l.pos++
		l.emit(token.RBracket)
		return lexInScript
	case '(':
		l.pos++
		pos := l.pos
		l.acceptRun(" \t")

		if isLetter(l.peek()) {
			var typ token.Type
			if l.hasPrefix("integer") {
				l.pos += len("integer")
				typ = token.IntCast
			} else if l.hasPrefix("int") {
				l.pos += len("int")
				typ = token.IntCast
			} else if l.hasPrefix("real") {
				l.pos += len("real")
				typ = token.DoubleCast
			} else if l.hasPrefix("double") {
				l.pos += len("double")
				typ = token.DoubleCast
			} else if l.hasPrefix("float") {
				l.pos += len("float")
				typ = token.DoubleCast
			} else if l.hasPrefix("string") {
				l.pos += len("string")
				typ = token.StringCast
			} else if l.hasPrefix("binary") {
				l.pos += len("binary")
				typ = token.StringCast
			} else if l.hasPrefix("array") {
				l.pos += len("array")
				typ = token.ArrayCast
			} else if l.hasPrefix("object") {
				l.pos += len("object")
				typ = token.ObjectOperator
			} else if l.hasPrefix("boolean") {
				l.pos += len("boolean")
				typ = token.BoolCast
			} else if l.hasPrefix("bool") {
				l.pos += len("bool")
				typ = token.BoolCast
			} else if l.hasPrefix("unset") {
				l.pos += len("unset")
				typ = token.Unset
			} else {
				goto LParen
			}

			l.acceptRun(" \t")
			if l.peek() == ')' {
				l.pos++
				l.emit(typ)
				return lexInScript
			}
		}
	LParen:
		l.pos = pos
		l.emit(token.LParen)
		return lexInScript
	case ')':
		l.pos++
		l.emit(token.RParen)
		return lexInScript
	case '|':
		l.pos++
		if c := l.peek(); c == '=' {
			l.pos++
			l.emit(token.OrEqual)
		} else if c == '|' {
			l.pos++
			l.emit(token.BooleanOr)
		} else {
			l.emit(token.Bar)
		}
		return lexInScript
	case '^':
		l.pos++
		if l.peek() == '=' {
			l.pos++
			l.emit(token.XorEqual)
		} else {
			l.emit(token.Caret)
		}
		return lexInScript
	case '&':
		l.pos++
		if c := l.peek(); c == '=' {
			l.pos++
			l.emit(token.AndEqual)
		} else if c == '&' {
			l.pos++
			l.emit(token.BooleanAnd)
		} else {
			l.emit(token.Ampersand)
		}
		return lexInScript
	case '+':
		l.pos++
		if c := l.peek(); c == '+' {
			l.pos++
			l.emit(token.Inc)
		} else if c == '=' {
			l.pos++
			l.emit(token.PlusEqual)
		} else {
			l.emit(token.Plus)
		}
		return lexInScript
	case '-':
		l.pos++
		if c := l.peek(); c == '-' {
			l.pos++
			l.emit(token.Dec)
		} else if c == '=' {
			l.pos++
			l.emit(token.MinusEqual)
		} else {
			l.emit(token.Minus)
		}
		return lexInScript
	case '*':
		l.pos++
		if c := l.peek(); c == '=' {
			l.pos++
			l.emit(token.MulEqual)
		} else if c == '*' {
			l.pos++
			if l.peek() == '=' {
				l.pos++
				l.emit(token.MulEqual)
			} else {
				l.emit(token.PowEqual)
			}
		} else {
			l.emit(token.Asterisk)
		}
		return lexInScript
	case '/':
		l.pos++
		if l.peek() == '=' {
			l.pos++
			l.emit(token.DivEqual)
		} else {
			l.emit(token.Slash)
		}
		return lexInScript
	case '=':
		l.pos++
		if c1 := l.peek(); c1 == '=' {
			l.pos++
			if c2 := l.peek(); c2 == '=' {
				l.pos++
				l.emit(token.IsIdentical)
			} else {
				l.emit(token.IsEqual)
			}
		} else if c1 == '>' {
			l.pos++
			l.emit(token.DoubleArrow)
		} else {
			l.emit(token.Assign)
		}
		return lexInScript
	case '%':
		l.pos++
		if l.peek() == '=' {
			l.pos++
			l.emit(token.ModEqual)
		} else {
			l.emit(token.Modulo)
		}
		return lexInScript
	case '!':
		l.pos++
		if c := l.peek(); c == '=' {
			l.pos++
			if l.peek() == '=' {
				l.pos++
				l.emit(token.IsNotIdentical)
			} else {
				l.emit(token.IsNotEqual)
			}
		} else {
			l.emit(token.Bang)
		}
		return lexInScript
	case '~':
		l.pos++
		l.emit(token.Tilde)
		return lexInScript
	case '<':
		l.pos++
		if c := l.peek(); c == '>' {
			l.pos++
			l.emit(token.IsNotEqual)
		} else if c == '=' {
			l.pos++
			if l.peek() == '>' {
				l.pos++
				l.emit(token.Spaceship)
			} else {
				l.emit(token.IsSmallerOrEqual)
			}
		} else if c == '<' {
			l.pos++
			if l.peek() == '=' {
				l.pos++
				l.emit(token.SlEqual)
			} else {
				l.emit(token.Sl)
			}
		} else {
			l.emit(token.Lt)
		}
		return lexInScript
	case '>':
		l.pos++
		if c := l.peek(); c == '=' {
			l.pos++
			l.emit(token.IsGreaterOrEqual)
		} else if c == '>' {
			l.pos++
			if l.peek() == '=' {
				l.pos++
				l.emit(token.SrEqual)
			} else {
				l.emit(token.Sr)
			}
		} else {
			l.emit(token.Gt)
		}
		return lexInScript
	case '?':
		l.pos++
		if c := l.peek(); c == '?' {
			l.pos++
			l.emit(token.Coalesce)
		} else if c == '>' { // ?>
			l.pos++
			if c := l.peek(); isNewline(c) {
				l.pos++
				if c == '\r' && l.peek() == '\n' {
					l.pos++
				}
			}
			l.emit(token.CloseTag)
			l.begin(modeInitial)
			return nil
		} else {
			l.emit(token.QuestionMark)
		}
		return lexInScript
	case '@':
		l.pos++
		l.emit(token.At)
		return lexInScript
		//
	case eof:
		l.pos++
		l.pop()
		return nil
	default:
		if isLabelStart(cur) {
			pos := l.pos
			l.pos++
			l.acceptRunLabel()
			ident := l.input[pos:l.pos]

			if ident == "yield" { // yield from
				pos := l.pos
				l.acceptRunSpace()
				if l.hasPrefix("from") {
					l.pos += len("from")
					if !isLabel(l.peek()) {
						l.emit(token.YieldFrom)
						return lexInScript
					}
				}
				l.pos = pos
				l.emit(token.Yield)
				return lexInScript
			}

			l.emit(token.LookupIdent(ident))
			return lexInScript
		} else if isDigit(cur) {
			l.readNumber()
			if l.accept(".") { // '123.'  -> dnumber
				l.acceptRun(digits)
				if p, ok := isExpNumSufix(l); ok {
					l.pos += p
					l.acceptRun(digits)
				}
				l.emit(token.Dnumber)
			} else {
				if p, ok := isExpNumSufix(l); ok {
					l.pos += p
					l.acceptRun(digits)
					l.emit(token.Dnumber)
				} else {
					l.emit(token.Lnumber)
				}
			}
			return lexInScript
		} else {
			return l.errorf("invalid character `%c` ascii(%d)", cur, cur)
		}
	}
}

func lexDoubleQuotes(l *Lexer) stateFn {

	if l.peek() == '"' {
		l.pos++
		l.emit(token.DoubleQuotes).begin(modeInScript)
		return nil
	}

	if embeddedVariables(l) {
		return nil
	}

	for r := l.next(); r != eof; r = l.next() {
		switch r {
		case '"':
			break
		case '$':
			if r2 := l.peek(); isLabelStart(r2) || r2 == '{' {
				break
			}
			continue
		case '{':
			if l.peek() == '$' {
				break
			}
			continue
		case '\\':
			l.next()
			fallthrough
		default:
			continue
		}

		l.backup()
		break
	}
	l.emit(token.EncapsedAndWhitespace)
	return lexDoubleQuotes
}

func lexVarOffset(l *Lexer) stateFn {
	if l.pos >= len(l.input) {
		l.emit(token.End).pop()
		return nil
	}
	switch cur := l.peek(); cur {
	case '[':
		l.pos++
		l.emit(token.LBracket)
		return lexVarOffset
	case ']':
		l.pos++
		l.emit(token.RBracket).pop()
		return nil
	case '$':
		l.pos++
		if isLabelStart(l.peek()) {
			l.acceptRunLabel()
			l.emit(token.Variable)
			return lexVarOffset
		}
		return l.errorf("Unexpected character in input:  '%c' (ASCII=%d) state=%d", cur, cur, l.mode)
	case ' ', '\n', '\r', '\t', '\\', '\'', '#':
		l.pos++
		l.emit(token.EncapsedAndWhitespace)
		return lexVarOffset
	default:
		if isLabelStart(cur) {
			l.acceptRunLabel()
			l.emit(token.String)
			return lexVarOffset
		} else if isDigit(cur) {
			l.readNumber()
			l.emit(token.NumString)
			return lexVarOffset
		} else {
			return l.errorf("Unexpected character in input:  '%c' (ASCII=%d) state=%d", cur, cur, l.mode)
		}
	}
}

func lexLookingForProperty(l *Lexer) stateFn {
	switch cur := l.peek(); cur {
	case ' ', '\t', '\r', '\n':
		l.acceptRunSpace()
		// l.emit(token.Whitespace)
		return lexLookingForProperty
	case '-':
		if l.peekN(1) == '>' {
			l.pos += len("->")
			l.emit(token.ObjectOperator)
			return lexLookingForProperty
		}
		fallthrough
	default:
		if isLabelStart(cur) {
			l.acceptRunLabel()
			l.emit(token.String).pop()
			return nil
		}
		l.pop()
		return nil
	}
}

func lexComment(l *Lexer) stateFn {
	if l.peek() == '#' {
		l.pos++
	} else {
		l.pos += 2 //
	}
	for r := l.next(); r != eof; r = l.next() {
		switch r {
		case '\r':
			if l.peek() == '\n' {
				l.next()
			}
		case '\n':
			//l.line++
			break
		case '?':
			if l.peek() == '>' {
				l.backup()
				break
			}
			fallthrough
		default:
			continue
		}
		break
	}
	l.emit(token.Comment)
	return lexInScript
}

func lexDocComment(l *Lexer) stateFn {
	var doc bool
	if l.hasPrefix("/**") {
		doc = true
		l.pos += len("/**")
	} else {
		l.pos += len("/*")
	}

	for r := l.next(); r != eof; r = l.next() {
		if r == '*' && l.peek() == '/' {
			// not consume
			break
		}
	}
	if l.more() {
		l.next()
	} else {
		return l.errorf("Unterminated comment starting line %d", l.line)
	}
	if doc {
		l.emit(token.DocComment)
	} else {
		l.emit(token.Comment)
	}
	return lexInScript
}

func lexLookingForVarname(l *Lexer) stateFn {
	if isLabelStart(l.peek()) {
		tmp := l.pos
		l.acceptRunLabel()
		if c := l.peek(); c == '[' || c == '}' {
			l.emit(token.StringVarname)
			goto end
		}
		l.pos = tmp
	}
end:
	l.pop()
	l.push(modeInScript)
	return nil
}

// embeddedVariables collect tokens and change mode
// true  return, means mode has changed
// false return, means mode has no changed
func embeddedVariables(l *Lexer) bool {
	if l.hasPrefix("${") {
		l.pos += len("${")
		l.emit(token.DollarOpenCurlyBraces).push(modeLookingForVarname)
		return true
	}
	if l.hasPrefix("{$") {
		l.pos++ // only consume {
		l.emit(token.CurlyOpen).push(modeInScript)
		return true
	}
	if l.peek() == '$' {
		l.pos++
		if isLabelStart(l.peek()) {
			l.acceptRunLabel()
			l.emit(token.Variable)
			if c := l.peek(); c == '[' {
				l.push(modeVarOffset)
				return true
			}
			if n := len("->"); l.hasPrefix("->") {
				l.pos += n
				if isLabelStart(l.peek()) {
					l.pos -= n
					l.push(modeLookingForProperty)
					return true
				}
				l.pos -= n
			}
			return true
		}
		l.pos--
	}
	return false
}

func lexBackquote(l *Lexer) stateFn {
	if l.peek() == '`' {
		l.pos++
		l.emit(token.Backquote).begin(modeInScript)
		return nil
	}
	if embeddedVariables(l) {
		return nil
	}
	for r := l.next(); r != eof; r = l.next() {
		switch r {
		case '`':
			break
		case '$':
			if r2 := l.peek(); isLabelStart(r2) || r2 == '{' {
				break
			}
			continue
		case '{':
			if l.peek() == '$' {
				break
			}
			continue
		case '\\':
			l.next()
			fallthrough
		default:
			continue
		}
		l.backup()
		break
	}
	l.emit(token.EncapsedAndWhitespace)
	return lexBackquote
}

func lexHeredoc(l *Lexer) stateFn {
	if !l.more() {
		l.begin(modeInScript)
		return nil
	}
	if embeddedVariables(l) {
		return nil
	}
	r := l.next()
	for ; r != eof; r = l.next() {
		switch r {
		case '\r':
			if l.peek() == '\n' {
				l.pos++
			}
			fallthrough
		case '\n':
			if l.hasPrefix(l.docLabel) {
				n, m := len(l.docLabel), 0
				if c := l.peekN(n); c == ';' {
					m++
				}
				if c := l.peekN(n + m); c == '\n' || c == '\r' {
					l.emit(token.EncapsedAndWhitespace)
					l.pos += n
					l.emit(token.EndHeredoc)
					l.begin(modeInScript)
					return nil
				}
			}
			continue
		case '$':
			if r2 := l.peek(); isLabelStart(r2) || r2 == '{' {
				break
			}
			continue
		case '{':
			if l.peek() == '$' {
				break
			}
			continue
		case '\\':
			l.next()
			fallthrough
		default:
			continue
		}
		l.backup()
		break
	}
	l.emit(token.EncapsedAndWhitespace)
	return lexHeredoc
}

func lexNowdoc(l *Lexer) stateFn {
	r := l.next()
	for ; r != eof; r = l.next() {
		switch r {
		case '\r':
			if l.peek() == '\n' {
				l.pos++
			}
			fallthrough
		case '\n':
			if l.hasPrefix(l.docLabel) {
				n, m := len(l.docLabel), 0
				if c := l.peekN(n); c == ';' {
					m++
				}
				if c := l.peekN(n + m); c == '\n' || c == '\r' {
					l.emit(token.EncapsedAndWhitespace)
					l.pos += n
					l.emit(token.EndHeredoc)
					break
				}
			}
			fallthrough
		default:
			continue
		}
		break
	}
	if r == eof {
		l.emit(token.EncapsedAndWhitespace)
	}
	l.begin(modeInScript)
	return nil
}

func isExpNumSufix(l *Lexer) (int, bool) {
	if c1 := l.peek(); c1 == 'e' || c1 == 'E' {
		p := 1
		if c2 := l.peekN(1); c2 == '+' || c2 == '-' {
			p++
		}
		if isDigit(l.peekN(p)) {
			return p, true
		}
	}
	return 0, false
}
