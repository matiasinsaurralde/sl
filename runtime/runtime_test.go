package runtime

import (
	"strings"
	"testing"

	lexer "github.com/matiasinsaurralde/sl/lexer"
	parser "github.com/matiasinsaurralde/sl/parser"
)

const testProgram = `programa test
var a, b : numerico
var c, d : numerico = 16
var e : numerico
var f : numerico = 32

inicio
imprimir(a)
imprimir(b)
imprimir(c)
imprimir(e)
fin
`

func TestRuntime(t *testing.T) {
	reader := strings.NewReader(testProgram)
	lexer, _ := lexer.New(reader)
	tokenSet := lexer.Parse()
	reader.Reset(testProgram)
	// parser, _ := parser.New(reader, tokenSet)
	p := parser.New(reader, tokenSet)
	ast := p.Parse()

	runtime := New(ast)
	runtime.Run()
}
