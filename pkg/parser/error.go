package parser

import (
	"fmt"

	"github.com/fredrikkvalvik/temp-lang/pkg/lexer"
	"github.com/fredrikkvalvik/temp-lang/pkg/token"
)

func (p *Parser) expectPeekError(expect token.TokenType) {
	err := p.expectError(&p.peekToken, expect)
	p.errors = append(p.errors, err)
}

func (p *Parser) expectError(tok *token.Token, expect token.TokenType) error {
	lcStr := lineColString(p.l, tok)

	return fmt.Errorf("%s expected `%s`, got=`%s`\n", lcStr, expect, tok.Type)
}

func (p *Parser) noParsletError(tok *token.Token) {
	lcStr := lineColString(p.l, tok)

	err := fmt.Errorf("%s could not parser tok=%s\n", lcStr, tok.Type)
	p.errors = append(p.errors, err)
}

// will advance until curToken.Type == "SEMICOLON"
func (p *Parser) recover() {
	for !p.atEnd() || p.curToken.Type == token.SEMICOLON {
		p.advance()
	}
}

func lineColString(l *lexer.Lexer, tok *token.Token) string {
	line, col := getLineCol(l, tok)
	return fmt.Sprintf("[%d:%d]", line, col)
}

func getLineCol(l *lexer.Lexer, tok *token.Token) (int, int) {
	line, col := l.GetTokenPosition(tok)

	return line, col
}
