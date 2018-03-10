package token

import (
	"fmt"
	"log"
	"testing"
	"unsafe"
)

func Test_MemoryUsageOfHugeTokens(t *testing.T) {
	tok := NewToken(OpenTag, "<?php", 1)
	log.Printf("Sizeof Token: %d\n", unsafe.Sizeof(tok))
}

func Test_MemoryAligment(t *testing.T) {
	type A struct {
		a bool
		c string
		b bool
	}
	type B struct {
		a bool
		b bool
		c string
	}
	fmt.Printf("Sizeof A: %d\n", unsafe.Sizeof(A{}))
	fmt.Printf("Sizeof B: %d\n", unsafe.Sizeof(B{}))
}
