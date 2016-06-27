package parser

import(
  "strings"
  "strconv"

  "github.com/matiasinsaurralde/sl/ast"
  "github.com/matiasinsaurralde/sl/token"
)

func Eval(x string, rootExpression Ast.Expression ) Ast.Expression {
  var ex Ast.Expression

  reader := strings.NewReader( x )

  var binaryExpr bool = false

  for {
    b, err := reader.ReadByte()
    ch := string(b)

    if isOperator(ch) {
      binaryExpr = true
      break
    }

    if err != nil {
      break
    }
  }

  reader.Seek( 0, 0 )

  if binaryExpr {
    ex = EvalBinaryExpr( x, rootExpression )
  } else {
    ex = EvalLiteral( x, rootExpression  )
  }

  return ex
}

func EvalBinaryExpr( x string, rootExpression Ast.Expression ) Ast.Expression {
  reader := strings.NewReader( x )

  var e Ast.Expression
  var be Ast.BinaryExpression

  be = Ast.BinaryExpression{}

  buf := make( []byte, 0 )

  operatorCount := 0

  for {
    b, err := reader.ReadByte()
    ch := string(b)

    buf = append( buf, b )
    stringBuffer := string(buf)

    if isOperator(ch) {
      operatorCount++
    }

    if isOperator(ch) && be.X == nil {
      be.Operator = ch
      be.X = EvalLiteral( stringBuffer, nil )
      buf = make([]byte, 0)
    }

    if isOperator(ch) && operatorCount > 1 {
      be.Y = EvalLiteral( stringBuffer, nil  )
      read := reader.Size() - int64(reader.Len())
      s := x[ read : len( x ) ]
      var x Ast.Expression
      x = &be
      subexpr := Ast.BinaryExpression{
        Operator: ch,
        X: x,
        Y: Eval(s, rootExpression),
      }
      e = &subexpr
      break
    }

    if err != nil {
      be.Y = EvalLiteral( stringBuffer, nil  )
      e = &be
      break
    }
  }

  return e
}

func EvalLiteral( x string, rootExpression Ast.Expression ) Ast.Expression {
  reader := strings.NewReader( x )

  buf := make( []byte, 0 )

  var tok Ast.Expression

  for {
    b, err := reader.ReadByte()
    ch := string(b)

    if isLetter(ch) || isNumber(ch) {
      buf = append( buf, b )
    }

    if err != nil {

      stringBuffer := string(buf)

      if isNumber(stringBuffer) {
        tok = &Ast.BasicLiteral{
          Value: stringBuffer,
          Kind: token.INT,
        }
      } else {
        tok = &Ast.BasicLiteral{
          Value: stringBuffer,
          Kind: token.STRING,
        }
      }

      break
    }
  }
  return tok
}

func isNumber(input string) bool {
  _, err := strconv.ParseInt( input, 10, 32 )
  return err == nil
}

func isOperator(input string) bool {
  operators := "+-*"
  return strings.Contains(operators, input)
}

func isLetter(input string) bool {
  abc := "abcdefghijklmnopqrstuvwxyz"
  return strings.Contains(abc, input)
}
