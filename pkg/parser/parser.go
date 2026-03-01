// Package parser implements a recursive-descent parser for the SL language.
package parser

import (
	"fmt"
	"strconv"

	"github.com/matiasinsaurralde/sl/pkg/ast"
	"github.com/matiasinsaurralde/sl/pkg/lexer"
)

// ParseError records a parse error.
type ParseError struct {
	Line int
	Col  int
	Msg  string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d col %d: %s", e.Line, e.Col, e.Msg)
}

// Parser holds the parser state.
type Parser struct {
	lex        *lexer.Lexer
	errors     []*ParseError
	parenDepth int // tracks nesting depth of '(' — used to resolve '=' ambiguity
}

// New creates a Parser for src.
func New(src string) *Parser {
	return &Parser{lex: lexer.New(src)}
}

func (p *Parser) peek() lexer.Token {
	return p.lex.Peek()
}

func (p *Parser) next() lexer.Token {
	return p.lex.Next()
}

func (p *Parser) expect(tt lexer.TokenType) (lexer.Token, error) {
	t := p.next()
	if t.Type == lexer.LPAREN {
		p.parenDepth++
	} else if t.Type == lexer.RPAREN && p.parenDepth > 0 {
		p.parenDepth--
	}
	if t.Type != tt {
		return t, p.errorf(t, "expected %s, got %q", tt, t.Literal)
	}
	return t, nil
}

func (p *Parser) eat(tt lexer.TokenType) bool {
	if p.peek().Type == tt {
		p.next()
		switch tt {
		case lexer.LPAREN:
			p.parenDepth++
		case lexer.RPAREN:
			if p.parenDepth > 0 {
				p.parenDepth--
			}
		}
		return true
	}
	return false
}

func (p *Parser) errorf(t lexer.Token, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	e := &ParseError{Line: t.Line, Col: t.Col, Msg: msg}
	p.errors = append(p.errors, e)
	return e
}

// Errors returns all parse errors.
func (p *Parser) Errors() []*ParseError {
	return p.errors
}

// ParseProgram parses an entire SL program.
func (p *Parser) ParseProgram() (*ast.Program, error) {
	prog := &ast.Program{}

	// Optional "programa <name>"
	if p.peek().Type == lexer.PROGRAMA {
		p.next()
		if p.peek().Type == lexer.IDENT {
			prog.Name = p.next().Literal
		}
	}

	// Global declarations (var/const/tipos in any order, multiple times)
	p.parseDecls(&prog.Consts, &prog.Types, &prog.Vars)

	// Main body: inicio ... fin
	if _, err := p.expect(lexer.INICIO); err != nil {
		return prog, err
	}
	prog.Body = p.parseStmtList(lexer.FIN)
	if _, err := p.expect(lexer.FIN); err != nil {
		return prog, err
	}

	// Subroutine definitions
	for p.peek().Type == lexer.SUB || p.peek().Type == lexer.SUBRUTINA {
		sub, err := p.parseSub()
		if err != nil {
			return prog, err
		}
		prog.Subs = append(prog.Subs, sub)
	}

	return prog, nil
}

// parseDecls parses zero or more var/const/tipos sections.
func (p *Parser) parseDecls(consts *[]*ast.ConstDecl, types *[]*ast.TiposDecl, vars *[]*ast.VarDecl) {
	for {
		switch p.peek().Type {
		case lexer.CONST, lexer.CONSTANTES:
			p.next()
			*consts = append(*consts, p.parseConstSection()...)
		case lexer.TIPOS:
			p.next()
			*types = append(*types, p.parseTiposSection()...)
		case lexer.VAR, lexer.VARIABLES:
			p.next()
			*vars = append(*vars, p.parseVarSection()...)
		default:
			return
		}
	}
}

// parseConstSection parses constant declarations until a non-const token.
func (p *Parser) parseConstSection() []*ast.ConstDecl {
	var decls []*ast.ConstDecl
	for p.peek().Type == lexer.IDENT {
		name := p.next().Literal
		if _, err := p.expect(lexer.ASSIGN); err != nil {
			break
		}
		init := p.parseExpr()
		decls = append(decls, &ast.ConstDecl{Name: name, Init: init})
	}
	return decls
}

// parseTiposSection parses type alias declarations.
// SL syntax: name = type (using '=' not ':')
func (p *Parser) parseTiposSection() []*ast.TiposDecl {
	var decls []*ast.TiposDecl
	for p.peek().Type == lexer.IDENT {
		name := p.next().Literal
		if _, err := p.expect(lexer.ASSIGN); err != nil {
			break
		}
		typ := p.parseType()
		decls = append(decls, &ast.TiposDecl{Name: name, Type: typ})
	}
	return decls
}

// parseVarSection parses variable declarations.
func (p *Parser) parseVarSection() []*ast.VarDecl {
	var decls []*ast.VarDecl
	for isVarDeclStart(p.peek()) {
		d := p.parseOneVarDecl()
		if d != nil {
			decls = append(decls, d)
		}
	}
	return decls
}

// isVarDeclStart returns true if the next token could start a var declaration.
func isVarDeclStart(t lexer.Token) bool {
	return t.Type == lexer.IDENT
}

// parseOneVarDecl parses a single variable declaration line.
// Forms:
//
//	name = expr                    (inferred type, single var)
//	name : type                    (typed, zero value)
//	name : type = expr             (typed, initialized)
//	name, name, ... : type         (multiple vars)
func (p *Parser) parseOneVarDecl() *ast.VarDecl {
	if p.peek().Type != lexer.IDENT {
		return nil
	}

	// Collect names
	names := []string{p.next().Literal}
	for p.peek().Type == lexer.COMMA {
		p.next()
		if p.peek().Type == lexer.IDENT {
			names = append(names, p.next().Literal)
		}
	}

	// name = expr  (no type annotation)
	if p.peek().Type == lexer.ASSIGN && len(names) == 1 {
		p.next()
		init := p.parseExpr()
		return &ast.VarDecl{Names: names, Type: nil, Init: init}
	}

	// name : type [= expr]
	if _, err := p.expect(lexer.COLON); err != nil {
		return &ast.VarDecl{Names: names}
	}
	typ := p.parseType()
	var init ast.Expr
	if p.peek().Type == lexer.ASSIGN {
		p.next()
		init = p.parseExpr()
	}
	return &ast.VarDecl{Names: names, Type: typ, Init: init}
}

// parseType parses a type expression.
func (p *Parser) parseType() ast.TypeNode {
	switch p.peek().Type {
	case lexer.NUMERICO:
		p.next()
		return &ast.SimpleType{Name: "numerico"}
	case lexer.CADENA:
		p.next()
		return &ast.SimpleType{Name: "cadena"}
	case lexer.LOGICO:
		p.next()
		return &ast.SimpleType{Name: "logico"}
	case lexer.VECTOR:
		return p.parseVectorType()
	case lexer.MATRIZ:
		return p.parseMatrixType()
	case lexer.REGISTRO:
		return p.parseRegistroType()
	case lexer.IDENT:
		name := p.next().Literal
		return &ast.NamedType{Name: name}
	}
	t := p.peek()
	_ = p.errorf(t, "expected type, got %q", t.Literal)
	p.next()
	return &ast.SimpleType{Name: "numerico"}
}

func (p *Parser) parseVectorType() ast.TypeNode {
	p.next() // consume 'vector'
	if _, err := p.expect(lexer.LBRACK); err != nil {
		return &ast.VectorType{Size: 0, ElemType: &ast.SimpleType{Name: "numerico"}}
	}
	size := 0
	if p.peek().Type == lexer.STAR {
		p.next() // consume *
	} else if p.peek().Type == lexer.NUMBER {
		n, _ := strconv.Atoi(p.next().Literal)
		size = n
	}
	p.eat(lexer.RBRACK)
	elemType := ast.TypeNode(&ast.SimpleType{Name: "numerico"})
	if isTypeStart(p.peek()) {
		elemType = p.parseType()
	}
	return &ast.VectorType{Size: size, ElemType: elemType}
}

func (p *Parser) parseMatrixType() ast.TypeNode {
	p.next() // consume 'matriz'
	if _, err := p.expect(lexer.LBRACK); err != nil {
		return &ast.MatrixType{Dims: []int{0, 0}, ElemType: &ast.SimpleType{Name: "numerico"}}
	}
	var dims []int
	for {
		if p.peek().Type == lexer.STAR {
			p.next()
			dims = append(dims, 0)
		} else if p.peek().Type == lexer.NUMBER {
			n, _ := strconv.Atoi(p.next().Literal)
			dims = append(dims, n)
		} else {
			dims = append(dims, 0)
		}
		if p.peek().Type == lexer.COMMA {
			p.next()
		} else {
			break
		}
	}
	p.eat(lexer.RBRACK)
	elemType := ast.TypeNode(&ast.SimpleType{Name: "numerico"})
	if isTypeStart(p.peek()) {
		elemType = p.parseType()
	}
	return &ast.MatrixType{Dims: dims, ElemType: elemType}
}

// isTypeStart returns true if t can start a type expression.
func isTypeStart(t lexer.Token) bool {
	switch t.Type {
	case lexer.NUMERICO, lexer.CADENA, lexer.LOGICO,
		lexer.VECTOR, lexer.MATRIZ, lexer.REGISTRO, lexer.IDENT:
		return true
	}
	return false
}

func (p *Parser) parseRegistroType() ast.TypeNode {
	p.next() // consume 'registro'
	p.eat(lexer.LBRACE)
	var fields []*ast.FieldDef
	for p.peek().Type == lexer.IDENT {
		names := []string{p.next().Literal}
		for p.peek().Type == lexer.COMMA {
			p.next()
			if p.peek().Type == lexer.IDENT {
				names = append(names, p.next().Literal)
			}
		}
		p.eat(lexer.COLON)
		typ := p.parseType()
		fields = append(fields, &ast.FieldDef{Names: names, Type: typ})
		p.eat(lexer.SEMI)
	}
	p.eat(lexer.RBRACE)
	return &ast.RegistroType{Fields: fields}
}

// parseSub parses a subroutine declaration.
func (p *Parser) parseSub() (*ast.SubDecl, error) {
	p.next() // consume 'sub' or 'subrutina'
	nameTok, err := p.expect(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	sub := &ast.SubDecl{Name: nameTok.Literal}

	// Parameters
	if p.peek().Type == lexer.LPAREN {
		p.next()
		sub.Params = p.parseParams()
		p.eat(lexer.RPAREN)
	}

	// Optional return type
	if p.peek().Type == lexer.RETORNA {
		p.next()
		sub.ReturnType = p.parseType()
	}

	// Local declarations
	p.parseDecls(&sub.Consts, &sub.Types, &sub.Vars)

	// Body
	if _, err := p.expect(lexer.INICIO); err != nil {
		return sub, err
	}
	sub.Body = p.parseStmtList(lexer.FIN)
	p.eat(lexer.FIN)

	return sub, nil
}

// parseParams parses a parameter list inside ().
// Groups are separated by ';' or by starting a new ref keyword.
func (p *Parser) parseParams() []*ast.ParamGroup {
	var groups []*ast.ParamGroup
	for p.peek().Type != lexer.RPAREN && p.peek().Type != lexer.EOF {
		g := p.parseParamGroup()
		if g != nil {
			groups = append(groups, g)
		}
		if p.peek().Type == lexer.SEMI {
			p.next()
		} else if p.peek().Type != lexer.RPAREN && p.peek().Type != lexer.EOF {
			// Allow implicit continuation without semicolon
			// if the next token is ref or ident (another group)
			if p.peek().Type != lexer.REF && p.peek().Type != lexer.IDENT {
				break
			}
		}
	}
	return groups
}

func (p *Parser) parseParamGroup() *ast.ParamGroup {
	byRef := false
	if p.peek().Type == lexer.REF {
		p.next()
		byRef = true
	}
	if p.peek().Type != lexer.IDENT {
		return nil
	}
	names := []string{p.next().Literal}
	for p.peek().Type == lexer.COMMA {
		p.next()
		if p.peek().Type == lexer.IDENT {
			names = append(names, p.next().Literal)
		}
	}
	if _, err := p.expect(lexer.COLON); err != nil {
		return &ast.ParamGroup{ByRef: byRef, Names: names}
	}
	typ := p.parseType()
	return &ast.ParamGroup{ByRef: byRef, Names: names, Type: typ}
}

// parseStmtList parses statements until stopAt token or EOF.
func (p *Parser) parseStmtList(stopAt ...lexer.TokenType) []ast.Stmt {
	var stmts []ast.Stmt
	for {
		t := p.peek()
		if t.Type == lexer.EOF {
			break
		}
		// Check stop tokens
		for _, stop := range stopAt {
			if t.Type == stop {
				return stmts
			}
		}
		// Also stop on sino (for si parsing)
		if t.Type == lexer.SINO {
			return stmts
		}

		stmt := p.parseStmt()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
		// Consume optional semicolons
		for p.peek().Type == lexer.SEMI {
			p.next()
		}
	}
	return stmts
}

// parseStmt parses a single statement.
func (p *Parser) parseStmt() ast.Stmt {
	t := p.peek()
	switch t.Type {
	case lexer.SI:
		return p.parseSi()
	case lexer.DESDE:
		return p.parseDesde()
	case lexer.MIENTRAS:
		return p.parseMientras()
	case lexer.REPETIR:
		return p.parseRepetir()
	case lexer.EVAL:
		return p.parseEval()
	case lexer.SALIR:
		p.next()
		return &ast.SalirStmt{Line: t.Line}
	case lexer.RETORNA:
		return p.parseRetorna()
	case lexer.SEMI:
		p.next()
		return nil
	}

	// Expression statement: assignment or call
	return p.parseExprStmt()
}

func (p *Parser) parseSi() ast.Stmt {
	line := p.peek().Line
	p.next() // consume 'si'
	p.eat(lexer.LPAREN)
	cond := p.parseExpr()
	p.eat(lexer.RPAREN)
	p.eat(lexer.LBRACE)

	then := p.parseStmtList(lexer.RBRACE, lexer.FIN)
	var elseIfs []*ast.ElseIf
	var elsePart []ast.Stmt

	for p.peek().Type == lexer.SINO {
		p.next() // consume 'sino'
		if p.peek().Type == lexer.SI {
			// sino si (cond) ...
			p.next() // consume 'si'
			p.eat(lexer.LPAREN)
			eicond := p.parseExpr()
			p.eat(lexer.RPAREN)
			eiBody := p.parseStmtList(lexer.RBRACE, lexer.SINO, lexer.FIN)
			elseIfs = append(elseIfs, &ast.ElseIf{Cond: eicond, Body: eiBody})
		} else {
			// plain sino
			elsePart = p.parseStmtList(lexer.RBRACE, lexer.FIN)
			break
		}
	}

	p.eat(lexer.RBRACE)
	return &ast.SiStmt{
		Line:    line,
		Cond:    cond,
		Then:    then,
		ElseIfs: elseIfs,
		Else:    elsePart,
	}
}

func (p *Parser) parseDesde() ast.Stmt {
	line := p.peek().Line
	p.next() // consume 'desde'
	varName := p.next().Literal
	p.eat(lexer.ASSIGN)
	start := p.parseExpr()
	if _, err := p.expect(lexer.HASTA); err != nil {
		return &ast.DesdeStmt{Line: line, Var: varName, Start: start, End: start}
	}
	end := p.parseExpr()
	var step ast.Expr
	if p.peek().Type == lexer.PASO {
		p.next()
		step = p.parseExpr()
	}
	p.eat(lexer.LBRACE)
	body := p.parseStmtList(lexer.RBRACE)
	p.eat(lexer.RBRACE)
	return &ast.DesdeStmt{Line: line, Var: varName, Start: start, End: end, Step: step, Body: body}
}

func (p *Parser) parseMientras() ast.Stmt {
	line := p.peek().Line
	p.next() // consume 'mientras'
	p.eat(lexer.LPAREN)
	cond := p.parseExpr()
	p.eat(lexer.RPAREN)
	p.eat(lexer.LBRACE)
	body := p.parseStmtList(lexer.RBRACE)
	p.eat(lexer.RBRACE)
	return &ast.MientrasStmt{Line: line, Cond: cond, Body: body}
}

func (p *Parser) parseRepetir() ast.Stmt {
	line := p.peek().Line
	p.next() // consume 'repetir'
	var body []ast.Stmt
	if p.peek().Type == lexer.LBRACE {
		p.next() // consume '{'
		body = p.parseStmtList(lexer.RBRACE)
		p.eat(lexer.RBRACE)
	} else {
		body = p.parseStmtList(lexer.HASTA)
	}
	p.eat(lexer.HASTA)
	p.eat(lexer.LPAREN)
	cond := p.parseExpr()
	p.eat(lexer.RPAREN)
	return &ast.RepetirStmt{Line: line, Body: body, Cond: cond}
}

func (p *Parser) parseEval() ast.Stmt {
	line := p.peek().Line
	p.next() // consume 'eval'
	p.eat(lexer.LBRACE)
	stmt := &ast.EvalStmt{Line: line}
	for p.peek().Type == lexer.CASO {
		p.next() // consume 'caso'
		p.eat(lexer.LPAREN)
		cond := p.parseExpr()
		p.eat(lexer.RPAREN)
		p.eat(lexer.COLON) // optional ':' after caso condition
		body := p.parseStmtList(lexer.CASO, lexer.SINO, lexer.RBRACE)
		stmt.Cases = append(stmt.Cases, &ast.EvalCase{Cond: cond, Body: body})
	}
	if p.peek().Type == lexer.SINO {
		p.next()
		p.eat(lexer.COLON) // optional ':' after sino
		stmt.Else = p.parseStmtList(lexer.RBRACE)
	}
	p.eat(lexer.RBRACE)
	return stmt
}

func (p *Parser) parseRetorna() ast.Stmt {
	line := p.peek().Line
	p.next() // consume 'retorna'
	var val ast.Expr
	// Parens are optional; retorna can appear with or without them.
	if p.peek().Type == lexer.LPAREN {
		p.next()
		val = p.parseExpr()
		p.eat(lexer.RPAREN)
	} else if !isStmtTerminator(p.peek()) {
		// No paren: parse bare expression (e.g. "retorna a")
		val = p.parseExpr()
	}
	return &ast.RetornaStmt{Line: line, Value: val}
}

// isStmtTerminator returns true for tokens that cannot start an expression
// and therefore signal the end of a retorna statement.
func isStmtTerminator(t lexer.Token) bool {
	switch t.Type {
	case lexer.FIN, lexer.SINO, lexer.CASO, lexer.HASTA,
		lexer.EOF, lexer.RBRACE, lexer.SEMI:
		return true
	}
	return false
}

// parseExprStmt parses an assignment or a call statement.
func (p *Parser) parseExprStmt() ast.Stmt {
	expr := p.parseExpr()
	if expr == nil {
		// skip unknown token
		t := p.next()
		_ = p.errorf(t, "unexpected token %q", t.Literal)
		return nil
	}

	// Check if it's an assignment
	if p.peek().Type == lexer.ASSIGN {
		line := p.peek().Line
		p.next()
		val := p.parseExpr()
		return &ast.AssignStmt{Line: line, Target: expr, Value: val}
	}

	// Must be a call statement or terminar/imprimir/leer handled via expression
	switch e := expr.(type) {
	case *ast.CallExpr:
		// Check for special built-in statements
		switch e.Name {
		case "imprimir":
			return &ast.ImprimirStmt{Line: e.Line, Args: e.Args}
		case "leer":
			return &ast.LeerStmt{Line: e.Line, Vars: e.Args}
		case "terminar":
			s := &ast.TerminarStmt{Line: e.Line}
			if len(e.Args) > 0 {
				s.Msg = e.Args[0]
			}
			return s
		}
		return &ast.CallStmt{Line: e.Line, Call: e}
	}

	return &ast.AssignStmt{Line: expr.GetLine(), Target: expr, Value: expr}
}

// ---- Expression parsing (recursive descent, precedence climbing) ----

// parseExpr parses the lowest-precedence expression (or).
func (p *Parser) parseExpr() ast.Expr {
	return p.parseOr()
}

func (p *Parser) parseOr() ast.Expr {
	left := p.parseAnd()
	for p.peek().Type == lexer.OR || p.peek().Type == lexer.OR2 {
		line := p.peek().Line
		p.next()
		right := p.parseAnd()
		left = &ast.BinaryExpr{Line: line, Op: lexer.OR, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseAnd() ast.Expr {
	left := p.parseNot()
	for p.peek().Type == lexer.AND || p.peek().Type == lexer.AND2 {
		line := p.peek().Line
		p.next()
		right := p.parseNot()
		left = &ast.BinaryExpr{Line: line, Op: lexer.AND, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseNot() ast.Expr {
	if p.peek().Type == lexer.NOT {
		line := p.peek().Line
		p.next()
		operand := p.parseRelational()
		return &ast.UnaryExpr{Line: line, Op: lexer.NOT, Operand: operand}
	}
	return p.parseRelational()
}

var relOps = map[lexer.TokenType]bool{
	lexer.EQ: true, lexer.NEQ: true, lexer.LT: true,
	lexer.LE: true, lexer.GT: true, lexer.GE: true,
}

func (p *Parser) parseRelational() ast.Expr {
	left := p.parseAddSub()
	t := p.peek()
	op := t.Type
	// SL uses bare '=' for equality when inside parentheses (e.g. si (x = 1)).
	// At statement level '=' is assignment and handled separately.
	if op == lexer.ASSIGN && p.parenDepth > 0 {
		op = lexer.EQ
	}
	if relOps[op] {
		line := t.Line
		p.next()
		right := p.parseAddSub()
		return &ast.BinaryExpr{Line: line, Op: op, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseAddSub() ast.Expr {
	left := p.parseMulDiv()
	for p.peek().Type == lexer.PLUS || p.peek().Type == lexer.MINUS {
		line := p.peek().Line
		op := p.next().Type
		right := p.parseMulDiv()
		left = &ast.BinaryExpr{Line: line, Op: op, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseMulDiv() ast.Expr {
	left := p.parseUnary()
	for p.peek().Type == lexer.STAR || p.peek().Type == lexer.SLASH || p.peek().Type == lexer.PERCENT {
		line := p.peek().Line
		op := p.next().Type
		right := p.parseUnary()
		left = &ast.BinaryExpr{Line: line, Op: op, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseUnary() ast.Expr {
	if p.peek().Type == lexer.MINUS || p.peek().Type == lexer.PLUS {
		line := p.peek().Line
		op := p.next().Type
		operand := p.parsePow()
		return &ast.UnaryExpr{Line: line, Op: op, Operand: operand}
	}
	return p.parsePow()
}

func (p *Parser) parsePow() ast.Expr {
	base := p.parseFactor()
	if p.peek().Type == lexer.CARET {
		line := p.peek().Line
		p.next()
		// right-associative: recurse into parsePow
		exp := p.parsePow()
		return &ast.BinaryExpr{Line: line, Op: lexer.CARET, Left: base, Right: exp}
	}
	return base
}

func (p *Parser) parseFactor() ast.Expr {
	t := p.peek()
	var expr ast.Expr

	switch t.Type {
	case lexer.NUMBER:
		p.next()
		v, _ := strconv.ParseFloat(t.Literal, 64)
		expr = &ast.NumberLit{Line: t.Line, Value: v}

	case lexer.STRING:
		p.next()
		expr = &ast.StringLit{Line: t.Line, Value: t.Literal}

	case lexer.LBRACE:
		expr = p.parseArrayLit()

	case lexer.LPAREN:
		p.parenDepth++
		p.next()
		expr = p.parseExpr()
		p.eat(lexer.RPAREN)

	case lexer.IDENT:
		p.next()
		name := t.Literal
		// Boolean predefined constants
		switch name {
		case "TRUE", "SI":
			expr = &ast.BoolLit{Line: t.Line, Value: true}
		case "FALSE", "NO":
			expr = &ast.BoolLit{Line: t.Line, Value: false}
		default:
			// Function call or variable
			if p.peek().Type == lexer.LPAREN {
				p.parenDepth++
				p.next()
				args := p.parseArgList()
				p.eat(lexer.RPAREN)
				expr = &ast.CallExpr{Line: t.Line, Name: name, Args: args}
			} else {
				expr = &ast.IdentExpr{Line: t.Line, Name: name}
			}
		}

	default:
		return nil
	}

	// Apply selectors: [index], [i,j], .campo
	expr = p.parseSelectors(expr)
	return expr
}

func (p *Parser) parseSelectors(base ast.Expr) ast.Expr {
	for {
		switch p.peek().Type {
		case lexer.LBRACK:
			line := p.peek().Line
			p.next()
			indices := []ast.Expr{p.parseExpr()}
			for p.peek().Type == lexer.COMMA {
				p.next()
				indices = append(indices, p.parseExpr())
			}
			p.eat(lexer.RBRACK)
			base = &ast.IndexExpr{Line: line, Array: base, Indices: indices}
		case lexer.DOT:
			line := p.peek().Line
			p.next()
			field := p.next().Literal
			base = &ast.FieldExpr{Line: line, Record: base, Field: field}
		default:
			return base
		}
	}
}

func (p *Parser) parseArgList() []ast.Expr {
	if p.peek().Type == lexer.RPAREN {
		return nil
	}
	args := []ast.Expr{p.parseExpr()}
	for p.peek().Type == lexer.COMMA {
		p.next()
		args = append(args, p.parseExpr())
	}
	return args
}

func (p *Parser) parseArrayLit() ast.Expr {
	line := p.peek().Line
	p.next() // consume {
	var elems []ast.Expr
	fill := false
	for p.peek().Type != lexer.RBRACE && p.peek().Type != lexer.EOF {
		if p.peek().Type == lexer.ELLIPSIS {
			p.next()
			fill = true
			break
		}
		elems = append(elems, p.parseExpr())
		if p.peek().Type == lexer.COMMA {
			p.next()
		} else {
			break
		}
		// Check for trailing ... after comma
		if p.peek().Type == lexer.ELLIPSIS {
			p.next()
			fill = true
			break
		}
	}
	p.eat(lexer.RBRACE)
	return &ast.ArrayLit{Line: line, Elems: elems, Fill: fill}
}

// Parse parses an SL program from source and returns the AST.
func Parse(src string) (*ast.Program, []*ParseError) {
	p := New(src)
	prog, _ := p.ParseProgram()
	return prog, p.errors
}
