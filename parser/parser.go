package parser

import(
  "github.com/matiasinsaurralde/sl/ast"
  "github.com/matiasinsaurralde/sl/token"

  "os"
  "bufio"

  "fmt"
)

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

  var expect token.Token
  expect = -1

  fmt.Println("Parser initialization")

  for {
    _, err := r.Read(buf)

    var keep bool
    keep = false

    // fmt.Println("Read ", n, "bytes, buf: ", buf, string(buf), string(v))

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
        default:
          v = make([]byte, 0)
        }
      }
    case 10: // "\n"
      switch expect {
      case token.IDENT:
        fmt.Println("\n- Found a program: ", string(v))
        f.ProgramName = string(v)
        expect = -1
      case token.COMMENT_END:
        keep = true
      }
    case 47: // "/"
      if expect == token.COMMENT_END {
        // fmt.Println("Comment:", string(v))
        fmt.Println( "\n- Found a comment...", string(v))
        f.Comments = append( f.Comments, Ast.Comment{string(v)} )
      }
      if expect == -1 {
        v = make([]byte, 0)
        expect = token.COMMENT_START
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
    /*
    switch buf[0] {
    case 32:
      fmt.Println("v = ", v, string(v))
      v = make([]byte, 0)
    case 34:
      fmt.Println("String starts or ends", string(v))
      v = make([]byte, 0)
    case 40:
      fmt.Println("Initial parenthesis", string(v))
      v = make([]byte, 0)
    case 41:
      fmt.Println("End parenthesis", string(v))
      v = make([]byte,0)
    case 10:
      fmt.Println("CRLF", v)
      v = make([]byte,0)
    default:
      v = append(v, buf[0])
    }
    */

    if err != nil {
      break
    }
  }

  return f, err
}
