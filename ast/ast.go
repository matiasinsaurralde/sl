package ast

import (
	"os"

	"github.com/matiasinsaurralde/sl/token"
)

// File represents the root of an SL program
type File struct {
	Name        string
	ProgramName string
	File        *os.File
	Scope       *Scope
	Nodes       []Node

	Comments []Comment
}

// Node is the interface that all AST nodes implement
type Node interface {
	Pos() token.Pos
	End() token.Pos
}

// Statement represents a statement node
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node
type Expression interface {
	Node
	expressionNode()
}

// Comment represents a comment in the source code
type Comment struct {
	Text     string
	StartPos token.Pos
	EndPos   token.Pos
}

func (c *Comment) Pos() token.Pos { return c.StartPos }
func (c *Comment) End() token.Pos { return c.EndPos }

// Program represents a program declaration
type Program struct {
	Name     string
	StartPos token.Pos
	EndPos   token.Pos
}

func (p *Program) Pos() token.Pos { return p.StartPos }
func (p *Program) End() token.Pos { return p.EndPos }
func (p *Program) statementNode() {}

// VariableDeclaration represents a variable declaration
type VariableDeclaration struct {
	Name     string
	Type     string
	Value    Expression
	StartPos token.Pos
	EndPos   token.Pos
}

func (v *VariableDeclaration) Pos() token.Pos { return v.StartPos }
func (v *VariableDeclaration) End() token.Pos { return v.EndPos }
func (v *VariableDeclaration) statementNode() {}

// SubroutineDeclaration represents a subroutine/function declaration
type SubroutineDeclaration struct {
	Name       string
	Parameters []*Parameter
	ReturnType string
	Body       *BlockStatement
	StartPos   token.Pos
	EndPos     token.Pos
}

func (s *SubroutineDeclaration) Pos() token.Pos { return s.StartPos }
func (s *SubroutineDeclaration) End() token.Pos { return s.EndPos }
func (s *SubroutineDeclaration) statementNode() {}

// Parameter represents a function parameter
type Parameter struct {
	Name     string
	Type     string
	StartPos token.Pos
	EndPos   token.Pos
}

func (p *Parameter) Pos() token.Pos { return p.StartPos }
func (p *Parameter) End() token.Pos { return p.EndPos }

// BlockStatement represents a block of statements
type BlockStatement struct {
	Statements []Statement
	StartPos   token.Pos
	EndPos     token.Pos
}

func (b *BlockStatement) Pos() token.Pos { return b.StartPos }
func (b *BlockStatement) End() token.Pos { return b.EndPos }
func (b *BlockStatement) statementNode() {}

// ExpressionStatement represents an expression used as a statement
type ExpressionStatement struct {
	Expression Expression
	StartPos   token.Pos
	EndPos     token.Pos
}

func (e *ExpressionStatement) Pos() token.Pos { return e.StartPos }
func (e *ExpressionStatement) End() token.Pos { return e.EndPos }
func (e *ExpressionStatement) statementNode() {}

// IfStatement represents an if-else statement
type IfStatement struct {
	Condition Expression
	Then      Statement
	Else      Statement
	StartPos  token.Pos
	EndPos    token.Pos
}

func (i *IfStatement) Pos() token.Pos { return i.StartPos }
func (i *IfStatement) End() token.Pos { return i.EndPos }
func (i *IfStatement) statementNode() {}

// WhileStatement represents a while loop
type WhileStatement struct {
	Condition Expression
	Body      Statement
	StartPos  token.Pos
	EndPos    token.Pos
}

func (w *WhileStatement) Pos() token.Pos { return w.StartPos }
func (w *WhileStatement) End() token.Pos { return w.EndPos }
func (w *WhileStatement) statementNode() {}

// ForStatement represents a for loop (desde-hasta-paso)
type ForStatement struct {
	Variable string
	Start    Expression
	EndExpr  Expression
	Step     Expression
	Body     Statement
	StartPos token.Pos
	EndPos   token.Pos
}

func (f *ForStatement) Pos() token.Pos { return f.StartPos }
func (f *ForStatement) End() token.Pos { return f.EndPos }
func (f *ForStatement) statementNode() {}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Value    Expression
	StartPos token.Pos
	EndPos   token.Pos
}

func (r *ReturnStatement) Pos() token.Pos { return r.StartPos }
func (r *ReturnStatement) End() token.Pos { return r.EndPos }
func (r *ReturnStatement) statementNode() {}

// TerminateStatement represents a terminate statement
type TerminateStatement struct {
	Message  Expression
	StartPos token.Pos
	EndPos   token.Pos
}

func (t *TerminateStatement) Pos() token.Pos { return t.StartPos }
func (t *TerminateStatement) End() token.Pos { return t.EndPos }
func (t *TerminateStatement) statementNode() {}

// Identifier represents an identifier
type Identifier struct {
	Name     string
	StartPos token.Pos
	EndPos   token.Pos
}

func (i *Identifier) Pos() token.Pos  { return i.StartPos }
func (i *Identifier) End() token.Pos  { return i.EndPos }
func (i *Identifier) expressionNode() {}

// Literal represents a literal value
type Literal struct {
	Type     token.Token
	Value    string
	StartPos token.Pos
	EndPos   token.Pos
}

func (l *Literal) Pos() token.Pos  { return l.StartPos }
func (l *Literal) End() token.Pos  { return l.EndPos }
func (l *Literal) expressionNode() {}

// BinaryExpression represents a binary operation
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
	StartPos token.Pos
	EndPos   token.Pos
}

func (b *BinaryExpression) Pos() token.Pos  { return b.StartPos }
func (b *BinaryExpression) End() token.Pos  { return b.EndPos }
func (b *BinaryExpression) expressionNode() {}

// CallExpression represents a function call
type CallExpression struct {
	Function  string
	Arguments []Expression
	StartPos  token.Pos
	EndPos    token.Pos
}

func (c *CallExpression) Pos() token.Pos  { return c.StartPos }
func (c *CallExpression) End() token.Pos  { return c.EndPos }
func (c *CallExpression) expressionNode() {}

// AssignmentExpression represents an assignment
type AssignmentExpression struct {
	Left     *Identifier
	Operator string
	Right    Expression
	StartPos token.Pos
	EndPos   token.Pos
}

func (a *AssignmentExpression) Pos() token.Pos  { return a.StartPos }
func (a *AssignmentExpression) End() token.Pos  { return a.EndPos }
func (a *AssignmentExpression) expressionNode() {}

// IfValExpression represents an ifval expression
type IfValExpression struct {
	Condition Expression
	Then      Expression
	Else      Expression
	StartPos  token.Pos
	EndPos    token.Pos
}

func (i *IfValExpression) Pos() token.Pos  { return i.StartPos }
func (i *IfValExpression) End() token.Pos  { return i.EndPos }
func (i *IfValExpression) expressionNode() {}

// Scope represents a variable scope
type Scope struct {
	Variables map[string]interface{}
	Parent    *Scope
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Variables: make(map[string]interface{}),
		Parent:    parent,
	}
}

func (s *Scope) Get(name string) (interface{}, bool) {
	if value, exists := s.Variables[name]; exists {
		return value, true
	}
	if s.Parent != nil {
		return s.Parent.Get(name)
	}
	return nil, false
}

func (s *Scope) Set(name string, value interface{}) {
	s.Variables[name] = value
}
