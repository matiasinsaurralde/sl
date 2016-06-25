package parser

import(
  "github.com/matiasinsaurralde/sl/ast"
  "github.com/matiasinsaurralde/sl/token"

  "os"
  "bufio"

  "strings"

  "fmt"
)

func parseBlockStatement( st *Ast.BlockStatement, body string ) {
  fmt.Println( "ParseBlockStatement:", st, "\n")
  reader := strings.NewReader(body)
  scanner := bufio.NewScanner( reader )

  for scanner.Scan() {
    fmt.Println("bufio.Scanner, statement:", scanner.Text() )
  }
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

  fmt.Println("Parser initialization")

  var node Ast.Node

  for {
    _, err := r.Read(buf)

    var keep bool
    keep = false

    // fmt.Println("Read ", n, "bytes, buf: ", buf, string(buf), string(v))

    // fmt.Println("expect", expect, string(v))

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
        default:
          v = make([]byte, 0)
        }
      }
    case 10: // "\n"
      tok := token.Lookup(string(v))
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

          parseBlockStatement(routineLiteral.Body, body )

          expect = -1
          v = make([]byte, 0)
        } else {
          keep = true
        }
      }

      switch tok {
      case token.START:
        expect = token.END
        v = make([]byte,0)

        routineLiteral := &Ast.RoutineLiteral{
          StartPos: token.Pos(currentPos),
          Body: &Ast.BlockStatement{
            List: make([]Ast.Statement, 0),
          },
        }

        node = routineLiteral
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

  return f, err
}
