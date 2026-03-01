package parser

import (
	"testing"

	"github.com/matiasinsaurralde/sl/pkg/ast"
)

func mustParse(t *testing.T, src string) *ast.Program {
	t.Helper()
	prog, errs := Parse(src)
	if len(errs) > 0 {
		t.Fatalf("parse errors for %q: %v", src, errs)
	}
	return prog
}

func TestParseEmptyProgram(t *testing.T) {
	prog := mustParse(t, "inicio\nfin")
	if len(prog.Body) != 0 {
		t.Errorf("expected empty body, got %d statements", len(prog.Body))
	}
}

func TestParseVarDecl(t *testing.T) {
	prog := mustParse(t, `
var
   x : numerico
   y : cadena
inicio
fin`)
	if len(prog.Vars) != 2 {
		t.Fatalf("expected 2 var decls, got %d", len(prog.Vars))
	}
}

func TestParseConstDecl(t *testing.T) {
	prog := mustParse(t, `
const
   PI = 3.14
inicio
fin`)
	if len(prog.Consts) != 1 {
		t.Fatalf("expected 1 const, got %d", len(prog.Consts))
	}
	if prog.Consts[0].Name != "PI" {
		t.Errorf("expected const name PI, got %q", prog.Consts[0].Name)
	}
}

func TestParseImprimirStmt(t *testing.T) {
	prog := mustParse(t, `inicio
   imprimir ("hello", "\n")
fin`)
	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(prog.Body))
	}
	if _, ok := prog.Body[0].(*ast.ImprimirStmt); !ok {
		t.Errorf("expected ImprimirStmt, got %T", prog.Body[0])
	}
}

func TestParseSiStmt(t *testing.T) {
	prog := mustParse(t, `
var x : numerico
inicio
   si (x > 0) {
      imprimir ("pos")
   sino
      imprimir ("neg")
   }
fin`)
	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(prog.Body))
	}
	si, ok := prog.Body[0].(*ast.SiStmt)
	if !ok {
		t.Fatalf("expected SiStmt, got %T", prog.Body[0])
	}
	if si.Else == nil {
		t.Error("expected else branch")
	}
}

func TestParseDesdeStmt(t *testing.T) {
	prog := mustParse(t, `
var i : numerico
inicio
   desde i=1 hasta 10 {
      imprimir (i)
   }
fin`)
	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(prog.Body))
	}
	if _, ok := prog.Body[0].(*ast.DesdeStmt); !ok {
		t.Errorf("expected DesdeStmt, got %T", prog.Body[0])
	}
}

func TestParseMientrasStmt(t *testing.T) {
	prog := mustParse(t, `
var x : numerico
inicio
   mientras (x < 10) {
      x = x + 1
   }
fin`)
	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(prog.Body))
	}
	if _, ok := prog.Body[0].(*ast.MientrasStmt); !ok {
		t.Errorf("expected MientrasStmt, got %T", prog.Body[0])
	}
}

func TestParseRepetirStmt(t *testing.T) {
	prog := mustParse(t, `
var x : numerico
inicio
   repetir {
      x = x + 1
   } hasta (x >= 5)
fin`)
	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(prog.Body))
	}
	if _, ok := prog.Body[0].(*ast.RepetirStmt); !ok {
		t.Errorf("expected RepetirStmt, got %T", prog.Body[0])
	}
}

func TestParseSubDecl(t *testing.T) {
	prog := mustParse(t, `
inicio
   dbl(3)
fin

sub dbl(x : numerico)
inicio
   retorna x * 2
fin`)
	if len(prog.Subs) != 1 {
		t.Fatalf("expected 1 sub, got %d", len(prog.Subs))
	}
	sub := prog.Subs[0]
	if sub.Name != "dbl" {
		t.Errorf("expected sub name dbl, got %q", sub.Name)
	}
	if len(sub.Params) != 1 {
		t.Fatalf("expected 1 param group, got %d", len(sub.Params))
	}
}

func TestParseSubRefParam(t *testing.T) {
	prog := mustParse(t, `
var x : numerico
inicio
   inc(x)
fin

sub inc(ref v : numerico)
inicio
   v = v + 1
fin`)
	if len(prog.Subs) != 1 {
		t.Fatalf("expected 1 sub, got %d", len(prog.Subs))
	}
	sub := prog.Subs[0]
	if !sub.Params[0].ByRef {
		t.Error("expected first param to be ByRef")
	}
}

func TestParseVectorType(t *testing.T) {
	prog := mustParse(t, `
var v : vector [10] numerico
inicio
fin`)
	if len(prog.Vars) != 1 {
		t.Fatalf("expected 1 var, got %d", len(prog.Vars))
	}
	vt, ok := prog.Vars[0].Type.(*ast.VectorType)
	if !ok {
		t.Fatalf("expected VectorType, got %T", prog.Vars[0].Type)
	}
	if vt.Size != 10 {
		t.Errorf("expected size 10, got %d", vt.Size)
	}
}

func TestParseMatrixTypeNoElemType(t *testing.T) {
	// Element type is optional — should default without parse error
	prog := mustParse(t, `
var M : matriz [5, 3]
inicio
fin`)
	if len(prog.Vars) != 1 {
		t.Fatalf("expected 1 var, got %d", len(prog.Vars))
	}
	mt, ok := prog.Vars[0].Type.(*ast.MatrixType)
	if !ok {
		t.Fatalf("expected MatrixType, got %T", prog.Vars[0].Type)
	}
	if len(mt.Dims) != 2 || mt.Dims[0] != 5 || mt.Dims[1] != 3 {
		t.Errorf("expected dims [5,3], got %v", mt.Dims)
	}
}

func TestParseRetornaNoParens(t *testing.T) {
	prog := mustParse(t, `
inicio
   dbl(2)
fin

sub dbl(x : numerico)
inicio
   retorna x * 2
fin`)
	sub := prog.Subs[0]
	if len(sub.Body) != 1 {
		t.Fatalf("expected 1 stmt in sub body, got %d", len(sub.Body))
	}
	ret, ok := sub.Body[0].(*ast.RetornaStmt)
	if !ok {
		t.Fatalf("expected RetornaStmt, got %T", sub.Body[0])
	}
	if ret.Value == nil {
		t.Error("expected non-nil return value")
	}
}

func TestParseEvalCaso(t *testing.T) {
	prog := mustParse(t, `
var x : numerico
inicio
   eval {
      caso (x = 1): imprimir ("one")
      caso (x = 2): imprimir ("two")
   }
fin`)
	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(prog.Body))
	}
	ev, ok := prog.Body[0].(*ast.EvalStmt)
	if !ok {
		t.Fatalf("expected EvalStmt, got %T", prog.Body[0])
	}
	if len(ev.Cases) != 2 {
		t.Errorf("expected 2 cases, got %d", len(ev.Cases))
	}
}

func TestParseRegistroType(t *testing.T) {
	prog := mustParse(t, `
tipos
   punto = registro {
      x, y : numerico
   }
var p : punto
inicio
fin`)
	if len(prog.Types) != 1 {
		t.Fatalf("expected 1 tipo, got %d", len(prog.Types))
	}
}

func TestParseArrayLitWithEllipsis(t *testing.T) {
	prog := mustParse(t, `
var M : matriz [5, 3] = {{1, 2, 3}, ...}
inicio
fin`)
	if len(prog.Vars) != 1 {
		t.Fatalf("expected 1 var, got %d", len(prog.Vars))
	}
	// Just ensure it parses without error
}

func TestParseSmartQuoteString(t *testing.T) {
	// Smart quotes as used in the example files
	src := "inicio\n   imprimir (\u201CHola\u201D)\nfin"
	prog := mustParse(t, src)
	if len(prog.Body) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(prog.Body))
	}
}
