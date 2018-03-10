package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/eaglewu/luban/compiler/lexer"
	"github.com/eaglewu/luban/compiler/token"
)

var file = flag.String("file", "", "File of lexical scanning")
var compare = flag.String("compare", "", "compare with json file genrate by native PHP")

const (
	COMPLEX = 1
	SIMPLE  = 2
)

type item struct {
	Type  int    `json:"type"`
	Line  int    `json:"l"`
	Token string `json:"t"`
	Value string `json:"v"`
}

func main() {
	flag.Parse()
	if *file == "" {
		fmt.Fprintf(os.Stderr, "Usage: lexer -file\n")
		os.Exit(-1)
	}
	input := readFile(*file)

	var items []item
	if *compare != "" {
		jsonData := readFile(*compare)
		if err := json.Unmarshal(jsonData, &items); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid JSON: %s\n", err.Error())
			os.Exit(-1)
		}
	}

	lexer := lexer.New(string(input))
	go lexer.Run()
	for n, tok := 0, lexer.NextToken(); tok.Type != token.Error; tok = lexer.NextToken() {
		if tok.Type == token.HaltCompiler || tok.Type == token.End {
			break
		}
		if *compare == "" {
			fmt.Printf(
				"Line: \033[36m%d\033[0m Token: \033[32m %s \033[0m ('%s')\n",
				tok.Line,
				tok.Type,
				tok.Literal,
			)
		} else {
			if n < len(items) {
				if err := compareItem(items[n], tok); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
					os.Exit(-1)
				}
			} else {
				// out of range
				fmt.Fprintf(os.Stderr, "ERROR: [%d] out of items range\n", n)
			}
		}
		n++
	}
}

func readFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(-1)
	}
	return data
}

func compareItem(i item, tok token.Token) error {
	if i.Type == COMPLEX {
		if i.Line != tok.Line {
			return fmt.Errorf("diffent line, Native:%d Go:%d", i.Line, tok.Line)
		}
	}
	if i.Value != tok.Literal {
		return fmt.Errorf(
			"diffent value, Line: %d Native:'\033[32m%s\033[0m' Go:'\033[32m%s\033[0m'",
			tok.Line, i.Value, tok.Literal,
		)
	}
	return nil
}
