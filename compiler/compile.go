package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"

	"github.com/eaglewu/luban/compiler/lexer"
	"github.com/eaglewu/luban/compiler/token"
)

var cpuprofile = flag.String("cpuprof", "", "write cpu profile to file")
var compilefile = flag.String("file", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	buf, err := ioutil.ReadFile(*compilefile)
	if err != nil {
		log.Fatal(err)
	}
	l := lexer.New(string(buf))
	go l.Run()

	//tokens := []token.Token{}
	for tok := l.NextToken(); tok.Type != token.Error && tok.Type != token.End; tok = l.NextToken() {
		//tokens = append(tokens, tok)
		// fmt.Printf("%+v\n", tok)
		// if tok.Line == 102 {
		// 	break
		// }
	}
	fmt.Println("done")
}
