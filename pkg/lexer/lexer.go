// Package lexer implements the SL language tokenizer.
package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer holds the tokenizer state.
type Lexer struct {
	src    []rune
	pos    int
	line   int
	col    int
	peeked *Token
}

// New creates a Lexer for src.
func New(src string) *Lexer {
	return &Lexer{src: []rune(src), pos: 0, line: 1, col: 1}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.src) {
		return 0
	}
	return l.src[l.pos]
}

func (l *Lexer) peek2() rune {
	if l.pos+1 >= len(l.src) {
		return 0
	}
	return l.src[l.pos+1]
}

func (l *Lexer) advance() rune {
	if l.pos >= len(l.src) {
		return 0
	}
	r := l.src[l.pos]
	l.pos++
	if r == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return r
}

func (l *Lexer) skipWhitespaceAndComments() {
	for l.pos < len(l.src) {
		ch := l.peek()
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			l.advance()
			continue
		}
		if ch == '/' {
			if l.peek2() == '/' {
				// single-line comment
				for l.pos < len(l.src) && l.peek() != '\n' {
					l.advance()
				}
				continue
			}
			if l.peek2() == '*' {
				// multi-line comment
				l.advance() // /
				l.advance() // *
				for l.pos < len(l.src) {
					if l.peek() == '*' && l.peek2() == '/' {
						l.advance() // *
						l.advance() // /
						break
					}
					l.advance()
				}
				continue
			}
		}
		break
	}
}

// isIdentStart returns true if r can start an SL identifier.
func isIdentStart(r rune) bool {
	if r == '_' {
		return true
	}
	if r == 'ñ' || r == 'Ñ' {
		return true
	}
	// ASCII letters only (no accented vowels)
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// isIdentPart returns true if r can continue an SL identifier.
func isIdentPart(r rune) bool {
	if isIdentStart(r) {
		return true
	}
	return r >= '0' && r <= '9'
}

// isDigit returns true for ASCII digits.
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (l *Lexer) makeToken(t TokenType, lit string, line, col int) Token {
	return Token{Type: t, Literal: lit, Line: line, Col: col}
}

// Next returns the next token, consuming it.
func (l *Lexer) Next() Token {
	if l.peeked != nil {
		t := *l.peeked
		l.peeked = nil
		return t
	}
	return l.nextToken()
}

// Peek returns the next token without consuming it.
func (l *Lexer) Peek() Token {
	if l.peeked == nil {
		t := l.nextToken()
		l.peeked = &t
	}
	return *l.peeked
}

func (l *Lexer) nextToken() Token {
	l.skipWhitespaceAndComments()

	if l.pos >= len(l.src) {
		return l.makeToken(EOF, "", l.line, l.col)
	}

	line, col := l.line, l.col
	ch := l.peek()

	// Ellipsis
	if ch == '.' && l.peek2() == '.' {
		l.advance()
		l.advance()
		if l.peek() == '.' {
			l.advance()
			return l.makeToken(ELLIPSIS, "...", line, col)
		}
		return l.makeToken(ILLEGAL, "..", line, col)
	}

	// Number
	if isDigit(ch) || (ch == '.' && isDigit(l.peek2())) {
		return l.readNumber(line, col)
	}

	// String — ASCII and Unicode smart quotes (all variants)
	if isDoubleQuote(ch) || isSingleQuote(ch) {
		return l.readString(line, col)
	}

	// Identifier or keyword
	if isIdentStart(ch) {
		return l.readIdent(line, col)
	}

	// Operators and delimiters
	l.advance()
	switch ch {
	case '+':
		return l.makeToken(PLUS, "+", line, col)
	case '-':
		return l.makeToken(MINUS, "-", line, col)
	case '*':
		return l.makeToken(STAR, "*", line, col)
	case '/':
		return l.makeToken(SLASH, "/", line, col)
	case '%':
		return l.makeToken(PERCENT, "%", line, col)
	case '^':
		return l.makeToken(CARET, "^", line, col)
	case '=':
		if l.peek() == '=' {
			l.advance()
			return l.makeToken(EQ, "==", line, col)
		}
		return l.makeToken(ASSIGN, "=", line, col)
	case '<':
		if l.peek() == '>' {
			l.advance()
			return l.makeToken(NEQ, "<>", line, col)
		}
		if l.peek() == '=' {
			l.advance()
			return l.makeToken(LE, "<=", line, col)
		}
		return l.makeToken(LT, "<", line, col)
	case '>':
		if l.peek() == '=' {
			l.advance()
			return l.makeToken(GE, ">=", line, col)
		}
		return l.makeToken(GT, ">", line, col)
	case '&':
		if l.peek() == '&' {
			l.advance()
			return l.makeToken(AND2, "&&", line, col)
		}
		return l.makeToken(ILLEGAL, "&", line, col)
	case '|':
		if l.peek() == '|' {
			l.advance()
			return l.makeToken(OR2, "||", line, col)
		}
		return l.makeToken(ILLEGAL, "|", line, col)
	case '(':
		return l.makeToken(LPAREN, "(", line, col)
	case ')':
		return l.makeToken(RPAREN, ")", line, col)
	case '[':
		return l.makeToken(LBRACK, "[", line, col)
	case ']':
		return l.makeToken(RBRACK, "]", line, col)
	case '{':
		return l.makeToken(LBRACE, "{", line, col)
	case '}':
		return l.makeToken(RBRACE, "}", line, col)
	case ',':
		return l.makeToken(COMMA, ",", line, col)
	case ';':
		return l.makeToken(SEMI, ";", line, col)
	case '.':
		return l.makeToken(DOT, ".", line, col)
	case ':':
		return l.makeToken(COLON, ":", line, col)
	}

	return l.makeToken(ILLEGAL, string(ch), line, col)
}

func (l *Lexer) readNumber(line, col int) Token {
	var sb strings.Builder
	for l.pos < len(l.src) && isDigit(l.peek()) {
		sb.WriteRune(l.advance())
	}
	if l.peek() == '.' && isDigit(l.peek2()) {
		sb.WriteRune(l.advance()) // .
		for l.pos < len(l.src) && isDigit(l.peek()) {
			sb.WriteRune(l.advance())
		}
	}
	// Optional exponent (e.g. 1.5e10)
	if l.peek() == 'e' || l.peek() == 'E' {
		sb.WriteRune(l.advance())
		if l.peek() == '+' || l.peek() == '-' {
			sb.WriteRune(l.advance())
		}
		for l.pos < len(l.src) && isDigit(l.peek()) {
			sb.WriteRune(l.advance())
		}
	}
	return l.makeToken(NUMBER, sb.String(), line, col)
}

// matchingCloseQuote returns the closing quote for a given opening quote.
func matchingCloseQuote(open rune) rune {
	switch open {
	case '\u201C': // "  LEFT DOUBLE QUOTATION MARK
		return '\u201D' // "  RIGHT DOUBLE QUOTATION MARK
	case '\u2018': // '  LEFT SINGLE QUOTATION MARK
		return '\u2019' // '  RIGHT SINGLE QUOTATION MARK
	default:
		return open // ASCII " or ' close with the same char
	}
}

func isDoubleQuote(r rune) bool {
	return r == '"' || r == '\u201C' || r == '\u201D'
}

func isSingleQuote(r rune) bool {
	return r == '\'' || r == '\u2018' || r == '\u2019'
}

func (l *Lexer) readString(line, col int) Token {
	open := l.advance() // consume opening quote
	closeQ := matchingCloseQuote(open)
	isDoubleOpen := isDoubleQuote(open)
	isSingleOpen := isSingleQuote(open)
	var sb strings.Builder
	for l.pos < len(l.src) {
		ch := l.peek()
		// Accept any same-family quote as closing (handles copy-paste smart-quote inconsistencies)
		isClose := ch == closeQ ||
			(isDoubleOpen && isDoubleQuote(ch)) ||
			(isSingleOpen && isSingleQuote(ch))
		if isClose {
			l.advance()
			break
		}
		if ch == '\\' {
			l.advance() // consume backslash
			esc := l.advance()
			switch esc {
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case '\\':
				sb.WriteByte('\\')
			case '\'', '\u2019':
				sb.WriteByte('\'')
			case '"', '\u201D':
				sb.WriteByte('"')
			default:
				sb.WriteRune(esc)
			}
			continue
		}
		sb.WriteRune(l.advance())
	}
	return l.makeToken(STRING, sb.String(), line, col)
}

func (l *Lexer) readIdent(line, col int) Token {
	var sb strings.Builder
	for l.pos < len(l.src) && isIdentPart(l.peek()) {
		sb.WriteRune(l.advance())
	}
	lit := sb.String()
	tt := LookupIdent(lit)
	return l.makeToken(tt, lit, line, col)
}

// TokenizeAll tokenizes the entire source and returns all tokens (for testing).
func TokenizeAll(src string) ([]Token, error) {
	l := New(src)
	var tokens []Token
	for {
		t := l.Next()
		tokens = append(tokens, t)
		if t.Type == EOF {
			break
		}
		if t.Type == ILLEGAL {
			return tokens, fmt.Errorf("illegal token %q at line %d col %d", t.Literal, t.Line, t.Col)
		}
	}
	return tokens, nil
}

// Ensure imports are used.
var _ = unicode.IsLetter
var _ = utf8.RuneError
var _ = fmt.Sprintf
