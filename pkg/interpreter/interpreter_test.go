package interpreter

import (
	"strings"
	"testing"

	"github.com/matiasinsaurralde/sl/pkg/parser"
)

// run executes src with the given stdin lines (joined by newlines) and returns
// the captured stdout. It fatally fails if parsing or execution fails.
func run(t *testing.T, src, stdin string) string {
	t.Helper()
	prog, errs := parser.Parse(src)
	if len(errs) > 0 {
		t.Fatalf("parse errors: %v", errs)
	}

	interp := New(nil)

	// Redirect stdout to a strings.Builder
	var out strings.Builder
	interp.stdout.w = &out

	// Provide stdin by writing to a temp reader
	if stdin != "" {
		interp.stdin.buf = splitByIFS(stdin+"\n", interp.ifs)
	}

	if err := interp.Run(prog); err != nil {
		if _, ok := err.(*TerminateError); ok {
			return out.String()
		}
		t.Fatalf("runtime error: %v", err)
	}
	return out.String()
}

// ---- Basic output ----

func TestImprimir(t *testing.T) {
	out := run(t, `inicio
   imprimir ("hello world")
fin`, "")
	if out != "hello world" {
		t.Errorf("want %q, got %q", "hello world", out)
	}
}

func TestImprimirMultiArg(t *testing.T) {
	out := run(t, `inicio
   imprimir ("a", " ", "b")
fin`, "")
	if out != "a b" {
		t.Errorf("want %q, got %q", "a b", out)
	}
}

// ---- Arithmetic ----

func TestArithmetic(t *testing.T) {
	cases := []struct {
		expr string
		want string
	}{
		{"2 + 3", "5"},
		{"10 - 3", "7"},
		{"4 * 5", "20"},
		{"10 / 4", "2.5"},
		{"7 % 3", "1"},
		{"2 ^ 8", "256"},
	}
	for _, c := range cases {
		src := `var x : numerico
inicio
   x = ` + c.expr + `
   imprimir (x)
fin`
		out := run(t, src, "")
		if out != c.want {
			t.Errorf("expr %q: want %q, got %q", c.expr, c.want, out)
		}
	}
}

// ---- String concatenation ----

func TestStringConcat(t *testing.T) {
	out := run(t, `var s : cadena
inicio
   s = "foo" + "bar"
   imprimir (s)
fin`, "")
	if out != "foobar" {
		t.Errorf("want %q, got %q", "foobar", out)
	}
}

// ---- Conditionals ----

func TestSiTrue(t *testing.T) {
	out := run(t, `inicio
   si (1 < 2) {
      imprimir ("yes")
   }
fin`, "")
	if out != "yes" {
		t.Errorf("want yes, got %q", out)
	}
}

func TestSiElse(t *testing.T) {
	out := run(t, `inicio
   si (2 < 1) {
      imprimir ("yes")
   sino
      imprimir ("no")
   }
fin`, "")
	if out != "no" {
		t.Errorf("want no, got %q", out)
	}
}

// ---- Loops ----

func TestDesdeLoop(t *testing.T) {
	out := run(t, `var i : numerico
inicio
   desde i=1 hasta 5 {
      imprimir (i, " ")
   }
fin`, "")
	if out != "1 2 3 4 5 " {
		t.Errorf("want %q, got %q", "1 2 3 4 5 ", out)
	}
}

func TestDesdePaso(t *testing.T) {
	out := run(t, `var i : numerico
inicio
   desde i=0 hasta 10 paso 2 {
      imprimir (i, " ")
   }
fin`, "")
	if out != "0 2 4 6 8 10 " {
		t.Errorf("want %q, got %q", "0 2 4 6 8 10 ", out)
	}
}

func TestMientrasLoop(t *testing.T) {
	out := run(t, `var n : numerico
inicio
   n = 1
   mientras (n <= 3) {
      imprimir (n, " ")
      n = n + 1
   }
fin`, "")
	if out != "1 2 3 " {
		t.Errorf("want %q, got %q", "1 2 3 ", out)
	}
}

func TestRepetirLoop(t *testing.T) {
	out := run(t, `var n : numerico
inicio
   n = 1
   repetir {
      imprimir (n, " ")
      n = n + 1
   } hasta (n > 3)
fin`, "")
	if out != "1 2 3 " {
		t.Errorf("want %q, got %q", "1 2 3 ", out)
	}
}

// ---- Sub calls ----

func TestSubCall(t *testing.T) {
	out := run(t, `inicio
   imprimir (double(4))
fin

sub double(x : numerico)
inicio
   retorna x * 2
fin`, "")
	if out != "8" {
		t.Errorf("want 8, got %q", out)
	}
}

func TestSubRefParam(t *testing.T) {
	out := run(t, `var x : numerico
inicio
   x = 5
   bump(x)
   imprimir (x)
fin

sub bump(ref v : numerico)
inicio
   v = v + 1
fin`, "")
	if out != "6" {
		t.Errorf("want 6, got %q", out)
	}
}

func TestRecursiveFactorial(t *testing.T) {
	out := run(t, `inicio
   imprimir (fact(5))
fin

sub fact(n : numerico)
inicio
   si (n <= 1) {
      retorna 1
   }
   retorna n * fact(n - 1)
fin`, "")
	if out != "120" {
		t.Errorf("want 120, got %q", out)
	}
}

// ---- Arrays ----

func TestVectorAssignRead(t *testing.T) {
	out := run(t, `var v : vector [3] numerico
inicio
   v[1] = 10
   v[2] = 20
   v[3] = 30
   imprimir (v[1], " ", v[2], " ", v[3])
fin`, "")
	if out != "10 20 30" {
		t.Errorf("want %q, got %q", "10 20 30", out)
	}
}

func TestAlen(t *testing.T) {
	out := run(t, `var v : vector [5] numerico
inicio
   imprimir (alen(v))
fin`, "")
	if out != "5" {
		t.Errorf("want 5, got %q", out)
	}
}

func TestDimOpenVector(t *testing.T) {
	out := run(t, `var v : vector [*] numerico
inicio
   dim(v, 4)
   imprimir (alen(v))
fin`, "")
	if out != "4" {
		t.Errorf("want 4, got %q", out)
	}
}

func TestMatrixAssignRead(t *testing.T) {
	out := run(t, `var M : matriz [2, 3] numerico
inicio
   M[1,1] = 1
   M[1,2] = 2
   M[1,3] = 3
   M[2,1] = 4
   M[2,2] = 5
   M[2,3] = 6
   imprimir (M[2,3])
fin`, "")
	if out != "6" {
		t.Errorf("want 6, got %q", out)
	}
}

// ---- Eval/Caso ----

func TestEvalCaso(t *testing.T) {
	out := run(t, `var x : numerico
inicio
   x = 2
   eval {
      caso (x = 1): imprimir ("one")
      caso (x = 2): imprimir ("two")
      caso (x = 3): imprimir ("three")
   }
fin`, "")
	if out != "two" {
		t.Errorf("want two, got %q", out)
	}
}

func TestEvalSinoNone(t *testing.T) {
	out := run(t, `var x : numerico
inicio
   x = 99
   eval {
      caso (x = 1): imprimir ("one")
      sino: imprimir ("other")
   }
fin`, "")
	if out != "other" {
		t.Errorf("want other, got %q", out)
	}
}

// ---- Built-in functions ----

func TestAbs(t *testing.T) {
	out := run(t, `inicio
   imprimir (abs(-7))
fin`, "")
	if out != "7" {
		t.Errorf("want 7, got %q", out)
	}
}

func TestSqrt(t *testing.T) {
	out := run(t, `inicio
   imprimir (sqrt(9))
fin`, "")
	if out != "3" {
		t.Errorf("want 3, got %q", out)
	}
}

func TestStrlen(t *testing.T) {
	out := run(t, `inicio
   imprimir (strlen("hello"))
fin`, "")
	if out != "5" {
		t.Errorf("want 5, got %q", out)
	}
}

func TestUpper(t *testing.T) {
	out := run(t, `inicio
   imprimir (upper("hello"))
fin`, "")
	if out != "HELLO" {
		t.Errorf("want HELLO, got %q", out)
	}
}

func TestInc(t *testing.T) {
	out := run(t, `var n : numerico
inicio
   n = 5
   inc(n)
   imprimir (n)
fin`, "")
	if out != "6" {
		t.Errorf("want 6, got %q", out)
	}
}

func TestDec(t *testing.T) {
	out := run(t, `var n : numerico
inicio
   n = 5
   dec(n, 2)
   imprimir (n)
fin`, "")
	if out != "3" {
		t.Errorf("want 3, got %q", out)
	}
}

func TestIntercambiar(t *testing.T) {
	out := run(t, `var a, b : numerico
inicio
   a = 1
   b = 2
   intercambiar(a, b)
   imprimir (a, " ", b)
fin`, "")
	if out != "2 1" {
		t.Errorf("want %q, got %q", "2 1", out)
	}
}

func TestTerminar(t *testing.T) {
	out := run(t, `inicio
   terminar ("done")
   imprimir ("never")
fin`, "")
	if !strings.Contains(out, "done") {
		t.Errorf("want output containing 'done', got %q", out)
	}
	if strings.Contains(out, "never") {
		t.Errorf("should not print 'never', got %q", out)
	}
}

// ---- Registro ----

func TestRegistro(t *testing.T) {
	out := run(t, `tipos
   punto = registro {
      x, y : numerico
   }
var p : punto
inicio
   p.x = 3
   p.y = 4
   imprimir (p.x, " ", p.y)
fin`, "")
	if out != "3 4" {
		t.Errorf("want %q, got %q", "3 4", out)
	}
}

// ---- Salir (break) ----

func TestSalir(t *testing.T) {
	out := run(t, `var i : numerico
inicio
   desde i=1 hasta 10 {
      si (i = 4) {
         salir
      }
      imprimir (i, " ")
   }
fin`, "")
	if out != "1 2 3 " {
		t.Errorf("want %q, got %q", "1 2 3 ", out)
	}
}

// ---- Leer ----

func TestLeer(t *testing.T) {
	out := run(t, `var a, b : numerico
inicio
   leer (a, b)
   imprimir (a + b)
fin`, "3,5")
	if out != "8" {
		t.Errorf("want 8, got %q", out)
	}
}

// ---- GCD example ----

func TestGCDExample(t *testing.T) {
	out := run(t, `var a, b : numerico
inicio
   leer (a, b)
   mientras (a <> b) {
      si (a > b) {
         a = a - b
      sino
         b = b - a
      }
   }
   imprimir (a)
fin`, "6,9")
	if out != "3" {
		t.Errorf("want 3, got %q", out)
	}
}

// ---- Matrix transpose ----

func TestMatrixTranspose(t *testing.T) {
	src := `var
   M : matriz [5, 3] = {{7, 12, 5},
                        {1, 4, 22},
                        {6, 20, 13},
                        ...
                       }
   T : matriz [*,*] numerico
inicio
   transponer (M, T)
   imprimir (alen(T), "\n")
   imprimir (T[1], "\n")
   imprimir (T[2], "\n")
   imprimir (T[3], "\n")
fin


sub transponer ( M : matriz [*,*] numerico
             ref R : matriz [*,*] numerico)
var
   cant_filas = alen (M)
   cant_cols  = alen (M [1])
   filas, cols : numerico
inicio
   dim (R, cant_cols, cant_filas)
   desde filas=1 hasta cant_filas {
      desde cols=1 hasta cant_cols {
         R [filas, cols] = M [cols, filas]
      }
   }
fin`
	out := run(t, src, "")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if lines[0] != "3" {
		t.Errorf("transpose: alen(T) want 3, got %q", lines[0])
	}
	if lines[1] != "7 1 6" {
		t.Errorf("transpose row 1: want %q, got %q", "7 1 6", lines[1])
	}
	if lines[2] != "12 4 20" {
		t.Errorf("transpose row 2: want %q, got %q", "12 4 20", lines[2])
	}
	if lines[3] != "5 22 13" {
		t.Errorf("transpose row 3: want %q, got %q", "5 22 13", lines[3])
	}
}
