package token

// Token es utilizado para enumerar los tokens.
type Token int

const (
	// UNKNOWN es utilizado por defecto:
	UNKNOWN Token = iota
	NIL

	// EOF representa el final del archivo.
	EOF

	// EQ representa =
	EQ

	// COLON representa :
	COLON

	// COMMA representa ,
	COMMA

	// IDENT representa un identificador.
	IDENT
	TYPE

	// EXPR representa una expresión.
	EXPR
	OP

	COMMENT

	PARENEXPR
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET

	BLOCK

	// LITINT representa un literal entero.
	LITINT
	LITFLOAT

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

	KWDESDE
	KWHASTA
	KWPASO
)

// Item almacena información de cada token (posición y tipo).
type Item struct {
	Type    Token
	Literal string

	Pos int64
	Ln  int64

	Skip bool
}

// IsOp indica si el token es un operador.
func (i *Item) IsOp() (ok bool) {
	switch i.Literal {
	case "+":
		ok = true
	case "-":
		ok = true
	case "*":
		ok = true
	case "/":
		ok = true
	}
	return ok
}
