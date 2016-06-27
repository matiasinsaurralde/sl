package parser_test

import(
  "github.com/matiasinsaurralde/sl/parser"
  "github.com/matiasinsaurralde/sl/token"
  "github.com/matiasinsaurralde/sl/ast"
  // "fmt"
  "testing"
)

func assertLiteral( literal *Ast.BasicLiteral, kind token.Token, t *testing.T ) {
  if literal.Kind != kind {
    t.Error( "Expected literal to be of type", kind, "got", literal.Kind )
  }
}

func assertOperator( x *Ast.BinaryExpression, op string, t *testing.T ) {
  if x.Operator != op {
    t.Error( "Expected operator to be", op, "got", x.Operator )
  }
}

func TestBasicIntLiteralEvaluation( t *testing.T ) {
  source := `1`
  var basicLiteral *Ast.BasicLiteral
  basicLiteral = parser.Eval(source, nil).(*Ast.BasicLiteral)
  assertLiteral( basicLiteral, token.INT, t)
}

func TestBasicBigIntLiteralEvaluation( t *testing.T ) {
  source := `199999999`
  var basicLiteral *Ast.BasicLiteral
  basicLiteral = parser.Eval(source, nil).(*Ast.BasicLiteral)
  assertLiteral( basicLiteral, token.INT, t )
}

func TestBasicStringLiteralEvaluation( t *testing.T ) {
  source := `test`
  var basicLiteral *Ast.BasicLiteral
  basicLiteral = parser.Eval(source, nil).(*Ast.BasicLiteral)
  assertLiteral( basicLiteral, token.STRING, t )
}

func TestBasicIntBinaryExpr( t *testing.T ) {
  source := `10+20`
  var expr *Ast.BinaryExpression
  expr = parser.Eval(source, nil).(*Ast.BinaryExpression)

  var xLiteral, yLiteral *Ast.BasicLiteral
  xLiteral = expr.X.(*Ast.BasicLiteral)
  yLiteral = expr.Y.(*Ast.BasicLiteral)

  assertLiteral( xLiteral, token.INT, t )
  assertLiteral( yLiteral, token.INT, t )

}

func TestBasicStringBinaryExpr( t *testing.T ) {
  source := `a+b`
  var expr *Ast.BinaryExpression
  expr = parser.Eval(source, nil).(*Ast.BinaryExpression)

  var xLiteral, yLiteral *Ast.BasicLiteral
  xLiteral = expr.X.(*Ast.BasicLiteral)
  yLiteral = expr.Y.(*Ast.BasicLiteral)

  assertLiteral( xLiteral, token.STRING, t )
  assertLiteral( yLiteral, token.STRING, t )

}

func TestBasicMixedBinaryExpr( t *testing.T ) {
  source := `a+1`
  var expr *Ast.BinaryExpression
  expr = parser.Eval(source, nil).(*Ast.BinaryExpression)

  var xLiteral, yLiteral *Ast.BasicLiteral
  xLiteral = expr.X.(*Ast.BasicLiteral)
  yLiteral = expr.Y.(*Ast.BasicLiteral)

  assertLiteral( xLiteral, token.STRING, t )
  assertLiteral( yLiteral, token.INT, t )

}

func TestComplexMixedBinaryExpr( t *testing.T ) {
  source := `a+1-b*2`
  var expr, xExpr, yExpr *Ast.BinaryExpression
  expr = parser.Eval(source, nil).(*Ast.BinaryExpression)

  xExpr = expr.X.(*Ast.BinaryExpression)
  yExpr = expr.Y.(*Ast.BinaryExpression)

  var xxLiteral, xyLiteral, yxLiteral, yyLiteral *Ast.BasicLiteral
  xxLiteral = xExpr.X.(*Ast.BasicLiteral)
  xyLiteral = xExpr.Y.(*Ast.BasicLiteral)
  yxLiteral = yExpr.X.(*Ast.BasicLiteral)
  yyLiteral = yExpr.Y.(*Ast.BasicLiteral)

  assertLiteral( xxLiteral, token.STRING, t )
  assertLiteral( xyLiteral, token.INT, t )
  assertLiteral( yxLiteral, token.STRING, t )
  assertLiteral( yyLiteral, token.INT, t )

  assertOperator( expr, "-", t )
  assertOperator( xExpr, "+", t )
  assertOperator( yExpr, "*", t )

}
