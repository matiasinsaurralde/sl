package lexer

// TokenType identifies the kind of lexical token.
type TokenType int

const (
	// Special
	ILLEGAL  TokenType = iota
	EOF                // end of file
	NUMBER             // numeric literal: 3, 3.14
	STRING             // string literal: "hello" or 'world'
	IDENT              // identifier
	ELLIPSIS           // ...

	// Keywords
	AND
	ARCHIVO
	CASO
	CONST
	CONSTANTES
	DESDE
	EVAL
	FIN
	HASTA
	INICIO
	LIB
	LIBEXT
	LOGICO
	MATRIZ
	MIENTRAS
	NOT
	NUMERICO
	OR
	PASO
	PROGRAMA
	CADENA
	REF
	REGISTRO
	REPETIR
	RETORNA
	SALIR
	SI
	SINO
	SUB
	SUBRUTINA
	TIPOS
	VAR
	VARIABLES
	VECTOR

	// Operators
	PLUS    // +
	MINUS   // -
	STAR    // *
	SLASH   // /
	PERCENT // %
	CARET   // ^
	EQ      // ==
	NEQ     // <>
	LT      // <
	LE      // <=
	GT      // >
	GE      // >=
	ASSIGN  // =
	AND2    // &&
	OR2     // ||

	// Delimiters
	LPAREN // (
	RPAREN // )
	LBRACK // [
	RBRACK // ]
	LBRACE // {
	RBRACE // }
	COMMA  // ,
	SEMI   // ;
	DOT    // .
	COLON  // :
)

var tokenNames = map[TokenType]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	NUMBER:     "NUMBER",
	STRING:     "STRING",
	IDENT:      "IDENT",
	ELLIPSIS:   "...",
	AND:        "and",
	ARCHIVO:    "archivo",
	CASO:       "caso",
	CONST:      "const",
	CONSTANTES: "constantes",
	DESDE:      "desde",
	EVAL:       "eval",
	FIN:        "fin",
	HASTA:      "hasta",
	INICIO:     "inicio",
	LIB:        "lib",
	LIBEXT:     "libext",
	LOGICO:     "logico",
	MATRIZ:     "matriz",
	MIENTRAS:   "mientras",
	NOT:        "not",
	NUMERICO:   "numerico",
	OR:         "or",
	PASO:       "paso",
	PROGRAMA:   "programa",
	CADENA:     "cadena",
	REF:        "ref",
	REGISTRO:   "registro",
	REPETIR:    "repetir",
	RETORNA:    "retorna",
	SALIR:      "salir",
	SI:         "si",
	SINO:       "sino",
	SUB:        "sub",
	SUBRUTINA:  "subrutina",
	TIPOS:      "tipos",
	VAR:        "var",
	VARIABLES:  "variables",
	VECTOR:     "vector",
	PLUS:       "+",
	MINUS:      "-",
	STAR:       "*",
	SLASH:      "/",
	PERCENT:    "%",
	CARET:      "^",
	EQ:         "==",
	NEQ:        "<>",
	LT:         "<",
	LE:         "<=",
	GT:         ">",
	GE:         ">=",
	ASSIGN:     "=",
	AND2:       "&&",
	OR2:        "||",
	LPAREN:     "(",
	RPAREN:     ")",
	LBRACK:     "[",
	RBRACK:     "]",
	LBRACE:     "{",
	RBRACE:     "}",
	COMMA:      ",",
	SEMI:       ";",
	DOT:        ".",
	COLON:      ":",
}

func (t TokenType) String() string {
	if s, ok := tokenNames[t]; ok {
		return s
	}
	return "UNKNOWN"
}

var keywords = map[string]TokenType{
	"and":        AND,
	"archivo":    ARCHIVO,
	"caso":       CASO,
	"const":      CONST,
	"constantes": CONSTANTES,
	"desde":      DESDE,
	"eval":       EVAL,
	"fin":        FIN,
	"hasta":      HASTA,
	"inicio":     INICIO,
	"lib":        LIB,
	"libext":     LIBEXT,
	"logico":     LOGICO,
	"matriz":     MATRIZ,
	"mientras":   MIENTRAS,
	"not":        NOT,
	"numerico":   NUMERICO,
	"or":         OR,
	"paso":       PASO,
	"programa":   PROGRAMA,
	"cadena":     CADENA,
	"ref":        REF,
	"registro":   REGISTRO,
	"repetir":    REPETIR,
	"retorna":    RETORNA,
	"salir":      SALIR,
	"si":         SI,
	"sino":       SINO,
	"sub":        SUB,
	"subrutina":  SUBRUTINA,
	"tipos":      TIPOS,
	"var":        VAR,
	"variables":  VARIABLES,
	"vector":     VECTOR,
}

// LookupIdent returns the keyword TokenType for s, or IDENT.
func LookupIdent(s string) TokenType {
	if t, ok := keywords[s]; ok {
		return t
	}
	return IDENT
}

// Token is a single lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

func (t Token) String() string {
	return t.Literal
}
