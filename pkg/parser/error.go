package parser

import (
	"errors"
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

var (
	ParseError = errors.New("Parse error")
)

func (p *Parser) expectPeekError(expect token.TokenType) {
	err := p.expectError(&p.peekToken, expect)
	p.errors = append(p.errors, err)
}
func (p *Parser) expectCurError(expect token.TokenType) {
	err := p.expectError(&p.curToken, expect)
	p.errors = append(p.errors, err)
}

func (p *Parser) expectError(tok *token.Token, expect token.TokenType) error {
	lcStr := lineColString(tok)

	return fmt.Errorf("%s expected `%s`, got=`%s`", lcStr, expect, tok.Type)
}

func (p *Parser) noParsletError(tok *token.Token) {
	lcStr := lineColString(tok)

	err := fmt.Errorf("%s could not parser tok=%s", lcStr, tok.Type)
	p.errors = append(p.errors, err)
}

// will advance until curToken.Type == "SEMICOLON"
func (p *Parser) recover() {
	for !p.atEnd() || p.curToken.Type == token.SEMICOLON {
		p.advance()
	}
}

func lineColString(tok *token.Token) string {
	line, col := tok.Pos.Position()
	return fmt.Sprintf("[%d:%d]", line, col)
}
