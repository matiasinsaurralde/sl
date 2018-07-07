package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	lexer "github.com/matiasinsaurralde/sl/lexer"
	logger "github.com/matiasinsaurralde/sl/log"
	parser "github.com/matiasinsaurralde/sl/parser"
	runtime "github.com/matiasinsaurralde/sl/runtime"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app       = kingpin.New("sl", "Un intérprete experimental de SL")
	debug     = app.Flag("debug", "Habilitar el modo de depuración").Bool()
	inputFile = app.Arg("programa", "Programa a ejecutar").Required().String()
	help      = app.HelpFlag.Short('h')

	log = logger.Logger
)

func readSource(input string) (string, error) {
	data, err := ioutil.ReadFile(input)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	if *debug {
		log.Level = logrus.DebugLevel
	} else {
		log.Level = logrus.InfoLevel
	}
	src, err := readSource(*inputFile)
	if err != nil {
		panic(err)
	}

	reader := strings.NewReader(src)
	lexer, _ := lexer.New(reader)
	tokenSet := lexer.Parse()
	reader.Reset(src)
	p := parser.New(reader, tokenSet)
	ast := p.Parse()
	runtime := runtime.New(ast)
	runtime.Init()
	runtime.Run()
}
