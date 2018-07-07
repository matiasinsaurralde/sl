package parser

import (
	"fmt"
	"io"
	"os"
	"strconv"

	slast "github.com/matiasinsaurralde/sl/ast"
	logger "github.com/matiasinsaurralde/sl/log"
	token "github.com/matiasinsaurralde/sl/token"
)

var (
	log = logger.Logger
)

// Parser contiene la estructura del parser.
type Parser struct {
	tokenSet []*token.Item
	ast      *slast.AST

	i int64
}

// BuildValue inicializa un valor en Go, basado en un token.
func BuildValue(name string, typeTok *token.Item, valueTok *token.Item) (interface{}, error) {
	if valueTok == nil {
		log.WithField("prefix", "parser").
			Debugf("buildValue: '%s', type = %s, value = nil", name, typeTok.Literal)
	} else {
		log.WithField("prefix", "parser").
			Debugf("buildValue: '%s', type = %s, value = %s", name, typeTok.Literal, valueTok.Literal)
	}

	switch typeTok.Literal {
	case "numerico":
		// Initialize with zero values (int64 by default):
		if valueTok == nil {
			var n int64
			return n, nil
		}
		if valueTok.Type == token.LITFLOAT {
			n, err := strconv.ParseFloat(valueTok.Literal, 64)
			if err != nil {
				panic("Couldn't parse float")
			}
			return n, err
		}
		n, err := strconv.Atoi(valueTok.Literal)
		if err != nil {
			panic("Type error")
		}
		return n, err
	case "cadena":
		var s string
		if valueTok == nil {
			s = ""
			return s, nil
		}
		s = valueTok.Literal
		return s, nil
	default:
		return nil, fmt.Errorf("Tipo desconocido: '%s'", typeTok.Literal)
	}
}

// New inicializa un nuevo parser, toma como parámetro un set de tokens.
func New(r io.Reader, tokenSet []*token.Item) *Parser {
	log.WithField("prefix", "parser").Debugf("Inicio: %d tokens encontrados", len(tokenSet))
	p := &Parser{
		tokenSet: tokenSet,
	}
	return p
}

// Next retorna el token siguiente, nil si no está definido.
func (p *Parser) Next() *token.Item {
	p.i++
	if p.i > int64(len(p.tokenSet)-1) {
		return nil
	}
	return p.tokenSet[p.i]
}

// Rewind retorna el token anterior.
func (p *Parser) Rewind() *token.Item {
	p.i--
	return p.tokenSet[p.i]
}

// RewindTo retorna el token de una posición dada.
func (p *Parser) RewindTo(n int64) {
	p.i = n
}

func (p *Parser) currentToken() *token.Item {
	return p.tokenSet[p.i]
}

func (p *Parser) parseGlobalDecl(scope *slast.Scope, name string, t *token.Item, v *token.Item) error {
	if scope == nil {
		scope = p.ast.GlobalScope
	}
	val, err := BuildValue(name, t, v)
	if err != nil {
		return err
	}
	var node slast.Node = &slast.Decl{
		Name:  name,
		Value: val,
	}
	scope.Objects[name] = node
	return nil
}

func solveExpr(e *slast.BinaryExpr) int {
	fmt.Println("\tsolveExpr", "X=", e.X)
	fmt.Println("\tsolveExpr", "Y=", e.Y)
	fmt.Println()
	// for {
	var X, Y slast.Node
	X = e.X
	Y = e.Y

	var xVal int
	var yVal int

	var n int

	// for {
	// for {
	switch X.(type) {
	case *slast.BasicLit:
		l := X.(*slast.BasicLit)
		xVal, _ = strconv.Atoi(l.Literal)
	case *slast.BinaryExpr:
		subExpr := X.(*slast.BinaryExpr)
		xVal = solveExpr(subExpr)
	}

	switch Y.(type) {
	case *slast.BasicLit:
		l := Y.(*slast.BasicLit)
		yVal, _ = strconv.Atoi(l.Literal)
	case *slast.BinaryExpr:
		subExpr := Y.(*slast.BinaryExpr)
		yVal = solveExpr(subExpr)
	}
	op := e.Op
	switch op {
	case "+":
		fmt.Println(xVal, "+", yVal)
		n = xVal + yVal
	case "-":
		fmt.Println(xVal, "-", yVal)
		n = xVal - yVal
	case "*":
		fmt.Println(xVal, "*", yVal)
		n = xVal * yVal
	case "/":
		fmt.Println(xVal, "/", yVal)
		n = xVal / yVal
	}
	// }
	// }
	// }
	fmt.Println("result = ", n)
	return n
}

func (p *Parser) parseBinaryExpr(op string) {
	fmt.Println("parseExpr is called")
	p.Rewind()
	expr := &slast.BinaryExpr{}
	for {
		if p.currentToken().Type == token.EOF {
			break
		}
		fmt.Println("cur token = ", p.currentToken().Literal)
		if p.currentToken().Type == token.OP {
			expr.Op = p.currentToken().Literal
			p.Next()
			continue
		}
		if expr.X == nil {
			// fmt.Println("X is ", p.currentToken().Literal)
			lit := &slast.BasicLit{Kind: slast.INT, Literal: p.currentToken().Literal}
			expr.X = lit
			p.Next()
			continue
		}
		if expr.Y == nil {
			// fmt.Println("Y is ", p.currentToken().Literal)
			lit := &slast.BasicLit{Kind: slast.INT, Literal: p.currentToken().Literal}
			expr.Y = lit
			p.Next()
			// continue
		}
		if p.currentToken().Type == token.OP {
			newExpr := &slast.BinaryExpr{
				X:  expr,
				Op: p.currentToken().Literal,
			}
			expr = newExpr
			p.Next()
		}
		// p.Next()
	}
	// n := solveExpr(expr)
	// fmt.Println("EXPR = X = ", expr.X, "y = ", expr.Y, "op =", expr.Op)
	fmt.Println("---")
	// fmt.Println(" =", n)
	// var aaa = expr.X
	// fmt.Println(aaa)

	os.Exit(0)
}

// Parse agrupa toda la lógica requerida por el parser.
func (p *Parser) Parse() *slast.AST {
	p.ast = slast.New()
	var nextTok *token.Item
	scope := p.ast.GlobalScope
	// expr := []token.Item{}
	// var isExpr bool
	for {
		// Break if total tokens were walked:
		if p.i == int64(len(p.tokenSet)) {
			break
		}
		tok := p.tokenSet[p.i]
		switch tok.Type {
		case token.EOF:
			// break
		case token.KWINICIO:
			s := &slast.Scope{
				Nodes:      make([]slast.Node, 0),
				OuterScope: p.ast.GlobalScope,
				Objects:    make(map[string]slast.Node, 0),
			}
			mainBlock := &slast.BlockStmt{Scope: s}
			p.ast.GlobalScope.Nodes = append(p.ast.GlobalScope.Nodes, mainBlock)
			scope = s
		// TODO: Requerir un identificador válido para el programa:
		case token.KWPROGRAMA:
		case token.KWVAR:
			nextTok = p.Next()
			if nextTok == nil {
				// Handle this
				continue
			}
			varName := nextTok.Literal
			nextTok = p.Next()
			var typeTok *token.Item
			var valTok *token.Item

			if nextTok.Type == token.COLON {
				// Handle single-var declarations:
				typeTok = p.Next()
				if nextTok.Literal != ":" {
					panic("Expected type with : prefix")
				}
				// TODO: check if typeTok.Literal is a valid type name.
				valuePrefix := p.Next().Literal
				if valuePrefix != "=" {
					// Value of this variable is not set?
					err := p.parseGlobalDecl(scope, varName, typeTok, nil)
					if err != nil {
						panic(err)
					}
					continue
				}
				valTok = p.Next()
				err := p.parseGlobalDecl(scope, varName, typeTok, valTok)
				if err != nil {
					panic(err)
				}

			} else if nextTok.Type == token.COMMA {
				// Handle multiple-var declarations:
				vars := []string{varName}

				for {
					nextTok = p.Next()
					if nextTok.Type == token.IDENT {
						vars = append(vars, nextTok.Literal)
					}
					if nextTok.Type == token.COLON {
						typeTok = p.Next()
						// TODO: handle errors
						if typeTok == nil {
							break
						}
						valTok = p.Next()
						if valTok.Type != token.EQ {
							p.Rewind()
							valTok = nil
							break
						}
						valTok = p.Next()
						if valTok != nil {
							break
						}
					}
					if nextTok == nil {
						break
					}
				}
				for _, v := range vars {
					err := p.parseGlobalDecl(scope, v, typeTok, valTok)
					if err != nil {
						panic(err)
					}
				}
			}
		case token.IDENT:
			// Guess if it's a function or what?
			if tok.Literal == "imprimir" {
				var node slast.Node
				next := p.Next()
				node = &slast.Stmt{
					Name:  tok.Literal,
					Value: next.Literal,
				}
				scope.Nodes = append(scope.Nodes, node)
			}
		case token.KWFIN:
			p.ast.GlobalScope.Nodes = append(p.ast.GlobalScope.Nodes, scope)
		case token.EXPR:
			var solved bool
			for !solved {
				nextTok = p.Next()
				if nextTok == nil || nextTok.Type == token.EOF {
					break
				}
				if nextTok.IsOp() {
					fmt.Println("OP FOUND =", nextTok.Literal)
				}
			}
		case token.OP:
			fmt.Println("Found OP:", tok)
			// isExpr = true
			p.parseBinaryExpr(tok.Literal)
			break
		case token.UNKNOWN:
			nextTok = p.Next()
			if nextTok == nil {
				continue
			}
			if nextTok.Type == token.EQ {
				ident := &slast.Ident{
					Name: tok.Literal,
				}
				assign := &slast.AssignStmt{
					X: ident,
				}
				nextTok = p.Next()
				if nextTok == nil {
					continue
				}
				lit := &slast.BasicLit{Kind: slast.INT, Literal: nextTok.Literal}
				assign.Y = lit
				scope.Nodes = append(scope.Nodes, assign)
			}
			// If we have a paren expression, then this is a function call:
			if nextTok.Type == token.PARENEXPR {
				node := &slast.Stmt{
					Name:  tok.Literal,
					Value: nextTok.Literal,
				}
				log.WithField("prefix", "parser").Debugf("Agregando llamada a '%s'.", node.Name)
				scope.Nodes = append(scope.Nodes, node)
			}
		}
		p.i++
	}
	return p.ast
}
