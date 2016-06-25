package Ast

import(
  "os"
)

type File struct {
  Name string
  ProgramName string
  File *os.File
  Scope *Scope

  Comments []Comment
}

type Comment struct {
  Text string
}

type Scope struct {
}

type FuncDeclaration struct {
  Name string
}
