package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/matiasinsaurralde/sl/parser"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app       = kingpin.New("sl", "Un intérprete experimental de SL")
	debug     = app.Flag("debug", "Habilitar el modo de depuración").Bool()
	inputFile = app.Arg("programa", "Programa a ejecutar").Required().String()
	help      = app.HelpFlag.Short('h')
)

func readSource(input string) (string, error) {
	data, err := ioutil.ReadFile(input)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

var code = `programa holamundo
var a: numerico
imprimir(a)
fin
`

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	src, err := readSource(*inputFile)
	if err != nil {
		panic(err)
	}

	reader := strings.NewReader(src)
	parser, _ := parser.New(reader)
	parser.Parse()
}
