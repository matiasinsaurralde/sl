package token

type Token int

type Pos int

const (
	_ = iota
	DUMMY
	EOL
	COMMENT
	COMMENT_START
	COMMENT_END

	IDENT
	INT
	STRING

	ASSIGN

	PROGRAM
	PRINT

	START
	END
	SUBR
	SUBR_NAME
	SUBR_RETURN
	SUBR_RETURN_TYPE

	VAR
	VAR_NAME
	VAR_TYPE
	VAR_VALUE

	LPAREN
	RPAREN

	EXPR
	OP
)

var tokens = [...]string{
	DUMMY:            "DUMMY",
	EOL:              "EOL",
	COMMENT:          "COMMENT",
	COMMENT_START:    "/*",
	COMMENT_END:      "*/",
	IDENT:            "IDENT",
	INT:              "INT",
	STRING:           "STRING",
	ASSIGN:           "=",
	PROGRAM:          "programa",
	PRINT:            "imprimir",
	START:            "inicio",
	END:              "fin",
	SUBR:             "subrutina",
	SUBR_NAME:        "SUBR_NAME",
	SUBR_RETURN:      "retorna",
	SUBR_RETURN_TYPE: "SUBR_RETURN_TYPE",
	VAR:              "var",
	VAR_NAME:         "VAR_NAME",
	VAR_TYPE:         "VAR_TYPE",
	VAR_VALUE:        "VAR_VALUE",
	LPAREN:           "(",
	RPAREN:           ")",
	EXPR:             "EXPR",
	OP:               "OP",
}

func Lookup(input string) Token {
	for i, v := range tokens {
		if v == input {
			return Token(i)
			break
		}
	}
	return -1
}

func Get(index int) string {
	for k, v := range tokens {
		if index == k {
			return v
			break
		}
	}
	return ""
}
