package parser

import(
  "github.com/matiasinsaurralde/sl/ast"
  "github.com/matiasinsaurralde/sl/token"

  "github.com/davecgh/go-spew/spew"
  // goparser "go/parser"
  // goast "go/ast"

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
      // fmt.Println(" * Statement:", rawStatement )
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

            if rawExpressions == "" {}

            var CallExpression *Ast.CallExpression
            CallExpression = node.(*Ast.CallExpression)

            // parseExpressions(&CallExpression.Args, rawExpressions)

            var statement Ast.Statement
            statement = CallExpression

            bs.List = append( bs.List, statement)

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
    // fmt.Println(" * Node:", node)

    declarations = append(declarations, node)

    return declarations
  }

  splits = strings.Split(body, "=")

  // x = value (literal)
  if len(splits) == 2 {
    expressions := make([]Ast.Expression, 0)
    genericDeclaration := &Ast.GenericDeclaration{
      Name: splits[0],
    }

    node = genericDeclaration
    // parseExpressions( &expressions, splits[1])

    if expressions == nil {}


    declarations = append(declarations, node)

    return declarations
  }

  return declarations
}

func Parse( input string ) (f *Ast.File, err error) {

  f = &Ast.File{
    Comments: make([]Ast.Comment, 0),
    Nodes: make([]Ast.Node, 0),
  }

  reader := strings.NewReader(input)

  buf := make([]byte, 0)

  currentPosition := 0

  var tok, expect token.Token

  var block string = ""

  expect = -1

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

    switch tok {
    case token.VAR:
      expect = token.VAR_NAME
      buf = make([]byte, 0)
    case token.START:
      expect = token.END
      buf = make([]byte, 0)
      ignore = true
      block = ""

      if subroutine {
        // fmt.Println(" - Block starts (subroutine)...")
      } else {
        // fmt.Println(" - Block starts (main)...")
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
        // fmt.Println("VAR_NAME =",stringBuf)
        buf = make([]byte, 0)
        expect = token.VAR_TYPE
        ignore = true
      }
      if ch == "=" {
        genericDeclaration = Ast.GenericDeclaration{
          Name: stringBuf,
          StartPos: token.Pos( currentPosition - len( stringBuf ) ),
        }
        // fmt.Println("VAR_NAME =",stringBuf)
        buf = make([]byte, 0)
        expect = token.VAR_VALUE
        ignore = true
      }
    case token.VAR_TYPE:
      if ch == "\n" || err == io.EOF {
        genericDeclaration.EndPos = token.Pos(currentPosition + len(stringBuf) )
        declaration := genericDeclaration

        // fmt.Println( "*** genericDeclaration:", genericDeclaration)
        // fmt.Println("VAR_TYPE =", stringBuf)

        f.Nodes = append( f.Nodes , &declaration )

        buf = make([]byte, 0)
        expect = token.VAR_NAME
        ignore = true
      }
    case token.VAR_VALUE:
      if ch == "\n" || err == io.EOF {
        genericDeclaration.EndPos = token.Pos(currentPosition + len(stringBuf) )
        declaration := genericDeclaration

        genericDeclaration.Values = Eval(stringBuf, nil )

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

  return f, err
}
