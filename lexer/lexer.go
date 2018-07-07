package parser

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"time"

	logger "github.com/matiasinsaurralde/sl/log"
	token "github.com/matiasinsaurralde/sl/token"
)

var (
	log = logger.Logger
)

var (
	numMatch = regexp.MustCompile(`^[0-9]*$`)
)

// Lexer es la estructura de datos principal del lexer.
type Lexer struct {
	r io.Reader

	scanCh  chan string
	tokenCh chan token.Item

	tokens []*token.Item
	done   chan bool

	tok token.Token
	pos int64
	ln  int64
	buf bytes.Buffer
}

// New inicializa un nuevo lexer.
func New(r io.Reader) (*Lexer, error) {
	return &Lexer{
		r:       r,
		tokens:  make([]*token.Item, 0),
		ln:      1,
		scanCh:  make(chan string, 1),
		tokenCh: make(chan token.Item, 1),
		done:    make(chan bool),
	}, nil
}

func (l *Lexer) isExpr() {

}

func (l *Lexer) emitToken(t *token.Item) {
	l.tokens = append(l.tokens, t)
}

func (l *Lexer) newLine() {
	l.ln++
}

func (l *Lexer) logToken(tokenName string, tok *token.Item) {
	log.WithField("prefix", "lexer").
		Debugf("Token '%s' encontrado en linea %d, posicion %d",
			tokenName, tok.Ln, tok.Pos)
}

func (l *Lexer) handleTokens() {
	for {
		tok := <-l.tokenCh
		if tok.Literal == "" {
			continue
		}
		switch tok.Literal {
		case "programa":
			tok.Type = token.KWPROGRAMA
			l.logToken("PROGRAMA", &tok)
			l.emitToken(&tok)
			continue
		case "inicio":
			tok.Type = token.KWINICIO
			l.logToken("INICIO", &tok)
			l.emitToken(&tok)
			continue
		case "fin":
			l.logToken("FIN", &tok)
			l.emitToken(&tok)
			continue
		case "var":
			tok.Type = token.KWVAR
			l.logToken("VAR", &tok)
			l.emitToken(&tok)
			continue
		case "subrutina":
			l.logToken("SUBRUTINA", &tok)
			l.emitToken(&tok)
			continue
		case "desde":
			l.logToken("DESDE", &tok)
			l.emitToken(&tok)
			continue
		case "hasta":
			l.logToken("HASTA", &tok)
			l.emitToken(&tok)
			continue
		case "paso":
			l.logToken("PASO", &tok)
			l.emitToken(&tok)
			continue
		default:
			switch tok.Type {
			case token.PARENEXPR:
				l.logToken("PARENEXPR", &tok)
				l.emitToken(&tok)
			default:
				l.emitToken(&tok)
			}
		}
	}
}

// Parse tokeniza y prepara el AST del programa.
func (l *Lexer) Parse() []*token.Item {
	log.WithField("prefix", "lexer").
		Debug("Inicio")
	l.buf = bytes.Buffer{}
	go l.handleTokens()
	go l.scan()
	var strLit bool
	var intLit, floatLit bool
	var parenExpr, blockStmt bool
	var singleLnComment bool
	for {
		select {
		case ch := <-l.scanCh:
			// Handle single line comments:
			if l.buf.String() == "//" && !singleLnComment {
				l.buf.WriteString(ch)
				singleLnComment = true
				continue
			}
			if singleLnComment && ch != "\n" {
				l.buf.WriteString(ch)
				continue
			}

			// Eat all chars inside parent expression:
			if parenExpr && ch != ")" {
				l.buf.WriteString(ch)
				continue
			}

			// Eat all chars inside block statements:
			if blockStmt && ch != "}" {
				l.buf.WriteString(ch)
				continue
			}

			// Handle string literals:
			if ch == `"` {
				if !strLit {
					strLit = true
					l.buf.Reset()
					continue
				}
				// l.tokenCh <- "STRLIT:" + l.buf.String()
				l.tokenCh <- token.Item{
					Type:    token.LITSTR,
					Literal: l.buf.String(),
					Pos:     l.pos,
					Ln:      l.ln,
				}
				l.buf.Reset()
				strLit = false
				continue
			}
			if strLit {
				l.buf.WriteString(ch)
				continue
			}

			// Handle digits:
			isNum := l.isNumber(ch)
			if isNum {
				intLit = true
				l.buf.WriteString(ch)
				continue
			}

			if !isNum && intLit {
				if ch == "." && !floatLit {
					l.buf.WriteString(ch)
					floatLit = true
					continue
				}
				if floatLit {
					floatLit = false
					intLit = false
					// l.tokenCh <- "FLOATLIT:" + l.buf.String()
					l.tokenCh <- token.Item{
						Type:    token.LITFLOAT,
						Literal: l.buf.String(),
						Pos:     l.pos,
						Ln:      l.ln,
					}
					l.buf.Reset()
					// continue
				} else {
					// l.tokenCh <- "INTLIT:" + l.buf.String()
					l.tokenCh <- token.Item{
						Type:    token.LITINT,
						Literal: l.buf.String(),
						Pos:     l.pos,
						Ln:      l.ln,
					}
					l.buf.Reset()
					intLit = false
				}
			}
			switch ch {
			case "\n":
				l.newLine()
				if singleLnComment {
					singleLnComment = false
					l.buf.Reset()
					continue
				}
				// l.tokenCh <- l.buf.String()
				l.tokenCh <- token.Item{
					Type:    token.COMMENT,
					Literal: l.buf.String(),
					Pos:     l.pos,
					Ln:      l.ln,
				}
			case " ":
				// l.tokenCh <- l.buf.String()
				l.tokenCh <- token.Item{
					Type:    token.UNKNOWN,
					Literal: l.buf.String(),
					Pos:     l.pos,
					Ln:      l.ln,
				}
			case "\t":
				// l.tokenCh <- l.buf.String()
				l.tokenCh <- token.Item{
					Type:    token.UNKNOWN,
					Literal: l.buf.String(),
					Pos:     l.pos,
					Ln:      l.ln,
				}
			case ":":
				if l.buf.Len() > 0 {
					// l.tokenCh <- l.buf.String()
					l.tokenCh <- token.Item{
						Type:    token.UNKNOWN,
						Literal: l.buf.String(),
						Pos:     l.pos,
						Ln:      l.ln,
					}
					l.buf.Reset()
				}
				// l.tokenCh <- "EXPECTED_TYPE"
				l.tokenCh <- token.Item{
					Type:    token.COLON,
					Literal: ":",
					Pos:     l.pos,
					Ln:      l.ln,
				}
			case "=":
				if l.buf.Len() > 0 {
					// l.tokenCh <- l.buf.String()
					l.tokenCh <- token.Item{
						Type:    token.UNKNOWN,
						Literal: l.buf.String(),
						Pos:     l.pos,
						Ln:      l.ln,
					}
					l.buf.Reset()
					// Fix line count?
					l.ln--
				}
				// l.tokenCh <- "EXPECT_DECL"
				l.tokenCh <- token.Item{
					Type:    token.EQ,
					Literal: "=",
					Pos:     l.pos,
					Ln:      l.ln,
				}
			case "(":
				/*
					l.tokenCh <- token.Item{
						Type:    token.LPAREN,
						Literal: "(",
					}*/
				if l.buf.Len() > 0 {
					l.tokenCh <- token.Item{
						Type:    token.UNKNOWN,
						Literal: l.buf.String(),
						Pos:     l.pos,
						Ln:      l.ln,
					}
					l.buf.Reset()
				}
				parenExpr = true
			case ")":
				l.tokenCh <- token.Item{
					Type:    token.PARENEXPR,
					Literal: l.buf.String(),
					Pos:     l.pos,
					Ln:      l.ln,
				}
				if parenExpr {
					parenExpr = false
				}
				l.buf.Reset()
			case "{":
				l.tokenCh <- token.Item{
					Type:    token.LBRACKET,
					Literal: "{",
					Pos:     l.pos,
					Ln:      l.ln,
				}
				blockStmt = true
			case "}":
				l.tokenCh <- token.Item{
					Type:    token.BLOCK,
					Literal: l.buf.String(),
					Pos:     l.pos,
					Ln:      l.ln,
				}
				if blockStmt {
					blockStmt = false
				}
				l.tokenCh <- token.Item{
					Type:    token.RBRACKET,
					Literal: "}",
					Pos:     l.pos,
					Ln:      l.ln,
				}
				l.buf.Reset()
			case "+":
				l.tokenCh <- token.Item{
					Type:    token.OP,
					Literal: "+",
					Pos:     l.pos,
					Ln:      l.ln,
				}
			case "-":
				l.tokenCh <- token.Item{
					Type:    token.OP,
					Literal: "-",
					Pos:     l.pos,
					Ln:      l.ln,
				}
			default:
				l.buf.WriteString(ch)
				continue
			}
			l.buf.Reset()
		case <-l.done:
			log.WithField("prefix", "lexer").
				Debug("Listo.")
			return l.tokens
		}
	}
}

func (l *Lexer) currentToken() token.Token {
	return l.tok
}

func (l *Lexer) clear() {
	l.tok = token.NIL
}

// Error es un helper para indicar errores detallados.
func (l *Lexer) Error(tok *token.Item) {
	log.Errorf("Error en linea %d posición %d, identificador inválido", tok.Ln, tok.Pos)
	os.Exit(1)
}

// Scan lee el io.Reader byte por byte.
func (l *Lexer) scan() {
	buf := make([]byte, 1)
	for {
		_, err := l.r.Read(buf)
		if err == io.EOF {
			break
		}
		s := string(buf)
		l.scanCh <- s
		l.pos++
	}
	l.scanCh <- "\n"
	time.Sleep(10 * time.Millisecond)
	l.done <- true
}

func (l *Lexer) isNumber(input string) bool {
	return numMatch.MatchString(input)
}
