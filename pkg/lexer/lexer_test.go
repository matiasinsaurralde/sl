package lexer

import (
	"testing"
)

func tokenize(src string) []Token {
	l := New(src)
	var tokens []Token
	for {
		tok := l.Next()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens
}

func types(tokens []Token) []TokenType {
	out := make([]TokenType, len(tokens))
	for i, t := range tokens {
		out[i] = t.Type
	}
	return out
}

func TestKeywords(t *testing.T) {
	src := "si sino mientras desde hasta sub retorna salir"
	toks := tokenize(src)
	want := []TokenType{SI, SINO, MIENTRAS, DESDE, HASTA, SUB, RETORNA, SALIR, EOF}
	got := types(toks)
	if len(got) != len(want) {
		t.Fatalf("want %d tokens, got %d: %v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("token[%d]: want %v, got %v", i, want[i], got[i])
		}
	}
}

func TestNumberLiterals(t *testing.T) {
	cases := []struct {
		src string
		val string
	}{
		{"42", "42"},
		{"3.14", "3.14"},
		{"0", "0"},
		{"1e3", "1e3"},
	}
	for _, c := range cases {
		toks := tokenize(c.src)
		if toks[0].Type != NUMBER {
			t.Errorf("%q: want NUMBER, got %v", c.src, toks[0].Type)
		}
		if toks[0].Literal != c.val {
			t.Errorf("%q: want lit %q, got %q", c.src, c.val, toks[0].Literal)
		}
	}
}

func TestStringLiterals(t *testing.T) {
	cases := []struct {
		src string
		val string
	}{
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{"\u201Csmart\u201D", "smart"},   // U+201C ... U+201D
		{"\u2018single\u2019", "single"}, // U+2018 ... U+2019
	}
	for _, c := range cases {
		toks := tokenize(c.src)
		if toks[0].Type != STRING {
			t.Errorf("%q: want STRING, got %v", c.src, toks[0].Type)
		}
		if toks[0].Literal != c.val {
			t.Errorf("%q: want lit %q, got %q", c.src, c.val, toks[0].Literal)
		}
	}
}

func TestLineComments(t *testing.T) {
	src := "a // ignore this\nb"
	toks := tokenize(src)
	tt := types(toks)
	want := []TokenType{IDENT, IDENT, EOF}
	if len(tt) != len(want) {
		t.Fatalf("want %d tokens, got %d: %v", len(want), len(tt), tt)
	}
	for i := range want {
		if tt[i] != want[i] {
			t.Errorf("token[%d]: want %v, got %v", i, want[i], tt[i])
		}
	}
}

func TestBlockComments(t *testing.T) {
	src := "a /* multi\nline */ b"
	toks := tokenize(src)
	tt := types(toks)
	want := []TokenType{IDENT, IDENT, EOF}
	if len(tt) != len(want) {
		t.Fatalf("want %d tokens, got %d: %v", len(want), len(tt), tt)
	}
}

func TestOperators(t *testing.T) {
	// In SL, '=' is ASSIGN token (used for both assignment and equality by context).
	// '==' would be EQ but SL doesn't typically use '=='; '<>' is NEQ.
	src := "+ - * / % ^ < <= > >= = <>"
	toks := tokenize(src)
	want := []TokenType{PLUS, MINUS, STAR, SLASH, PERCENT, CARET, LT, LE, GT, GE, ASSIGN, NEQ, EOF}
	got := types(toks)
	if len(got) != len(want) {
		t.Fatalf("want %d tokens, got %d: %v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("token[%d]: want %v, got %v", i, want[i], got[i])
		}
	}
}

func TestEllipsis(t *testing.T) {
	toks := tokenize("...")
	if toks[0].Type != ELLIPSIS {
		t.Errorf("want ELLIPSIS, got %v", toks[0].Type)
	}
}

func TestIdentifier(t *testing.T) {
	toks := tokenize("myVar_1")
	if toks[0].Type != IDENT {
		t.Errorf("want IDENT, got %v", toks[0].Type)
	}
	if toks[0].Literal != "myVar_1" {
		t.Errorf("want lit myVar_1, got %q", toks[0].Literal)
	}
}

func TestLookupIdent(t *testing.T) {
	if LookupIdent("si") != SI {
		t.Error("si should be SI keyword")
	}
	if LookupIdent("foobar") != IDENT {
		t.Error("foobar should be IDENT")
	}
}

func TestLineNumbers(t *testing.T) {
	src := "a\nb\nc"
	toks := tokenize(src)
	if toks[0].Line != 1 {
		t.Errorf("first token: want line 1, got %d", toks[0].Line)
	}
	if toks[1].Line != 2 {
		t.Errorf("second token: want line 2, got %d", toks[1].Line)
	}
	if toks[2].Line != 3 {
		t.Errorf("third token: want line 3, got %d", toks[2].Line)
	}
}
