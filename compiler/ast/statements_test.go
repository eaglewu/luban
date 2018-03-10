package ast

import (
	"fmt"
	"testing"

	"github.com/eaglewu/luban/compiler/token"
)

func Test_IfStatement(t *testing.T) {

	exp := &IntegerLiteral{
		Value: 123,
		BaseNode: &BaseNode{
			Token: token.Token{
				Line:    1,
				Type:    token.Lnumber,
				Literal: "123",
			},
		},
	}
	fmt.Println("0", exp)
	fmt.Println("1", &ConditionalExpression{
		Condition: exp,
		Consequence: &BlockStatement{
			Statements: []Statement{},
		},
	})

	stmt := &IfStatement{
		Conditionals: []*ConditionalExpression{
			&ConditionalExpression{
				Condition: exp,
				Consequence: &BlockStatement{
					Statements: []Statement{},
				},
			},
		},
		Alternative: nil,
	}

	fmt.Printf("%s", stmt.String())
}
