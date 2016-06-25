package parser

import(
  "github.com/matiasinsaurralde/sl/ast"

  "os"
)

func ParseFile( filename string ) ( f *Ast.File, err error ) {
  
  var file *os.File

  file, err = os.Open( filename )

  f = &Ast.File{
    Name: filename,
    File: file,
  }

  return f, err
}
