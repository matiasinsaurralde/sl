package Ast

import(
  "github.com/matiasinsaurralde/sl/token"
  "os"
)

type File struct {
  Name string
  ProgramName string
  File *os.File
  Scope *Scope

  Comments []Comment
}

type Node interface {
  Pos() token.Pos
  End() token.Pos
}

type RoutineLiteral struct {
  // Type *RoutineType
  // Body BlockStatement
  Body *BlockStatement

  StartPos token.Pos
  EndPos token.Pos
}

func( r *RoutineLiteral ) Pos() token.Pos {
  return r.StartPos
}

func( r *RoutineLiteral ) End() token.Pos {
  return r.EndPos
}

type RoutineType struct {

}

type BlockStatement struct {
  Start token.Pos
  List []Statement
  End token.Pos
}

type Statement interface {
  Node
  StatementNode()
}

type PrintStatement struct {
  Print token.Pos
}

type Comment struct {
  Text string
  StartPos token.Pos
  EndPos token.Pos
}

func( c *Comment ) Pos() token.Pos {
  return c.StartPos
}

func (c *Comment ) End() token.Pos {
  return c.EndPos
}

type Scope struct {
}

type FuncDeclaration struct {
  Name string
}
