package token

type Token int

type Pos int

const(
  _ = iota
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

  VAR

  LPAREN
  RPAREN

  EXPR
)

var tokens = [...]string{
  EOL: "EOL",
  COMMENT: "COMMENT",
  COMMENT_START: "/*",
  COMMENT_END: "*/",
  IDENT: "IDENT",
  INT: "INT",
  STRING: "STRING",
  ASSIGN: "=",
  PROGRAM: "programa",
  PRINT: "imprimir",
  START: "inicio",
  END: "fin",
  VAR: "var",
  LPAREN: "(",
  RPAREN: ")",
  EXPR: "EXPR",
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
