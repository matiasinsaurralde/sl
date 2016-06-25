package parser

import(
  "github.com/matiasinsaurralde/sl/ast"
  "github.com/matiasinsaurralde/sl/token"

  "github.com/davecgh/go-spew/spew"
  // goparser "go/parser"
  // goast "go/ast"

  // "os"
  "bufio"
  "io"

  "strings"

  "fmt"
)

func parseBlockStatement( body *string ) *Ast.BlockStatement {
  reader := strings.NewReader(*body)
  scanner := bufio.NewScanner( reader )

  bs := Ast.BlockStatement{}

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

            bs.List = append( bs.List, statement)

            // blockStatement.List = append(blockStatement.List, statement)

            expect = -1
          }
        default:
          buf = append(buf, byte(ch))
        }
      }
    }
  }

  return &bs
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

/*
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
*/

func Parse( input string ) (f *Ast.File, err error) {

  f = &Ast.File{
    Comments: make([]Ast.Comment, 0),
    Nodes: make([]Ast.Node, 0),
  }

  reader := strings.NewReader(input)

  buf := make([]byte, 0)

  currentPosition := 0

  spew.Dump(1)

  var tok, expect token.Token

  var block string = ""

  expect = -1

  // var node Ast.Node
  var genericDeclaration Ast.GenericDeclaration
  var subroutineDeclaration Ast.SubroutineDeclaration
  var mainDeclaration Ast.MainDeclaration
  var subroutine bool = false

  for {

    var ignore bool = false

    b, err := reader.ReadByte()

    ch := string(b)
    stringBuf := string(buf)

    tok = token.Lookup(stringBuf)

    // spew.Dump(stringBuf)
    // spew.Dump(stringBuf, expect)
    // fmt.Println("")

    // spew.Dump(stringBuf, tok)

    switch tok {
    case token.VAR:
      fmt.Println("*** Definicion", stringBuf, expect, ch)
      expect = token.VAR_NAME
      buf = make([]byte, 0)
    case token.START:
      fmt.Println("bexpect", expect)
      expect = token.END
      buf = make([]byte, 0)
      ignore = true
      block = ""

      if subroutine {
        fmt.Println(" - Block starts (subroutine)...")
      } else {
        fmt.Println(" - Block starts (main)...")
        mainDeclaration = Ast.MainDeclaration{
          StartPos: token.Pos(currentPosition),
        }
      }

    case token.SUBR:
      expect = token.SUBR_NAME
      buf = make([]byte, 0)
    case token.SUBR_RETURN:
      expect = token.SUBR_RETURN_TYPE
      buf = make([]byte, 0)
    }


    switch expect {
    case token.VAR_NAME:
      if ch == "\n" {
        ignore = true
      }
      if ch == ":" {
        genericDeclaration = Ast.GenericDeclaration{
          Name: stringBuf,
          StartPos: token.Pos( currentPosition - len( stringBuf ) ),
        }
        fmt.Println("VAR_NAME =",stringBuf)
        buf = make([]byte, 0)
        expect = token.VAR_TYPE
        ignore = true
      }
      if ch == "=" {
        genericDeclaration = Ast.GenericDeclaration{
          Name: stringBuf,
          StartPos: token.Pos( currentPosition - len( stringBuf ) ),
        }
        fmt.Println("VAR_NAME =",stringBuf)
        buf = make([]byte, 0)
        expect = token.VAR_VALUE
        ignore = true
      }
    case token.VAR_TYPE:
      if ch == "\n" || err == io.EOF {
        genericDeclaration.EndPos = token.Pos(currentPosition + len(stringBuf) )
        declaration := genericDeclaration

        fmt.Println( "*** genericDeclaration:", genericDeclaration)
        fmt.Println("VAR_TYPE =", stringBuf)

        f.Nodes = append( f.Nodes , &declaration )

        buf = make([]byte, 0)
        expect = token.VAR_NAME
        ignore = true
      }
    case token.VAR_VALUE:
      if ch == "\n" || err == io.EOF {
        genericDeclaration.EndPos = token.Pos(currentPosition + len(stringBuf) )
        declaration := genericDeclaration

        fmt.Println( "*** genericDeclaration:", genericDeclaration)
        fmt.Println("VAR_VALUE =", stringBuf)

        f.Nodes = append( f.Nodes , &declaration )

        buf = make([]byte, 0)
        expect = token.VAR_NAME
        ignore = true
      }
    case token.START:
      // fmt.Println("Expecting START", stringBuf )
      if ch == "\n" {
        ignore = true
      }

      if len(stringBuf) == 1 {
        startLinebreak := stringBuf[0:1]
        if startLinebreak == "\n" {
          buf = buf[ 1 : len(buf) ]
        }
      }
    case token.END:
      length := len(stringBuf)
      if( length >= 4 ) {

        lastChars := stringBuf[length-4 : length]

        if lastChars == "fin\n" {
          block = stringBuf[0 : length - 4]
          fmt.Println(" - Block ends with contents:")

          if subroutine {
            subroutineDeclaration.EndPos = token.Pos( currentPosition - len(stringBuf))
            subroutineDeclaration.Body = parseBlockStatement(&block)
            fmt.Println( "*** subroutineDeclaration", subroutineDeclaration )
            declaration := subroutineDeclaration
            f.Nodes = append( f.Nodes, &declaration )
          } else {
            mainDeclaration.EndPos = token.Pos( currentPosition )
            mainDeclaration.Body = parseBlockStatement(&block)
            declaration := mainDeclaration
            f.Nodes = append( f.Nodes, &declaration )
          }

          spew.Dump(block)
          fmt.Println("")
          buf = make([]byte, 0)
          subroutine = false
          expect = -1
        }
      }
      // lastChars := stringBuf[ len(stringBuf) - 3 : len(stringBuf) ]
      // fmt.Println( "***", lastChars)
    case token.SUBR_NAME:
      charTok := token.Lookup(ch)
      switch charTok {
      case token.LPAREN:
        fmt.Println(" - Subroutine is declared:", stringBuf )

        subroutineDeclaration = Ast.SubroutineDeclaration{
          Name: stringBuf,
          StartPos: token.Pos(currentPosition - len(stringBuf) ),
        }

        subroutine = true
        ignore = true
        buf = make([]byte, 0)
      case token.RPAREN:
        fmt.Println(" - Subroutine declaration with contents:")
        spew.Dump(block)
        fmt.Println("\n")
        ignore = true
        buf = make([]byte, 0)
        expect = token.START
      }
      // spew.Dump(ch)
    case token.SUBR_RETURN_TYPE:
      if ch == "\n" {
        buf = make([]byte, 0)
        fmt.Println( " - Subroutine returns: *", stringBuf, "*" )
        expect = token.START
      }
    case -1:
      if ch == "\n" {
        ignore = true
      }
      if len(stringBuf) == 1 {
        startLinebreak := stringBuf[0:1]
        if startLinebreak == "\n" {
          buf = buf[ 1 : len(buf) ]
        }
      }
    }

    if ch == " " {
      ignore = true
    }

    /*
    if len(stringBuf) > 2 {
      if stringBuf[0:1] == "\n" {
        buf = buf[1:len(buf)]
        ignore = true
      }
    }
    */

    if !ignore {
      buf = append( buf, b )
    }

    if err != nil {
      if err == io.EOF {
        fmt.Println("*eof")
      }
      break
    }

    currentPosition++

  }

  // fmt.Println("buf", string(buf))

  return f, err
}

func ParseExpression(x string) ( expr Ast.Expression, err error ) {
  // node := Ast.BasicLiteral{}

  buf := []string{}
  for i, c := range x {
    fmt.Println( i, c )
    ch := string(c)
    appendToBuffer := false
    switch ch {
    case " ":
      appendToBuffer = false
    default:
      appendToBuffer = true
    }

    if appendToBuffer {
      buf = append(buf, ch)
    }
  }

  fmt.Println("Buffer:", buf, len(buf))
  // expr = &node
  return expr, err
}
