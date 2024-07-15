package parser

import (
	"fmt"

	"github.com/arsenzairov/cs153/scanner"
	"github.com/arsenzairov/cs153/token"
)

type Tag int

const (
	FN_PROTO Tag = iota
	FN_DECL
)

type Node struct {
	tag   Tag // enum of possible AST node types. fn_decl for function declaration. integer_literal for literal integer. builin_call for calling builtin
	token token.Token
}

type parseFunc func()

type ParseRule struct {
	prefix     parseFunc
	infix      parseFunc
	precedence token.Precedence
}

type Parser struct {
	s         *scanner.Scanner
	peekToken token.Token
	curToken  token.Token
	rules     map[token.Tag]*ParseRule
	errors    []string
}

func NewParser(scanner *scanner.Scanner) *Parser {
	p := &Parser{s: scanner}

	p.registerRule(token.SUB, &ParseRule{p.unary, nil, token.PREC_OPER_LOW})

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerRule(tok token.Tag, rule *ParseRule) {
	p.rules[tok] = rule
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.s.Scan()
	if p.peekToken.Tag == token.ILLEGAL {
		p.errorAt(p.peekToken, "invalid token")
	}
}

func (p *Parser) peekIs(tok token.Tag) bool {
	if p.peekToken.Tag != tok {
		return false
	}
	return true
}

func (p *Parser) curIs(tok token.Tag) bool {
	if p.curToken.Tag != tok {
		return false
	}
	return true
}

func (p *Parser) expectPeek(tok token.Tag) bool {
	if !p.peekIs(tok) {
		p.errorAt(p.peekToken, fmt.Sprintf("Expected peek=%q, got=%q", tok, p.peekToken.Tag))
		return false
	}

	p.nextToken()
	return true
}

func (p *Parser) errorAt(tok token.Token, message string) error {
	return fmt.Errorf("[line %d] Error at %d: %s\n", tok.Pos.Line, tok.Pos.Column, message)
}

func (p *Parser) unary() {
	operator := p.curToken
	p.nextToken()

	p.parsePrecedence(PREC_UNARY)

	switch operator.Type {
	case token.MINUS:
		// p.emitByte()
	default:
		return
	}
}

func (p *Parser) binary() {
	operator := p.curToken
	p.nextToken()

	precedence := p.rules[operator.Tag].precedence
	p.parsePrecedence(precedence + 1)

	switch operator.Tag {
	case token.PLUS:
		//
	default:
		return
	}
}

func (p *Parser) parsePrecedence(precedence Precedence) error {
	prefixRule := p.rules[p.curToken.Tag].prefix
	if prefixRule == nil {
		return p.errorAt(p.curToken, "Expect expression.")
	}

	prefixRule()

	for precedence <= p.rules[p.peekToken.Tag].precedence && p.peekToken.Tag != token.EOF {
		p.nextToken()
		infixRule := p.rules[p.peekToken.Tag].prefix
		infixRule()
	}

	return nil
}

func (p *Parser) Parse() error {
	return p.parsePrecedence(PREC_NONE)
}
