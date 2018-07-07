package ast

type Kind int

const (
	// UNKNOWN es utilizado por defecto:
	UNKNOWN Kind = iota
	STR
	INT
)

// AST representa el AST.
type AST struct {
	GlobalScope *Scope
}

// New inicializa un nuevo AST.
func New() *AST {
	ast := &AST{
		GlobalScope: &Scope{
			OuterScope: nil,
			Objects:    make(map[string]Node),
			Nodes:      make([]Node, 0),
		},
	}
	return ast
}

// Node representa un nodo en el AST.
type Node interface {
	GetValue() interface{}
}

type NodeStruct struct {
	Children []*Node
	Parent   *Node
	Node
}

// Scope representa un scope.
type Scope struct {
	OuterScope *Scope
	Objects    map[string]Node
	Nodes      []Node
	Node
}

// BlockStmt representa un bloque de sentencias.
type BlockStmt struct {
	*Scope
	Node
}

// Stmt representa una sentencia.
type Stmt struct {
	Name  string
	Value string
	Node
}

// Decl representa una declaración de variables.
type Decl struct {
	Name  string
	Value interface{}
	Node
}

type BinaryExpr struct {
	X  Node
	Y  Node
	Op string
	Node
}

type AssignStmt struct {
	X Node
	Y Node
	Node
}

type Ident struct {
	Name string
	Node
}

func (i *Ident) GetValue() interface{} {
	return i.Name
}

type BasicLit struct {
	Kind    Kind
	Literal string
	Node
}

func (l *BasicLit) GetValue() interface{} {
	return l.Literal
}

func (d *Decl) GetValue() interface{} {
	return d.Value
}
