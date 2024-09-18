package token_test

import (
	"fmt"
	"testing"

	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

func TestTokenPrint(t *testing.T) {
	tokens := []token.Token{
		token.NewToken(token.LET, "let", nil, 0, 0),
	}

	for _, tok := range tokens {
		fmt.Println(tok)
	}
}