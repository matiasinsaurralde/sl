package parser

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	logger "github.com/matiasinsaurralde/sl/log"
)

var (
	log = logger.Logger
)

// Token es utilizado para enumerar los tokens.
type Token int

var (
	exprMatch = regexp.MustCompile(`(\+|\*|\/|-)`)
)

const (
	// UNKNOWN es utilizado por defecto:
	UNKNOWN Token = iota

	// IDENT representa un identificador.
	IDENT

	// EXPR representa una expresión.
	EXPR

	// LITINT representa un literal entero.
	LITINT

	// LITSTR representa un literal string/cadena.
	LITSTR

	// KWPROGRAMA indica el inicio de un programa.
	KWPROGRAMA

	// KWINICIO indica el inicio de una rutina.
	KWINICIO

	// KWFIN indica el fin de una rutina.
	KWFIN

	// KWVAR indica una declaración de variable.
	KWVAR

	// KWSUBRUTINA indica una subrutina.
	KWSUBRUTINA
)

// TokenInstance almacena información de cada token (posición y tipo).
type TokenInstance struct {
	Type    Token
	Literal string

	pos int64
	ln  int64
}

// Parser es la estructura de datos principal del parser.
type Parser struct {
	r io.Reader

	tokens []*TokenInstance
	pos    int64
	ln     int64
	chbuf  []byte
	buf    bytes.Buffer
}

// New inicializa un nuevo parser.
func New(r io.Reader) (*Parser, error) {
	return &Parser{
		r:      r,
		tokens: make([]*TokenInstance, 0),
		ln:     1,
	}, nil
}

func (p *Parser) isExpr() {

}

func (p *Parser) lookupToken() {
	if p.buf.Len() == 0 {
		return
	}
	s := strings.TrimSpace(p.buf.String())
	var t Token
	switch s {
	case "programa":
		t = KWPROGRAMA
	case "inicio":
		t = KWINICIO
	case "fin":
		t = KWFIN
	case "var":
		t = KWVAR
	case "subrutina":
		t = KWSUBRUTINA
	default:
		// ¿Es un entero?
		_, err := strconv.Atoi(s)
		if err == nil {
			t = LITINT
			break
		}

		// ¿Es una expresión?
		if exprMatch.MatchString(s) {
			t = EXPR
			break
		}

		// TODO: mejorar el reconocimiento de cadenas
		firstCh := s[0:1]
		lastCh := s[len(s)-1:]
		if firstCh == `"` && lastCh == `"` {
			t = LITSTR
			// fmt.Println("LIT_STR =", s)
			break
		}

		// Por defecto tomamos los tokens como identificadores:
		t = IDENT
		// fmt.Println("IDENT = ", s)
	}
	tok := &TokenInstance{
		Type:    t,
		Literal: s,
		pos:     p.pos,
		ln:      p.ln,
	}
	p.tokens = append(p.tokens, tok)
	// fmt.Println("FOUND = ", s, t)
}

// Parse tokeniza y prepara el AST del programa.
func (p *Parser) Parse() {
	p.chbuf = make([]byte, 1)
	p.buf = bytes.Buffer{}
	var ch string
	var lit bool
	for p.Scan() {
		ch = string(p.chbuf)
		if lit {
			// fmt.Println("Writing lit!", p.buf.String())
			p.buf.Write(p.chbuf)
			if ch == `"` {
				lit = false
			}
			continue
		}
		if ch == "\n" {
			p.lookupToken()
			p.buf.Reset()
			p.ln++
		} else if ch == " " {
			p.lookupToken()
			p.buf.Reset()
		} else if ch == `"` {
			// fmt.Println("STRING LIT starts")
			p.lookupToken()
			p.buf.Reset()
			p.buf.Write(p.chbuf)
			lit = true
		} else if ch == "(" {
			p.lookupToken()
			// fmt.Println("LBRACKET")
			p.buf.Reset()
		} else if ch == ")" {
			p.lookupToken()
			// fmt.Println("RBRACKET")
			p.buf.Reset()
		} else {
			p.buf.Write(p.chbuf)
		}
		// fmt.Println("buf is ", p.buf.String())
	}

	log.Debugf("%d tokens encontrados", len(p.tokens))
	fmt.Println("Found", len(p.tokens), "tokens")
	i := 0
	for {
		if i == len(p.tokens)-1 {
			break
		}
		tok := p.tokens[i]
		nextTok := p.tokens[i+1]
		if nextTok == nil {
			break
		}
		switch tok.Type {
		// Requerir un identificador válido para el programa:
		case KWPROGRAMA:
			if nextTok.Type != IDENT {
				p.Error(tok)
			}
		}
		fmt.Println(tok, " ->", nextTok)
		i++
	}
}

// Error es un helper para indicar errores detallados.
func (p *Parser) Error(tok *TokenInstance) {
	log.Errorf("Error en linea %d posición %d, identificador inválido", tok.ln, tok.pos)
	os.Exit(1)
}

// Scan lee el io.Reader byte por byte.
func (p *Parser) Scan() bool {
	_, err := p.r.Read(p.chbuf)
	if err == io.EOF {
		return false
	}
	p.pos++
	return true
}
