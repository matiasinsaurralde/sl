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
}

type ExpressionStatement struct {

}

type CallExpression struct {
  Function string
  Args []Expression
  StartPos token.Pos
  EndPos token.Pos
}

func( c *CallExpression ) Pos() token.Pos {
  return c.StartPos
}

func( c *CallExpression ) End() token.Pos {
  return c.EndPos
}

type PrintStatement struct {
  Print token.Pos
}

type Comment struct {
  Text string
  StartPos token.Pos
  EndPos token.Pos
}

type GenericDeclaration struct {
  Name string
  Values []Expression
  StartPos token.Pos
  EndPos token.Pos
}

func( r *GenericDeclaration ) Pos() token.Pos {
  return r.StartPos
}

func( r *GenericDeclaration ) End() token.Pos {
  return r.EndPos
}


type BasicLiteral struct {
  Value string
  Kind token.Token
  StartPos token.Pos
  EndPos token.Pos
}

func( b *BasicLiteral ) Pos() token.Pos {
  return b.StartPos
}

func( b *BasicLiteral ) End() token.Pos {
  return b.EndPos
}

type BinaryExpression struct {
  X Expression
  StartPos token.Pos
  EndPos token.Pos
  Y Expression
}

func( b *BinaryExpression ) Pos() token.Pos {
  return b.StartPos
}

func( b *BinaryExpression ) End() token.Pos {
  return b.EndPos
}

type Expression interface {
  Node
}

type Variable struct {
  Name string
  Value interface{}
}

func( c *Comment ) Pos() token.Pos {
  return c.StartPos
}

func (c *Comment ) End() token.Pos {
  return c.EndPos
}

type Declaration struct {
  StartPos token.Pos
  EndPos token.Pos
}

func( d *Declaration ) Pos() token.Pos {
  return d.StartPos
}

func( d *Declaration ) End() token.Pos {
  return d.EndPos
}

type Scope struct {
  Declarations []Node
  Nodes []Node
}

type FuncDeclaration struct {
  Name string
}
