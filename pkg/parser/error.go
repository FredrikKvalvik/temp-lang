package parser

import (
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

func (p *Parser) expectCurError(expect token.TokenType) {
	err := p.expectError(p.curToken, expect)
	p.errors = append(p.errors, err)
}

func (p *Parser) expectPeekError(expect token.TokenType) {
	err := p.expectError(p.peekToken, expect)
	p.errors = append(p.errors, err)
}

func (p *Parser) expectError(tok token.Token, expect token.TokenType) error {
	col := p.l.GetTokenColumn(tok)
	line := tok.Line
	got := tok.Type

	return fmt.Errorf("[%d:%d] expected `%s`, got=`%s`\n", line, col, expect, got)
}
