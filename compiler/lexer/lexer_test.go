package lexer

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/eaglewu/luban/compiler/token"
)

var input = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>Document</title>
</head>
<body><?php 
	if ($abc) {
		echo 'YES';
	}
	123;0b01001;0x80;017;
	$obj = new stdClass();
	$obj->name = 'Hello';
	// comment
	# comment
	/**
	docComment
	*/
	/* comment */
	"hello";
	"hello$a $a[123] {$a[name]} {$a['name']} {$a[$a[$b]]} ${a}";
	$a=<<<'DOT'
hello
DOT;
	$a=<<<"TEXT"
hello $a ${a}
TEXT;
`

// ?></body>
// </html>`

type testToken struct {
	expectedType    token.Type
	expectedLiteral string
	expectedLine    int
}

var tests = []testToken{
	{token.InlineHtml, `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>Document</title>
</head>
<body>`, 1},
	{token.OpenTag, "<?php ", 9},
	{token.Whitespace, "\n\t", 9},
	{token.If, "if", 10},
	{token.Whitespace, " ", 10},
	{token.LParen, "(", 10},
	{token.Variable, "$abc", 10},
	{token.RParen, ")", 10},
	{token.Whitespace, " ", 10},
	{token.LBrace, "{", 10},
	{token.Whitespace, "\n\t\t", 10},
	{token.Echo, "echo", 11},
	{token.Whitespace, " ", 11},
	{token.ConstantEncapsedString, "'YES'", 11},
	{token.Semicolon, ";", 11},
	{token.Whitespace, "\n\t", 11},
	{token.RBrace, "}", 12},
	{token.Whitespace, "\n\t", 12},
	{token.Lnumber, "123", 13},
	{token.Semicolon, ";", 13},
	{token.Lnumber, "0b01001", 13},
	{token.Semicolon, ";", 13},
	{token.Lnumber, "0x80", 13},
	{token.Semicolon, ";", 13},
	{token.Lnumber, "017", 13},
	{token.Semicolon, ";", 13},
	{token.Whitespace, "\n\t", 13},
	{token.Variable, "$obj", 14},
	{token.Whitespace, " ", 14},
	{token.Assign, "=", 14},
	{token.Whitespace, " ", 14},
	{token.New, "new", 14},
	{token.Whitespace, " ", 14},
	{token.String, "stdClass", 14},
	{token.LParen, "(", 14},
	{token.RParen, ")", 14},
	{token.Semicolon, ";", 14},
	{token.Whitespace, "\n\t", 14},
	{token.Variable, "$obj", 15},
	{token.ObjectOperator, "->", 15},
	{token.String, "name", 15},
	{token.Whitespace, " ", 15},
	{token.Assign, "=", 15},
	{token.Whitespace, " ", 15},
	{token.ConstantEncapsedString, "'Hello'", 15},
	{token.Semicolon, ";", 15},
	{token.Whitespace, "\n\t", 15},
	{token.Comment, "// comment\n", 16},
	{token.Whitespace, "\t", 17},
	{token.Comment, "# comment\n", 17},
	{token.Whitespace, "\t", 18},
	{token.DocComment, "/**\n\tdocComment\n\t*/", 18},
	{token.Whitespace, "\n\t", 20},
	{token.Comment, "/* comment */", 21},
	{token.Whitespace, "\n\t", 21},
	{token.ConstantEncapsedString, "\"hello\"", 22},
	{token.Semicolon, ";", 22},
	{token.Whitespace, "\n\t", 22},
	{token.DoubleQuotes, "\"", 23},
	{token.EncapsedAndWhitespace, "hello", 23},
	{token.Variable, "$a", 23},
	{token.EncapsedAndWhitespace, " ", 23},
	{token.Variable, "$a", 23},
	{token.LBracket, "[", 23},
	{token.NumString, "123", 23},
	{token.RBracket, "]", 23},
	{token.EncapsedAndWhitespace, " ", 23},
	{token.CurlyOpen, "{", 23},
	{token.Variable, "$a", 23},
	{token.LBracket, "[", 23},
	{token.String, "name", 23},
	{token.RBracket, "]", 23},
	{token.RBrace, "}", 23},
	{token.EncapsedAndWhitespace, " ", 23},

	{token.CurlyOpen, "{", 23},
	{token.Variable, "$a", 23},
	{token.LBracket, "[", 23},
	{token.ConstantEncapsedString, "'name'", 23},
	{token.RBracket, "]", 23},
	{token.RBrace, "}", 23},

	{token.EncapsedAndWhitespace, " ", 23},
	{token.CurlyOpen, "{", 23},
	{token.Variable, "$a", 23},
	{token.LBracket, "[", 23},
	{token.Variable, "$a", 23},
	{token.LBracket, "[", 23},
	{token.Variable, "$b", 23},
	{token.RBracket, "]", 23},
	{token.RBracket, "]", 23},
	{token.RBrace, "}", 23},
	{token.EncapsedAndWhitespace, " ", 23},
	{token.DollarOpenCurlyBraces, "${", 23},
	{token.StringVarname, "a", 23},
	{token.RBrace, "}", 23},

	{token.DoubleQuotes, "\"", 23},
	{token.Semicolon, ";", 23},
	{token.Whitespace, "\n\t", 23},
	{token.Variable, "$a", 24},
	{token.Assign, "=", 24},
	{token.StartHeredoc, "<<<'DOT'\n", 24},
	{token.EncapsedAndWhitespace, "hello\n", 25},
	{token.EndHeredoc, "DOT", 26},
	{token.Semicolon, ";", 26},
	{token.Whitespace, "\n\t", 26},
	{token.Variable, "$a", 27},
	{token.Assign, "=", 27},
	{token.StartHeredoc, "<<<\"TEXT\"\n", 27},
	{token.EncapsedAndWhitespace, "hello ", 28},
	{token.Variable, "$a", 28},
	{token.EncapsedAndWhitespace, " ", 28},
	{token.DollarOpenCurlyBraces, "${", 28},
	{token.StringVarname, "a", 28},
	{token.RBrace, "}", 28},
	{token.EncapsedAndWhitespace, "\n", 28},
	{token.EndHeredoc, "TEXT", 29},

	// {token.Whitespace, "\n", 13},
	// {token.CloseTag, "?>", 14},
	// {token.InlineHtml, "</body>\n</html>", 14},
}

func Test_NextToken(t *testing.T) {
	l := lex(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if err := compareToken(tt, tok); err != nil {
			fmt.Printf("%+v\n", tok)
			fmt.Printf("%s\n", l)
			t.Fatalf("tests[%d] - %s", i, err.Error())
		}
	}
}

func Test_Backquote(t *testing.T) {
	script := "<?php $a=`hello $b ${b} $c[123]`;"
	toks := []testToken{
		{token.OpenTag, "<?php ", 1},
		{token.Variable, "$a", 1},
		{token.Assign, "=", 1},
		{token.Backquote, "`", 1},
		{token.EncapsedAndWhitespace, "hello ", 1},
		{token.Variable, "$b", 1},
		{token.EncapsedAndWhitespace, " ", 1},
		{token.DollarOpenCurlyBraces, "${", 1},
		{token.StringVarname, "b", 1},
		{token.RBrace, "}", 1},
		{token.EncapsedAndWhitespace, " ", 1},
		{token.Variable, "$c", 1},
		{token.LBracket, "[", 1},
		{token.NumString, "123", 1},
		{token.RBracket, "]", 1},
		{token.Backquote, "`", 1},
		{token.Semicolon, ";", 1},
	}
	l := lex(script)
	for i, tt := range toks {
		tok := l.NextToken()
		if err := compareToken(tt, tok); err != nil {
			fmt.Printf("%s\n", l)
			t.Fatalf("tests[%d] - %s", i, err.Error())
		}
	}
}

func Test_Unicode(t *testing.T) {
	script := "<?php $🙂=`🚗🚴🚣🌺 $🎨$中国`; class 中国{ public $flag='🇨🇳' }"
	toks := []testToken{
		{token.OpenTag, "<?php ", 1},
		{token.Variable, "$🙂", 1},
		{token.Assign, "=", 1},
		{token.Backquote, "`", 1},
		{token.EncapsedAndWhitespace, "🚗🚴🚣🌺 ", 1},
		{token.Variable, "$🎨", 1},
		{token.Variable, "$中国", 1},
		{token.Backquote, "`", 1},
		{token.Semicolon, ";", 1},
		{token.Whitespace, " ", 1},
		{token.Class, "class", 1},
		{token.Whitespace, " ", 1},
		{token.String, "中国", 1},
		{token.LBrace, "{", 1},
		{token.Whitespace, " ", 1},
		{token.Public, "public", 1},
		{token.Whitespace, " ", 1},
		{token.Variable, "$flag", 1},
		{token.Assign, "=", 1},
		{token.ConstantEncapsedString, "'🇨🇳'", 1},
		{token.Whitespace, " ", 1},
		{token.RBrace, "}", 1},
	}
	l := lex(script)
	for i, tt := range toks {
		tok := l.NextToken()
		if err := compareToken(tt, tok); err != nil {
			fmt.Printf("%+v\n", tok)
			fmt.Printf("%s\n", l)
			t.Fatalf("tests[%d] - %s", i, err.Error())
		}
	}
}

func Test_StringWithKeywordPrefix(t *testing.T) {
	script := "<?php function_exists(); !=="
	toks := []testToken{
		{token.OpenTag, "<?php ", 1},
		{token.String, "function_exists", 1},
		{token.LParen, "(", 1},
		{token.RParen, ")", 1},
		{token.Semicolon, ";", 1},
		{token.Whitespace, " ", 1},
		{token.IsNotIdentical, "!==", 1},
	}
	l := lex(script)
	for i, tt := range toks {
		tok := l.NextToken()
		if err := compareToken(tt, tok); err != nil {
			fmt.Printf("%+v\n", tok)
			fmt.Printf("%s\n", l)
			t.Fatalf("tests[%d] - %s", i, err.Error())
		}
	}
}

func Test_Number(t *testing.T) {
	script := "<?php 123;.3;2.33;012.33;32.;036;038;0x00Ae;0xFF;0x0000Ab1cF;4.1E+6;4.1E-6;4.1E6;4.1Ex;"
	toks := []testToken{
		{token.OpenTag, "<?php ", 1},
		{token.Lnumber, "123", 1},
		{token.Semicolon, ";", 1},
		{token.Dnumber, ".3", 1},
		{token.Semicolon, ";", 1},
		{token.Dnumber, "2.33", 1},
		{token.Semicolon, ";", 1},
		{token.Dnumber, "012.33", 1},
		{token.Semicolon, ";", 1},
		{token.Dnumber, "32.", 1},
		{token.Semicolon, ";", 1},
		{token.Lnumber, "036", 1},
		{token.Semicolon, ";", 1},
		{token.Lnumber, "03", 1},
		{token.Lnumber, "8", 1}, // only test octal, sytax error
		{token.Semicolon, ";", 1},
		{token.Lnumber, "0x00Ae", 1},
		{token.Semicolon, ";", 1},
		{token.Lnumber, "0xFF", 1},
		{token.Semicolon, ";", 1},
		{token.Lnumber, "0x0000Ab1cF", 1},
		{token.Semicolon, ";", 1},
		{token.Dnumber, "4.1E+6", 1},
		{token.Semicolon, ";", 1},
		{token.Dnumber, "4.1E-6", 1},
		{token.Semicolon, ";", 1},
		{token.Dnumber, "4.1E6", 1},
		{token.Semicolon, ";", 1},
		{token.Dnumber, "4.1", 1},
		{token.String, "Ex", 1},
		{token.Semicolon, ";", 1},
	}
	l := lex(script)
	for i, tt := range toks {
		tok := l.NextToken()
		if err := compareToken(tt, tok); err != nil {
			fmt.Printf("%+v\n", tok)
			fmt.Printf("%s\n", l)
			t.Fatalf("tests[%d] - %s", i, err.Error())
		}
	}
}

func Benchmark_Test_Scripts(b *testing.B) {
	buf, err := ioutil.ReadFile("../test-scripts/run-tests.php")
	if err != nil {
		log.Fatal(err)
	}
	str := string(buf)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l := lex(str)
			for tok := l.NextToken(); tok.Type != token.Error && tok.Type != token.End; tok = l.NextToken() {
			}
		}()
	}
	wg.Wait()
}

func Test_Scripts(t *testing.T) {
	buf, err := ioutil.ReadFile("../test-scripts/run-tests.php")
	if err != nil {
		log.Fatal(err)
	}
	str := string(buf)

	start := time.Now()

	l := lex(str)
	for tok := l.NextToken(); tok.Type != token.Error && tok.Type != token.End; tok = l.NextToken() {
	}
	fmt.Println("Latency: ", time.Now().Sub(start))
}
func Benchmark_NextToken(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		l := lex(input)
		for i, tt := range tests {
			tok := l.NextToken()
			if err := compareToken(tt, tok); err != nil {
				fmt.Printf("%s\n", l)
				b.Fatalf("tests[%d] - %s", i, err.Error())
			}
		}
	}
}

func compareToken(t testToken, got token.Token) error {
	if got.Type != t.expectedType {
		return fmt.Errorf("tokentype wrong. expected=%q, got=%q", t.expectedType, got.Type)
	}
	if got.Literal != t.expectedLiteral {
		return fmt.Errorf("literal wrong. expected=%q, got=%q", t.expectedLiteral, got.Literal)
	}
	if got.Line != t.expectedLine {
		return fmt.Errorf("line number wrong. expected=%d, got=%d", t.expectedLine, got.Line)
	}
	return nil
}
