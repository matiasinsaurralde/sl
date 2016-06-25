package parser

import(
  "github.com/matiasinsaurralde/sl/ast"
  "github.com/matiasinsaurralde/sl/token"

  "os"
  "bufio"

  "strings"

  "fmt"
)

func parseBlockStatement( blockStatement *Ast.BlockStatement, body string ) {
  reader := strings.NewReader(body)
  scanner := bufio.NewScanner( reader )

  var expect token.Token
  expect = -1

  var node Ast.Node

  for scanner.Scan() {
    rawStatement := scanner.Text()
    rawStatement = strings.Replace(rawStatement, "\n", "", -1)
    if len(rawStatement) > 0 {
      fmt.Println(" * Statement:", rawStatement )
      buf := make([]byte, 0)

      for _, ch := range rawStatement {

        switch ch {
        case 32: // " "
        case 40: // (
          tok := token.Lookup(string(buf))

          switch tok {
          case token.PRINT:
            CallExpression := &Ast.CallExpression{
              Function: string(buf),
            }
            node = CallExpression
            expect = token.EXPR
            buf = make([]byte, 0)
          default:
            buf = make([]byte, 0)
          }
        case 41: // ")"
          if expect == token.EXPR {
            rawExpressions := string(buf)

            var CallExpression *Ast.CallExpression
            CallExpression = node.(*Ast.CallExpression)

            parseExpressions(&CallExpression.Args, rawExpressions)

            var statement Ast.Statement
            statement = CallExpression

            blockStatement.List = append(blockStatement.List, statement)

            expect = -1
          }
        default:
          buf = append(buf, byte(ch))
        }
      }
    }
  }
}

func parseExpression( expr string ) Ast.Expression {
  // e := Ast.Expression{}
  var e Ast.Expression
  if strings.Index( expr, "+") > 0 {

    splits := strings.Split(expr, "+")

    b := &Ast.BinaryExpression{
      X: &Ast.BasicLiteral{Value: splits[0]},
      Y: &Ast.BasicLiteral{Value: splits[1]},
    }
    e = b

    // b.X = literalX
    // b.Y = literalY
    return e
  }

  // Basic literal?
  b := &Ast.BasicLiteral{
    Value: expr,
  }
  return b
}

func parseExpressions( expressions *[]Ast.Expression, expr string ) {
  e := parseExpression(expr)
  *expressions = append( *expressions, e )
}

func parseDeclarations( body string ) []Ast.Node {

  declarations := make([]Ast.Node, 0)

  body = strings.Replace(body, "\n", "", -1)
  body = strings.Replace(body, " ", "", -1)

  if len(body) == 0 {
    return declarations
  }

  var splits []string

  splits = strings.Split(body, ":")

  var node Ast.Node

  // x : type (no value)
  if len(splits) == 2 {
    node = &Ast.GenericDeclaration{
      Name: splits[0],
    }
    fmt.Println(" * Node:", node)

    declarations = append(declarations, node)

    return declarations
  }

  splits = strings.Split(body, "=")

  // x = value (literal)
  if len(splits) == 2 {
    expressions := make([]Ast.Expression, 0)
    genericDeclaration := &Ast.GenericDeclaration{
      Name: splits[0],
      Values: make([]Ast.Expression, 0),
    }

    node = genericDeclaration

    fmt.Println(" * Node:", node)

    parseExpressions( &expressions, splits[1])
    // fmt.Println("Expressions:", splits[1], expressions[0])
    // fmt.Println("* Node:", node)

    declarations = append(declarations, node)

    return declarations
  }

  return declarations
}

func ParseFile( filename string ) ( f *Ast.File, err error ) {

  var file *os.File

  file, err = os.Open( filename )

  f = &Ast.File{
    Name: filename,
    File: file,
    Comments: make([]Ast.Comment, 0),
  }

  r := bufio.NewReader(file)
  buf := make([]byte, 1)
  v := make([]byte, 0)

  startPos := 0
  currentPos := 0

  var expect token.Token
  expect = -1

  fmt.Println("Parse...")

  var node Ast.Node

  var globalScope Ast.Scope
  globalScope = Ast.Scope{}

  for {
    _, err := r.Read(buf)

    var keep bool
    keep = false

    switch buf[0] {
    case 32: // " "

      tok := token.Lookup(string(v))

      if expect == -1 {
        switch tok {
        case token.PROGRAM:
          v = make([]byte, 0)
          expect = token.IDENT
        default:
          v = make([]byte, 0)
        }
      } else {
        switch expect {
        case token.COMMENT_END:
          keep = true
        case token.END:
          keep = true
        case token.VAR:
          // fmt.Println("?Var", string(v), len(v))
          keep = true
        default:
          v = make([]byte, 0)
        }
      }
    case 10: // "\n"

      value := string(v)
      value = strings.Replace(value, "\n", "", -1)

      tok := token.Lookup(value)

      switch expect {
      case token.IDENT:
        fmt.Println("\n- Found a program: ", string(v))
        f.ProgramName = string(v)
        expect = -1
      case token.COMMENT_END:
        keep = true
      case token.END:

        s := string(v)

        endSt := s[ len(s)-3 : len(s) ]
        tok = token.Lookup(string(endSt))

        if tok == token.END {

          routineLiteral := node.(*Ast.RoutineLiteral)
          routineLiteral.EndPos = token.Pos(currentPos)

          body := s[0:len(s)-3]

          parseBlockStatement( routineLiteral.Body, body )

          expect = -1
          v = make([]byte, 0)

          // fmt.Println("\nBlock ends")

        } else {
          keep = true
        }
      case token.VAR:
        keep = true
        // fmt.Println("Var declaration...", string(v))
        declarations := parseDeclarations(string(v))

        for _, d := range declarations {
          globalScope.Declarations = append( globalScope.Declarations, d )
        }

        v = make([]byte, 0)
      }

      switch tok {
      case token.START:
        fmt.Println("\n- Found a block\n")

        expect = token.END
        v = make([]byte,0)

        routineLiteral := &Ast.RoutineLiteral{
          StartPos: token.Pos(currentPos),
          Body: &Ast.BlockStatement{
            List: make([]Ast.Statement, 0),
          },
        }

        node = routineLiteral

        globalScope.Nodes = append(globalScope.Nodes, node)
      case token.VAR:
        fmt.Println("- Declaration:\n")
        v = make([]byte, 0)
        keep = true
        expect = token.VAR
      }
    case 47: // "/"
      if expect == -1 {
        v = make([]byte, 0)
        expect = token.COMMENT_START

        startPos = currentPos
      }
      if expect == token.COMMENT_END {
        fmt.Println( "\n- Found a comment...", string(v))
        comment := Ast.Comment{
          StartPos: token.Pos(startPos), EndPos: token.Pos(currentPos), Text: string(v),
        }
        f.Comments = append( f.Comments, comment )
        v = make([]byte, 0)
        expect = -1
      }
    case 42: // "*"
      if expect == token.COMMENT_START {
        expect = token.COMMENT_END
        v = make([]byte, 0)
      }
    default:
      v = append(v, buf[0])
    }

    if keep {
      v = append(v, buf[0])
    }

    currentPos++

    if err != nil {
      break
    }
  }

  f.Scope  = &globalScope

  fmt.Println( "\n- Scope:", f.Scope, "\n")

  return f, err
}
