package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/matiasinsaurralde/sl/ast"
	"github.com/matiasinsaurralde/sl/token"
)

type Parser struct {
	lexer     *token.Lexer
	curToken  token.TokenInfo
	peekToken token.TokenInfo
	errors    []string
}

func NewParser(input io.Reader) *Parser {
	lexer := token.NewLexer(input)
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) curTokenIs(t token.Token) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Token) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Token) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Token) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func Parse(input io.Reader) (*ast.File, error) {
	parser := NewParser(input)
	file := parser.parseFile()

	if len(parser.errors) > 0 {
		return nil, fmt.Errorf("parsing errors: %s", strings.Join(parser.errors, "; "))
	}

	return file, nil
}

func (p *Parser) parseFile() *ast.File {
	file := &ast.File{
		Nodes:    []ast.Node{},
		Comments: []ast.Comment{},
	}

	for p.curToken.Type != token.EOF {
		switch p.curToken.Type {
		case token.COMMENT:
			comment := p.parseComment()
			file.Comments = append(file.Comments, *comment)
		case token.EOL:
			p.nextToken() // skip EOL
		case token.PROGRAMA:
			program := p.parseProgram()
			file.Nodes = append(file.Nodes, program)
		case token.VAR:
			declarations := p.parseVariableDeclarations()
			for _, decl := range declarations {
				file.Nodes = append(file.Nodes, decl)
			}
		case token.SUBR:
			subroutine := p.parseSubroutineDeclaration()
			file.Nodes = append(file.Nodes, subroutine)
		case token.INICIO:
			block := p.parseBlockStatement()
			file.Nodes = append(file.Nodes, block)
		default:
			p.nextToken()
		}
	}

	return file
}

func (p *Parser) parseComment() *ast.Comment {
	comment := &ast.Comment{
		Text:     p.curToken.Literal,
		StartPos: p.curToken.Pos,
		EndPos:   p.curToken.Pos + token.Pos(len(p.curToken.Literal)),
	}
	p.nextToken()
	return comment
}

func (p *Parser) parseProgram() *ast.Program {
	program := &ast.Program{
		StartPos: p.curToken.Pos,
	}

	p.nextToken() // consume 'programa'

	if p.curTokenIs(token.IDENT) {
		program.Name = p.curToken.Literal
		p.nextToken()
	}

	program.EndPos = p.curToken.Pos
	return program
}

func (p *Parser) parseVariableDeclarations() []ast.Node {
	var declarations []ast.Node

	p.nextToken() // consume 'var'

	for !p.curTokenIs(token.INICIO) && !p.curTokenIs(token.SUBR) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.EOL) {
			p.nextToken()
			continue
		}

		decl := p.parseVariableDeclaration()
		if decl != nil {
			declarations = append(declarations, decl)
		}
	}

	return declarations
}

func (p *Parser) parseVariableDeclaration() ast.Node {
	startPos := p.curToken.Pos
	name := p.curToken.Literal

	if !p.curTokenIs(token.IDENT) {
		p.nextToken()
		return nil
	}

	p.nextToken()

	decl := &ast.VariableDeclaration{
		Name:     name,
		StartPos: startPos,
	}

	if p.curTokenIs(token.COLON) {
		// Type declaration: var name : type
		p.nextToken()
		if p.curTokenIs(token.NUMERICO) {
			decl.Type = p.curToken.Literal
			p.nextToken()
		}
	} else if p.curTokenIs(token.ASSIGN) {
		// Value assignment: var name = value
		p.nextToken()
		decl.Value = p.parseExpression()
	}

	decl.EndPos = p.curToken.Pos
	return decl
}

func (p *Parser) parseSubroutineDeclaration() *ast.SubroutineDeclaration {
	sub := &ast.SubroutineDeclaration{
		StartPos: p.curToken.Pos,
	}

	p.nextToken() // consume 'sub'

	if p.curTokenIs(token.IDENT) {
		sub.Name = p.curToken.Literal
		p.nextToken()
	}

	// Parse parameters
	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		sub.Parameters = p.parseParameters()
		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
	}

	// Parse return type
	if p.curTokenIs(token.RETORNA) {
		p.nextToken()
		if p.curTokenIs(token.NUMERICO) {
			sub.ReturnType = p.curToken.Literal
			p.nextToken()
		}
	}

	// Parse body
	if p.curTokenIs(token.INICIO) {
		sub.Body = p.parseBlockStatement()
	}

	sub.EndPos = p.curToken.Pos
	return sub
}

func (p *Parser) parseParameters() []*ast.Parameter {
	var params []*ast.Parameter

	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.IDENT) {
			param := &ast.Parameter{
				Name:     p.curToken.Literal,
				StartPos: p.curToken.Pos,
			}
			p.nextToken()

			if p.curTokenIs(token.COLON) {
				p.nextToken()
				if p.curTokenIs(token.NUMERICO) {
					param.Type = p.curToken.Literal
					p.nextToken()
				}
			}

			param.EndPos = p.curToken.Pos
			params = append(params, param)

			if p.curTokenIs(token.COMMA) {
				p.nextToken()
			}
		} else {
			p.nextToken()
		}
	}

	return params
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		StartPos:   p.curToken.Pos,
		Statements: []ast.Statement{},
	}

	p.nextToken() // consume 'inicio' or '{'

	stmtIdx := 0
	for !p.curTokenIs(token.FIN) && !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.EOL) {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
			stmtIdx++
		}
	}

	if p.curTokenIs(token.FIN) || p.curTokenIs(token.RBRACE) {
		p.nextToken()
	}

	block.EndPos = p.curToken.Pos
	return block
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.SI:
		return p.parseIfStatement()
	case token.MIENTRAS:
		return p.parseWhileStatement()
	case token.DESDE:
		return p.parseForStatement()
	case token.RETORNA:
		return p.parseReturnStatement()
	case token.TERMINAR:
		return p.parseTerminateStatement()
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			return p.parseAssignmentStatement()
		}
		return p.parseExpressionStatement()
	case token.ASSIGN:
		// Handle case where we're already at the '=' token
		// This can happen if the parser advanced incorrectly
		return p.parseExpressionStatement()
	case token.IMPRIMIR, token.LEER:
		return p.parseExpressionStatement()
	default:
		p.nextToken()
		return nil
	}
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	ifStmt := &ast.IfStatement{
		StartPos: p.curToken.Pos,
	}

	p.nextToken() // consume 'si'

	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		ifStmt.Condition = p.parseExpression()
		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
	}

	if p.curTokenIs(token.LBRACE) {
		p.nextToken()
		ifStmt.Then = p.parseBlockStatement()
	} else {
		ifStmt.Then = p.parseStatement()
	}

	if p.curTokenIs(token.SINO) {
		p.nextToken()
		if p.curTokenIs(token.LBRACE) {
			p.nextToken()
			ifStmt.Else = p.parseBlockStatement()
		} else {
			ifStmt.Else = p.parseStatement()
		}
	}

	ifStmt.EndPos = p.curToken.Pos
	return ifStmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	whileStmt := &ast.WhileStatement{
		StartPos: p.curToken.Pos,
	}

	p.nextToken() // consume 'mientras'

	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		whileStmt.Condition = p.parseExpression()
		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
	}

	if p.curTokenIs(token.LBRACE) {
		p.nextToken()
		whileStmt.Body = p.parseBlockStatement()
	} else {
		whileStmt.Body = p.parseStatement()
	}

	whileStmt.EndPos = p.curToken.Pos
	return whileStmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	forStmt := &ast.ForStatement{
		StartPos: p.curToken.Pos,
	}

	p.nextToken() // consume 'desde'

	if p.curTokenIs(token.IDENT) {
		forStmt.Variable = p.curToken.Literal
		p.nextToken()
	}

	if p.curTokenIs(token.ASSIGN) {
		p.nextToken()
		forStmt.Start = p.parseExpression()
	}

	if p.curTokenIs(token.HASTA) {
		p.nextToken()
		forStmt.EndExpr = p.parseExpression()
	}

	if p.curTokenIs(token.PASO) {
		p.nextToken()
		forStmt.Step = p.parseExpression()
	}

	if p.curTokenIs(token.LBRACE) {
		p.nextToken()
		forStmt.Body = p.parseBlockStatement()
	} else {
		forStmt.Body = p.parseStatement()
	}

	forStmt.EndPos = p.curToken.Pos
	return forStmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{
		StartPos: p.curToken.Pos,
	}

	p.nextToken() // consume 'retorna'

	if !p.curTokenIs(token.EOL) && !p.curTokenIs(token.EOF) {
		returnStmt.Value = p.parseExpression()
	}

	returnStmt.EndPos = p.curToken.Pos
	return returnStmt
}

func (p *Parser) parseTerminateStatement() *ast.TerminateStatement {
	terminateStmt := &ast.TerminateStatement{
		StartPos: p.curToken.Pos,
	}

	p.nextToken() // consume 'terminar'

	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		terminateStmt.Message = p.parseExpression()
		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
	}

	terminateStmt.EndPos = p.curToken.Pos
	return terminateStmt
}

func (p *Parser) parseAssignmentStatement() *ast.ExpressionStatement {
	assignment := &ast.AssignmentExpression{
		Left: &ast.Identifier{
			Name:     p.curToken.Literal,
			StartPos: p.curToken.Pos,
		},
		StartPos: p.curToken.Pos,
	}

	p.nextToken() // consume identifier

	if p.curTokenIs(token.ASSIGN) {
		assignment.Operator = p.curToken.Literal
		p.nextToken() // consume '='
		assignment.Right = p.parseExpression()
	}

	assignment.EndPos = assignment.Right.End()
	assignment.Left.EndPos = assignment.EndPos

	return &ast.ExpressionStatement{
		Expression: assignment,
		StartPos:   assignment.StartPos,
		EndPos:     assignment.EndPos,
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Expression: p.parseExpression(),
		StartPos:   p.curToken.Pos,
	}
	if stmt.Expression != nil {
		stmt.EndPos = stmt.Expression.End()
	}
	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseBinaryExpression(0)
}

func (p *Parser) parseBinaryExpression(precedence int) ast.Expression {
	left := p.parseUnaryExpression()

	for !p.curTokenIs(token.EOL) && !p.curTokenIs(token.EOF) && precedence < p.getPrecedence(p.curToken.Type) {
		// Don't treat assignment as a binary operator
		if p.curTokenIs(token.ASSIGN) {
			break
		}

		operator := p.curToken.Literal
		p.nextToken()

		right := p.parseBinaryExpression(p.getPrecedence(token.Lookup(operator)) + 1)

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
			StartPos: left.Pos(),
			EndPos:   right.End(),
		}
	}

	return left
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	switch p.curToken.Type {
	case token.PLUS, token.MINUS:
		operator := p.curToken.Literal
		startPos := p.curToken.Pos
		p.nextToken()

		operand := p.parseUnaryExpression()

		return &ast.BinaryExpression{
			Left: &ast.Literal{
				Type:     token.INT,
				Value:    "0",
				StartPos: startPos,
				EndPos:   startPos,
			},
			Operator: operator,
			Right:    operand,
			StartPos: startPos,
			EndPos:   operand.End(),
		}
	default:
		return p.parsePrimaryExpression()
	}
}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	switch p.curToken.Type {
	case token.IDENT, token.IMPRIMIR, token.LEER:
		if p.peekTokenIs(token.LPAREN) {
			return p.parseCallExpression()
		}
		if p.curToken.Type == token.IDENT {
			return p.parseIdentifier()
		}
		// For imprimir/leer without parens, treat as nil
		return nil
	case token.INT, token.FLOAT:
		return p.parseLiteral()
	case token.STRING:
		return p.parseLiteral()
	case token.LPAREN:
		p.nextToken()
		expr := p.parseExpression()
		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
		return expr
	case token.IFVAL:
		return p.parseIfValExpression()
	case token.ASSIGN:
		// Handle assignment token - this should not happen in normal parsing
		// but if it does, we need to handle it gracefully
		p.nextToken()
		return nil
	default:
		p.nextToken()
		return nil
	}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	ident := &ast.Identifier{
		Name:     p.curToken.Literal,
		StartPos: p.curToken.Pos,
		EndPos:   p.curToken.Pos + token.Pos(len(p.curToken.Literal)),
	}
	p.nextToken()
	return ident
}

func (p *Parser) parseLiteral() *ast.Literal {
	lit := &ast.Literal{
		Type:     p.curToken.Type,
		Value:    p.curToken.Literal,
		StartPos: p.curToken.Pos,
		EndPos:   p.curToken.Pos + token.Pos(len(p.curToken.Literal)),
	}
	p.nextToken()
	return lit
}

func (p *Parser) parseCallExpression() *ast.CallExpression {
	call := &ast.CallExpression{
		Function:  p.curToken.Literal,
		StartPos:  p.curToken.Pos,
		Arguments: []ast.Expression{},
	}

	p.nextToken() // consume function name
	p.nextToken() // consume '('

	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		arg := p.parseExpression()
		if arg != nil {
			call.Arguments = append(call.Arguments, arg)
		}

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
	}

	call.EndPos = p.curToken.Pos
	return call
}

func (p *Parser) parseIfValExpression() *ast.IfValExpression {
	ifVal := &ast.IfValExpression{
		StartPos: p.curToken.Pos,
	}

	p.nextToken() // consume 'ifval'

	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		ifVal.Condition = p.parseExpression()

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
			ifVal.Then = p.parseExpression()

			if p.curTokenIs(token.COMMA) {
				p.nextToken()
				ifVal.Else = p.parseExpression()
			}
		}

		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
	}

	ifVal.EndPos = p.curToken.Pos
	return ifVal
}

func (p *Parser) getPrecedence(t token.Token) int {
	switch t {
	case token.MULTIPLY, token.DIVIDE, token.MODULO:
		return 4
	case token.PLUS, token.MINUS:
		return 3
	case token.EQ, token.NEQ, token.LT, token.LTE, token.GT, token.GTE:
		return 2
	case token.AND, token.OR:
		return 1
	case token.ASSIGN:
		return 0 // Assignment has lowest precedence
	default:
		return -1
	}
}
