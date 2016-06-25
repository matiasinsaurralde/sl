package Ast

import(
  "os"
)

type File struct {
  Name string
  File *os.File
}

type FuncDeclaration struct {
  Name string
}
