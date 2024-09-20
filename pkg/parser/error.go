package parser

import (
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

func (p *Parser) expectCurError(expect token.TokenType) error {
	return expectError(p.curToken.Line, expect, p.curToken.Type)
}
func (p *Parser) expectPeekError(expect token.TokenType) error {
	return expectError(p.peekToken.Line, expect, p.peekToken.Type)
}
func expectError(line int, expect, got token.TokenType) error {
	return fmt.Errorf("[line: %d] exptected `%s`, got=`%s`\n", line, expect, got)
}
