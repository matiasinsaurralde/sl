package token

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

type Lexer struct {
	reader *bufio.Reader
	pos    int
	line   int
	col    int
}

type TokenInfo struct {
	Type    Token
	Literal string
	Pos     Pos
	Line    int
	Col     int
}

func NewLexer(input io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(input),
		pos:    0,
		line:   1,
		col:    0,
	}
}

func (l *Lexer) NextToken() TokenInfo {
	l.skipWhitespace()

	pos := Pos(l.pos)
	line := l.line
	col := l.col

	ch := l.readChar()

	switch ch {
	case 0:
		return TokenInfo{Type: EOF, Literal: "", Pos: pos, Line: line, Col: col}
	case '(':
		return TokenInfo{Type: LPAREN, Literal: "(", Pos: pos, Line: line, Col: col}
	case ')':
		return TokenInfo{Type: RPAREN, Literal: ")", Pos: pos, Line: line, Col: col}
	case '{':
		return TokenInfo{Type: LBRACE, Literal: "{", Pos: pos, Line: line, Col: col}
	case '}':
		return TokenInfo{Type: RBRACE, Literal: "}", Pos: pos, Line: line, Col: col}
	case '[':
		return TokenInfo{Type: LBRACKET, Literal: "[", Pos: pos, Line: line, Col: col}
	case ']':
		return TokenInfo{Type: RBRACKET, Literal: "]", Pos: pos, Line: line, Col: col}
	case ',':
		return TokenInfo{Type: COMMA, Literal: ",", Pos: pos, Line: line, Col: col}
	case ';':
		return TokenInfo{Type: SEMICOLON, Literal: ";", Pos: pos, Line: line, Col: col}
	case ':':
		return TokenInfo{Type: COLON, Literal: ":", Pos: pos, Line: line, Col: col}
	case '+':
		return TokenInfo{Type: PLUS, Literal: "+", Pos: pos, Line: line, Col: col}
	case '-':
		return TokenInfo{Type: MINUS, Literal: "-", Pos: pos, Line: line, Col: col}
	case '*':
		return TokenInfo{Type: MULTIPLY, Literal: "*", Pos: pos, Line: line, Col: col}
	case '/':
		if l.peekChar() == '*' {
			return l.readComment()
		}
		return TokenInfo{Type: DIVIDE, Literal: "/", Pos: pos, Line: line, Col: col}
	case '%':
		return TokenInfo{Type: MODULO, Literal: "%", Pos: pos, Line: line, Col: col}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			return TokenInfo{Type: LTE, Literal: "<=", Pos: pos, Line: line, Col: col}
		} else if l.peekChar() == '>' {
			l.readChar()
			return TokenInfo{Type: NEQ, Literal: "<>", Pos: pos, Line: line, Col: col}
		}
		return TokenInfo{Type: LT, Literal: "<", Pos: pos, Line: line, Col: col}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			return TokenInfo{Type: GTE, Literal: ">=", Pos: pos, Line: line, Col: col}
		}
		return TokenInfo{Type: GT, Literal: ">", Pos: pos, Line: line, Col: col}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			return TokenInfo{Type: EQ, Literal: "==", Pos: pos, Line: line, Col: col}
		}
		return TokenInfo{Type: ASSIGN, Literal: "=", Pos: pos, Line: line, Col: col}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			return TokenInfo{Type: AND, Literal: "&&", Pos: pos, Line: line, Col: col}
		}
		return TokenInfo{Type: ILLEGAL, Literal: string(ch), Pos: pos, Line: line, Col: col}
	case '"':
		return l.readString()
	case '\n':
		return TokenInfo{Type: EOL, Literal: "\n", Pos: pos, Line: line, Col: col}
	default:
		if unicode.IsLetter(rune(ch)) {
			return l.readIdentifier(ch)
		} else if unicode.IsDigit(rune(ch)) {
			return l.readNumber(ch)
		}
		return TokenInfo{Type: ILLEGAL, Literal: string(ch), Pos: pos, Line: line, Col: col}
	}
}

func (l *Lexer) readChar() byte {
	ch, err := l.reader.ReadByte()
	if err != nil {
		return 0
	}
	l.pos++
	l.col++
	if ch == '\n' {
		l.line++
		l.col = 0
	}
	return ch
}

func (l *Lexer) peekChar() byte {
	ch, err := l.reader.ReadByte()
	if err != nil {
		return 0
	}
	l.reader.UnreadByte()
	return ch
}

func (l *Lexer) skipWhitespace() {
	for {
		ch := l.peekChar()
		if ch == 0 {
			break
		}
		if !unicode.IsSpace(rune(ch)) {
			break
		}
		l.readChar()
	}
}

func (l *Lexer) readComment() TokenInfo {
	startPos := Pos(l.pos - 1)
	startLine := l.line
	startCol := l.col - 1

	var buf bytes.Buffer
	buf.WriteString("/*")

	// Read until we find */
	for {
		ch := l.readChar()
		if ch == 0 {
			break
		}
		buf.WriteByte(ch)
		if ch == '*' && l.peekChar() == '/' {
			buf.WriteByte(l.readChar())
			break
		}
	}

	return TokenInfo{
		Type:    COMMENT,
		Literal: buf.String(),
		Pos:     startPos,
		Line:    startLine,
		Col:     startCol,
	}
}

func (l *Lexer) readString() TokenInfo {
	startPos := Pos(l.pos - 1)
	startLine := l.line
	startCol := l.col - 1

	var buf bytes.Buffer
	buf.WriteByte('"')

	for {
		ch := l.readChar()
		if ch == 0 {
			break
		}
		buf.WriteByte(ch)
		if ch == '"' {
			break
		}
	}

	return TokenInfo{
		Type:    STRING,
		Literal: buf.String(),
		Pos:     startPos,
		Line:    startLine,
		Col:     startCol,
	}
}

func (l *Lexer) readIdentifier(first byte) TokenInfo {
	startPos := Pos(l.pos - 1)
	startLine := l.line
	startCol := l.col - 1

	var buf bytes.Buffer
	buf.WriteByte(first)

	for {
		ch := l.peekChar()
		if ch == 0 {
			break
		}
		if !unicode.IsLetter(rune(ch)) && !unicode.IsDigit(rune(ch)) && ch != '_' {
			break
		}
		buf.WriteByte(l.readChar())
	}

	literal := buf.String()
	tokenType := Lookup(literal)

	return TokenInfo{
		Type:    tokenType,
		Literal: literal,
		Pos:     startPos,
		Line:    startLine,
		Col:     startCol,
	}
}

func (l *Lexer) readNumber(first byte) TokenInfo {
	startPos := Pos(l.pos - 1)
	startLine := l.line
	startCol := l.col - 1

	var buf bytes.Buffer
	buf.WriteByte(first)

	hasDecimal := false

	for {
		ch := l.peekChar()
		if ch == 0 {
			break
		}
		if ch == '.' && !hasDecimal {
			hasDecimal = true
			buf.WriteByte(l.readChar())
		} else if unicode.IsDigit(rune(ch)) {
			buf.WriteByte(l.readChar())
		} else {
			break
		}
	}

	literal := buf.String()
	tokenType := INT
	if hasDecimal {
		tokenType = FLOAT
	}

	return TokenInfo{
		Type:    Token(tokenType),
		Literal: literal,
		Pos:     startPos,
		Line:    startLine,
		Col:     startCol,
	}
}
