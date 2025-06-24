package token

type Token int

type Pos int

const (
	_ = iota
	ILLEGAL
	EOF
	COMMENT

	// Literals
	IDENT
	INT
	FLOAT
	STRING

	// Operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /
	MODULO   // %

	// Comparison operators
	EQ  // ==
	NEQ // <>
	LT  // <
	LTE // <=
	GT  // >
	GTE // >=

	// Logical operators
	AND // &&
	OR  // or

	// Delimiters
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :

	// Keywords
	PROGRAMA // programa
	VAR      // var
	CONST    // const
	INICIO   // inicio
	FIN      // fin
	SUBR     // sub
	RETORNA  // retorna
	NUMERICO // numerico

	// Control flow
	SI       // si
	SINO     // sino
	MIENTRAS // mientras
	REPETIR  // repetir
	HASTA    // hasta
	DESDE    // desde
	PASO     // paso
	TERMINAR // terminar

	// Built-in functions
	IMPRIMIR // imprimir
	LEER     // leer
	INT_FUNC // int
	IFVAL    // ifval

	// Special
	EOL // end of line
)

var tokens = [...]string{
	ILLEGAL:   "ILLEGAL",
	EOF:       "EOF",
	COMMENT:   "COMMENT",
	IDENT:     "IDENT",
	INT:       "INT",
	FLOAT:     "FLOAT",
	STRING:    "STRING",
	ASSIGN:    "=",
	PLUS:      "+",
	MINUS:     "-",
	MULTIPLY:  "*",
	DIVIDE:    "/",
	MODULO:    "%",
	EQ:        "==",
	NEQ:       "<>",
	LT:        "<",
	LTE:       "<=",
	GT:        ">",
	GTE:       ">=",
	AND:       "&&",
	OR:        "or",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACKET:  "[",
	RBRACKET:  "]",
	COMMA:     ",",
	SEMICOLON: ";",
	COLON:     ":",
	PROGRAMA:  "programa",
	VAR:       "var",
	CONST:     "const",
	INICIO:    "inicio",
	FIN:       "fin",
	SUBR:      "sub",
	RETORNA:   "retorna",
	NUMERICO:  "numerico",
	SI:        "si",
	SINO:      "sino",
	MIENTRAS:  "mientras",
	REPETIR:   "repetir",
	HASTA:     "hasta",
	DESDE:     "desde",
	PASO:      "paso",
	TERMINAR:  "terminar",
	IMPRIMIR:  "imprimir",
	LEER:      "leer",
	INT_FUNC:  "int",
	IFVAL:     "ifval",
	EOL:       "EOL",
}

func Lookup(input string) Token {
	for i, v := range tokens {
		if v == input {
			return Token(i)
		}
	}
	return IDENT // Default to identifier if not found
}

func Get(index int) string {
	if index >= 0 && index < len(tokens) {
		return tokens[index]
	}
	return ""
}

func (t Token) String() string {
	return Get(int(t))
}
