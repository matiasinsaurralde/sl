package runtime

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	slast "github.com/matiasinsaurralde/sl/ast"
	logger "github.com/matiasinsaurralde/sl/log"
)

var (
	log = logger.Logger
)

// type imprimirFn func(string)

type Runtime struct {
	ast *slast.AST

	globals map[string]interface{}
}

/*
func (r *Runtime) call(t *parser.ExprStmt) {
	v, found := r.globals[t.Name]
	if !found {
		panic("Undefined " + t.Name)
	}
	switch t.Name {
	case "imprimir":
		f := v.(func(string))
		f(t.Expr)
	case "leer":
		f := v.(func())
		f()
	}
}*/

type genericFn func(...interface{})

func (r *Runtime) Init() {
	// registerStdio(r)
	r.globals["imprimir"] = func(args ...interface{}) {
		if len(args) > 1 {
			panic("Too many args")
		}
		v := args[0]
		switch v.(type) {
		case int:
			v, _ := args[0].(int)
			fmt.Println(v)
		case float64:
			v, _ := args[0].(float64)
			fmt.Println(v)
		case string:
			v, _ := args[0].(string)
			fmt.Println(v)
		case nil:
			r.error("Valor es nulo al llamar 'imprimir'")
		}
	}
	r.globals["leer"] = func() {
	}
}

func New(ast *slast.AST) *Runtime {
	return &Runtime{
		ast:     ast,
		globals: make(map[string]interface{}),
	}
}

func (r *Runtime) resolveValue(scope *slast.Scope, n string) (interface{}, error) {
	localVal, ok := scope.Objects[n]
	if ok {
		log.WithField("prefix", "runtime").
			Debugf("Resolviendo valor local '%s'.", n)
		return localVal.GetValue(), nil
	}
	globalVal, ok := r.ast.GlobalScope.Objects[n]
	if ok {
		log.WithField("prefix", "runtime").
			Debugf("Resolviendo valor global '%s'.", n)
		return globalVal.GetValue(), nil
	}
	return nil, errors.New("")
}

func (r *Runtime) symLookup(name string) (interface{}, error) {
	log.WithField("prefix", "runtime").
		Debugf("Resolviendo simbolo '%s'.", name)
	sym, ok := r.globals[name]
	if !ok {
		return nil, errors.New("Simbolo no encontrado")
	}
	return sym, nil
}

func (r *Runtime) evaluate(scope *slast.Scope, n slast.Node) {
	switch n.(type) {
	case *slast.AssignStmt:
		s := n.(*slast.AssignStmt)
		varName := s.X.GetValue().(string)
		v, exists := scope.Objects[varName]
		if !exists {
			panic("La variable " + varName + "no existe")
		}
		// existingType := reflect.TypeOf(v.GetValue())
		var node slast.Node
		switch v.GetValue().(type) {
		case int:
			lit := s.Y.(*slast.BasicLit)
			y, _ := strconv.Atoi(lit.GetValue().(string))
			node = &slast.Decl{
				Name:  varName,
				Value: y,
			}
		case string:
			lit := s.Y.(*slast.BasicLit)
			y, _ := lit.GetValue().(string)
			node = &slast.Decl{
				Name:  varName,
				Value: y,
			}
		}
		scope.Objects[varName] = node
	case *slast.Stmt:
		s := n.(*slast.Stmt)
		key := s.Value
		val, err := r.resolveValue(scope, key)
		if err != nil {
			r.errorf("Variable '%s' no definida", key)
		}
		sym, err := r.symLookup(s.Name)
		if err != nil {
			r.errorf("Simbolo '%s' no encontrado", s.Name)
		}
		f := sym.(func(...interface{}))
		f(val)
	}
}

func (r *Runtime) Run() {
	log.WithField("prefix", "runtime").
		Debugf("%d sentencia(s) encontradas.", len(r.ast.GlobalScope.Nodes))

	log.WithField("prefix", "runtime").
		Debugf("Iniciando")

	for _, v := range r.ast.GlobalScope.Nodes {
		block, ok := v.(*slast.BlockStmt)
		if !ok {
			panic(1)
		}
		for _, v := range block.Nodes {
			r.evaluate(block.Scope, v)
		}
	}
}

func (r *Runtime) error(message string) {
	log.WithField("prefix", "runtime").
		Error(message)
	os.Exit(1)
}

func (r *Runtime) errorf(format string, args ...interface{}) {
	log.WithField("prefix", "runtime").
		Errorf(format, args...)
	os.Exit(1)
}
