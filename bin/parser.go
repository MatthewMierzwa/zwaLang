package bin

import (
	"fmt"
)

// AST Node types
type Node interface{}

type VarDeclaration struct {
	Name  string
	Value int
}

type Assignment struct {
	VarName    string
	Expression Node
}

type PrintStatement struct {
	Expression Expression
}

type Expression interface {
	Node
}

type NumberLiteral struct {
	Value int
}

type Variable struct {
	Name string
}

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

// Parser structure
type Parser struct {
	lex       *Lexer
	curToken  Token
	peekToken Token
}

// NewParser initializes a parser with the given lexer
func NewParser(lex *Lexer) *Parser {
	p := &Parser{lex: lex}
	p.nextToken()
	p.nextToken()
	return p
}

// -- Receiver Functions --

// NextToken advances the parser to the next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekTokenIsOperator() bool {
	return p.peekToken.Type == PLUS || p.peekToken.Type == MINUS || p.peekToken.Type == ASTERISK || p.peekToken.Type == SLASH || p.peekToken.Type == MODULO
}

func (p *Parser) ParseProgram() []Node {
	var nodes []Node
	for p.curToken.Type != EOF {
		node := p.parseStatement()
		if node != nil {
			nodes = append(nodes, node)
		}
		p.nextToken()
	}
	return nodes
}

func (p *Parser) parseStatement() Node {
	switch p.curToken.Type {
	case IDENT:
		return p.parseAssignmentOrDeclaration()
	case SHOW:
		return p.parsePrintStatement()
	default:
		fmt.Printf("Unknown statement: %+v\n", p.curToken.Literal)
		return nil
	}
}

func (p *Parser) parseAssignmentOrDeclaration() Node {
	identifier := p.curToken.Literal
	p.nextToken()

	if p.curToken.Type == COLON {
		// Declaration
		p.nextToken()
		if p.curToken.Type != NATURAL_NUMBER_TYPE {
			fmt.Printf("Expected type of variable (Currently only supports type: number) after colon: %+v\n", p.curToken.Literal)
			return nil
		}
		p.nextToken()
		if p.curToken.Type != ASSIGN {
			fmt.Printf("Expected assignment operator after type declaration: %+v\n", p.curToken.Literal)
			return nil
		}
		p.nextToken()
		value := p.curToken.Literal
		return &VarDeclaration{Name: identifier, Value: atoi(value)}
	} else if p.curToken.Type == ASSIGN {
		p.nextToken()
		expression := p.parseExpression()
		return &Assignment{VarName: identifier, Expression: expression}
	}

	return nil
}

// parseExpression handles parsing addition and subtraction
func (p *Parser) parseExpression() Expression {
	left := p.parseTerm()

	for p.curToken.Type == PLUS || p.curToken.Type == MINUS {
		operator := p.curToken.Literal
		p.nextToken()
		right := p.parseTerm()
		left = &BinaryExpression{Left: left, Operator: operator, Right: right}
	}

	return left
}

// parseTerm handles parsing multiplication and division
func (p *Parser) parseTerm() Expression {
	left := p.parseFactor()

	for p.curToken.Type == ASTERISK || p.curToken.Type == SLASH || p.curToken.Type == MODULO {
		operator := p.curToken.Literal
		p.nextToken()
		right := p.parseFactor()
		left = &BinaryExpression{Left: left, Operator: operator, Right: right}
	}

	return left
}

// parseFactor handles parsing numbers, variables, and parentheses
func (p *Parser) parseFactor() Expression {
	switch p.curToken.Type {
	case LPAREN:
		p.nextToken() // Skip '('
		expr := p.parseExpression()
		p.nextToken()
		if p.curToken.Type != RPAREN {
			fmt.Println("Expected )")
			return nil
		}
		p.nextToken() // Skip ')'
		return expr
	case IDENT:
		ident := p.curToken.Literal
		// if the next token is a operator, then we need to parse the next term
		if p.peekTokenIsOperator() {
			p.nextToken()
		}
		return &Variable{Name: ident}
	case NUMBER:
		value := atoi(p.curToken.Literal)
		// if the next token is a operator, then we need to parse the next term
		if p.peekTokenIsOperator() {
			p.nextToken()
		}
		return &NumberLiteral{Value: value}
	default:
		fmt.Printf("Unexpected token: %+v\n", p.curToken.Literal)
		return nil
	}
}

func (p *Parser) parsePrintStatement() Node {
	p.nextToken()
	return &PrintStatement{Expression: p.parseExpression()}
}

// -- Helper Functions --

func atoi(str string) int {
	var result int
	fmt.Sscanf(str, "%d", &result)
	return result
}