// Package ast defines the abstract syntax tree for the SL language.
package ast

import "github.com/matiasinsaurralde/sl/pkg/lexer"

// Node is the base interface for all AST nodes.
type Node interface {
	nodeTag()
}

// ---- Type nodes ----

// TypeNode represents a type expression.
type TypeNode interface {
	Node
	typeTag()
}

// SimpleType is numerico, cadena, or logico.
type SimpleType struct {
	Name string // "numerico", "cadena", "logico"
}

func (*SimpleType) nodeTag() {}
func (*SimpleType) typeTag() {}

// VectorType is vector [N] T or vector [*] T.
type VectorType struct {
	Size     int // 0 means open (*)
	ElemType TypeNode
}

func (*VectorType) nodeTag() {}
func (*VectorType) typeTag() {}

// MatrixType is matriz [d1, d2, ...] T where 0 means open (*).
type MatrixType struct {
	Dims     []int // 0 = open (*)
	ElemType TypeNode
}

func (*MatrixType) nodeTag() {}
func (*MatrixType) typeTag() {}

// RegistroType is registro { fields }.
type RegistroType struct {
	Fields []*FieldDef
}

func (*RegistroType) nodeTag() {}
func (*RegistroType) typeTag() {}

// FieldDef is a single field in a registro.
type FieldDef struct {
	Names []string
	Type  TypeNode
}

// NamedType is a reference to a user-defined tipo.
type NamedType struct {
	Name string
}

func (*NamedType) nodeTag() {}
func (*NamedType) typeTag() {}

// ---- Declaration nodes ----

// VarDecl declares one variable (or group) at a given scope.
type VarDecl struct {
	Names []string
	Type  TypeNode // nil if inferred
	Init  Expr     // nil if no initializer
}

func (*VarDecl) nodeTag() {}

// ConstDecl declares a named constant.
type ConstDecl struct {
	Name string
	Init Expr
}

func (*ConstDecl) nodeTag() {}

// TiposDecl declares a named type alias.
type TiposDecl struct {
	Name string
	Type TypeNode
}

func (*TiposDecl) nodeTag() {}

// SubDecl declares a subroutine or function.
type SubDecl struct {
	Name       string
	Params     []*ParamGroup
	ReturnType TypeNode // nil for procedures
	Consts     []*ConstDecl
	Types      []*TiposDecl
	Vars       []*VarDecl
	Body       []Stmt
}

func (*SubDecl) nodeTag() {}

// ParamGroup is a group of params with the same ref-ness and type.
type ParamGroup struct {
	ByRef bool
	Names []string
	Type  TypeNode
}

// ---- Program ----

// Program is the root AST node.
type Program struct {
	Name   string // from "programa <name>", may be empty
	Consts []*ConstDecl
	Types  []*TiposDecl
	Vars   []*VarDecl
	Body   []Stmt
	Subs   []*SubDecl
}

func (*Program) nodeTag() {}

// ---- Statement nodes ----

// Stmt is the base interface for statements.
type Stmt interface {
	Node
	stmtTag()
}

// AssignStmt is target = value.
type AssignStmt struct {
	Line   int
	Target Expr // LValue
	Value  Expr
}

func (*AssignStmt) nodeTag() {}
func (*AssignStmt) stmtTag() {}

// CallStmt is a subroutine call as a statement.
type CallStmt struct {
	Line int
	Call *CallExpr
}

func (*CallStmt) nodeTag() {}
func (*CallStmt) stmtTag() {}

// ImprimirStmt is imprimir(args...).
type ImprimirStmt struct {
	Line int
	Args []Expr
}

func (*ImprimirStmt) nodeTag() {}
func (*ImprimirStmt) stmtTag() {}

// LeerStmt is leer(vars...).
type LeerStmt struct {
	Line int
	Vars []Expr // lvalue expressions
}

func (*LeerStmt) nodeTag() {}
func (*LeerStmt) stmtTag() {}

// SiStmt is si (cond) { ... [sino si (cond) ...]* [sino ...] }.
type SiStmt struct {
	Line    int
	Cond    Expr
	Then    []Stmt
	ElseIfs []*ElseIf
	Else    []Stmt // nil if no sino
}

func (*SiStmt) nodeTag() {}
func (*SiStmt) stmtTag() {}

// ElseIf is a sino si (cond) branch.
type ElseIf struct {
	Cond Expr
	Body []Stmt
}

// DesdeStmt is desde var=start hasta end [paso step] { body }.
type DesdeStmt struct {
	Line  int
	Var   string
	Start Expr
	End   Expr
	Step  Expr // nil means step = 1
	Body  []Stmt
}

func (*DesdeStmt) nodeTag() {}
func (*DesdeStmt) stmtTag() {}

// MientrasStmt is mientras (cond) { body }.
type MientrasStmt struct {
	Line int
	Cond Expr
	Body []Stmt
}

func (*MientrasStmt) nodeTag() {}
func (*MientrasStmt) stmtTag() {}

// RepetirStmt is repetir body hasta (cond).
type RepetirStmt struct {
	Line int
	Body []Stmt
	Cond Expr
}

func (*RepetirStmt) nodeTag() {}
func (*RepetirStmt) stmtTag() {}

// EvalStmt is eval { caso (bool) stmts ... [sino stmts] }.
type EvalStmt struct {
	Line  int
	Cases []*EvalCase
	Else  []Stmt
}

func (*EvalStmt) nodeTag() {}
func (*EvalStmt) stmtTag() {}

// EvalCase is a single caso (bool) stmts block.
type EvalCase struct {
	Cond Expr
	Body []Stmt
}

// SalirStmt exits the innermost loop.
type SalirStmt struct {
	Line int
}

func (*SalirStmt) nodeTag() {}
func (*SalirStmt) stmtTag() {}

// RetornaStmt returns a value from a function.
type RetornaStmt struct {
	Line  int
	Value Expr // nil for void
}

func (*RetornaStmt) nodeTag() {}
func (*RetornaStmt) stmtTag() {}

// TerminarStmt terminates the program.
type TerminarStmt struct {
	Line int
	Msg  Expr // optional
}

func (*TerminarStmt) nodeTag() {}
func (*TerminarStmt) stmtTag() {}

// ---- Expression nodes ----

// Expr is the base interface for expressions.
type Expr interface {
	Node
	exprTag()
	GetLine() int
}

// BinaryExpr is left op right.
type BinaryExpr struct {
	Line  int
	Op    lexer.TokenType
	Left  Expr
	Right Expr
}

func (*BinaryExpr) nodeTag()       {}
func (*BinaryExpr) exprTag()       {}
func (e *BinaryExpr) GetLine() int { return e.Line }

// UnaryExpr is op operand.
type UnaryExpr struct {
	Line    int
	Op      lexer.TokenType
	Operand Expr
}

func (*UnaryExpr) nodeTag()       {}
func (*UnaryExpr) exprTag()       {}
func (e *UnaryExpr) GetLine() int { return e.Line }

// IdentExpr is a simple variable or constant reference.
type IdentExpr struct {
	Line int
	Name string
}

func (*IdentExpr) nodeTag()       {}
func (*IdentExpr) exprTag()       {}
func (e *IdentExpr) GetLine() int { return e.Line }

// NumberLit is a numeric literal.
type NumberLit struct {
	Line  int
	Value float64
}

func (*NumberLit) nodeTag()       {}
func (*NumberLit) exprTag()       {}
func (e *NumberLit) GetLine() int { return e.Line }

// StringLit is a string literal.
type StringLit struct {
	Line  int
	Value string
}

func (*StringLit) nodeTag()       {}
func (*StringLit) exprTag()       {}
func (e *StringLit) GetLine() int { return e.Line }

// BoolLit is TRUE/FALSE/SI/NO.
type BoolLit struct {
	Line  int
	Value bool
}

func (*BoolLit) nodeTag()       {}
func (*BoolLit) exprTag()       {}
func (e *BoolLit) GetLine() int { return e.Line }

// ArrayLit is { expr, expr, ..., ... }.
type ArrayLit struct {
	Line  int
	Elems []Expr // each elem may itself be an ArrayLit for matrices
	Fill  bool   // true if last element is "..." (fill rest with last)
}

func (*ArrayLit) nodeTag()       {}
func (*ArrayLit) exprTag()       {}
func (e *ArrayLit) GetLine() int { return e.Line }

// CallExpr is a function/subroutine call.
type CallExpr struct {
	Line int
	Name string
	Args []Expr
}

func (*CallExpr) nodeTag()       {}
func (*CallExpr) exprTag()       {}
func (e *CallExpr) GetLine() int { return e.Line }

// IndexExpr is arr[i] or arr[i, j, ...].
type IndexExpr struct {
	Line    int
	Array   Expr
	Indices []Expr
}

func (*IndexExpr) nodeTag()       {}
func (*IndexExpr) exprTag()       {}
func (e *IndexExpr) GetLine() int { return e.Line }

// FieldExpr is record.field.
type FieldExpr struct {
	Line   int
	Record Expr
	Field  string
}

func (*FieldExpr) nodeTag()       {}
func (*FieldExpr) exprTag()       {}
func (e *FieldExpr) GetLine() int { return e.Line }
