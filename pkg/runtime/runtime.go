package runtime

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/matiasinsaurralde/sl/pkg/ast"
)

type Runtime struct {
	Scope  *ast.Scope
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func NewRuntime() *Runtime {
	return &Runtime{
		Scope:  ast.NewScope(nil),
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func NewRuntimeWithIO(stdin io.Reader, stdout io.Writer, stderr io.Writer) *Runtime {
	return &Runtime{
		Scope:  ast.NewScope(nil),
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (rt *Runtime) RunFile(filename string, fileAst *ast.File) {
	for _, node := range fileAst.Nodes {
		rt.evalNode(node)
	}
}

func (rt *Runtime) evalNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.BlockStatement:
		rt.evalBlock(n)
	case *ast.VariableDeclaration:
		var val interface{} = 0 // Default to 0 for numerico variables
		if n.Value != nil {
			val = rt.evalExpr(n.Value)
		}
		rt.Scope.Set(n.Name, val)
	case *ast.ConstantDeclaration:
		val := rt.evalExpr(n.Value)
		rt.Scope.Set(n.Name, val)
	case *ast.ExpressionStatement:
		rt.evalExpr(n.Expression)
	case *ast.ForStatement:
		rt.evalForStatement(n)
	case *ast.WhileStatement:
		rt.evalWhileStatement(n)
	case *ast.RepeatStatement:
		rt.evalRepeatStatement(n)
	case *ast.IfStatement:
		rt.evalIfStatement(n)
	case *ast.TerminateStatement:
		rt.evalTerminateStatement(n)
	}
}

func (rt *Runtime) evalBlock(block *ast.BlockStatement) {
	for _, stmt := range block.Statements {
		rt.evalNode(stmt)
	}
}

func (rt *Runtime) evalExpr(expr ast.Expression) interface{} {
	switch e := expr.(type) {
	case *ast.Literal:
		if e.Type.String() == "INT" {
			var v int
			fmt.Sscanf(e.Value, "%d", &v)
			return v
		}
		if e.Type.String() == "STRING" {
			// Unquote string literals to remove quotes and interpret escape sequences
			if len(e.Value) >= 2 && e.Value[0] == '"' && e.Value[len(e.Value)-1] == '"' {
				if unquoted, err := strconv.Unquote(e.Value); err == nil {
					return unquoted
				}
			}
		}
		return e.Value
	case *ast.Identifier:
		val, _ := rt.Scope.Get(e.Name)
		return val
	case *ast.AssignmentExpression:
		val := rt.evalExpr(e.Right)
		rt.Scope.Set(e.Left.Name, val)
		return val
	case *ast.CallExpression:
		return rt.evalCall(e)
	case *ast.IndexExpression:
		return rt.evalIndexExpression(e)
	case *ast.BinaryExpression:
		// Handle assignment if operator is '='
		if e.Operator == "=" {
			if ident, ok := e.Left.(*ast.Identifier); ok {
				// Evaluate right side
				rightVal := rt.evalExpr(e.Right)

				// Store the result
				rt.Scope.Set(ident.Name, rightVal)
				return rightVal
			}
		}
		left := rt.evalExpr(e.Left)
		right := rt.evalExpr(e.Right)
		return evalBinary(left, right, e.Operator)
	default:
		return nil
	}
}

func (rt *Runtime) evalCall(call *ast.CallExpression) interface{} {
	switch call.Function {
	case "imprimir":
		for _, arg := range call.Arguments {
			val := rt.evalExpr(arg)
			// Handle string escape sequences
			if str, ok := val.(string); ok {
				// Strip quotes if present and unquote to interpret escape sequences
				if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
					if unquoted, err := strconv.Unquote(str); err == nil {
						fmt.Fprint(rt.stdout, unquoted)
					} else {
						fmt.Fprint(rt.stdout, str)
					}
				} else {
					fmt.Fprint(rt.stdout, str)
				}
			} else {
				fmt.Fprint(rt.stdout, val)
			}
		}
		fmt.Fprintln(rt.stdout)
	case "leer":
		reader := bufio.NewReader(rt.stdin)
		for _, arg := range call.Arguments {
			if id, ok := arg.(*ast.Identifier); ok {
				fmt.Fprintf(rt.stdout, "Ingrese valor para %s: ", id.Name)
				text, _ := reader.ReadString('\n')
				var v int
				fmt.Sscanf(text, "%d", &v)
				rt.Scope.Set(id.Name, v)
			}
		}
	case "int":
		// Convert to integer
		if len(call.Arguments) > 0 {
			val := rt.evalExpr(call.Arguments[0])
			switch v := val.(type) {
			case float64:
				return int(v)
			case string:
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					return int(f)
				}
			}
		}
	}
	return nil
}

func (rt *Runtime) evalIndexExpression(indexExpr *ast.IndexExpression) interface{} {
	left := rt.evalExpr(indexExpr.Left)
	index := rt.evalExpr(indexExpr.Index)

	// Handle string indexing
	if str, ok := left.(string); ok {
		if idx, ok := index.(int); ok {
			if idx >= 0 && idx < len(str) {
				return string(str[idx])
			}
		}
	}

	return nil
}

func evalBinary(left, right interface{}, op string) interface{} {
	// Handle nil values
	if left == nil {
		left = 0
	}
	if right == nil {
		right = 0
	}

	// Handle logical operators first (they work with any types)
	switch op {
	case "or":
		return isTrue(left) || isTrue(right)
	case "and":
		return isTrue(left) && isTrue(right)
	}

	// Handle string concatenation
	if op == "+" {
		_, leftIsStr := left.(string)
		_, rightIsStr := right.(string)
		if leftIsStr || rightIsStr {
			return fmt.Sprintf("%v%v", left, right)
		}
	}

	li, lok := left.(int)
	ri, rok := right.(int)
	if lok && rok {
		switch op {
		case "+":
			return li + ri
		case "-":
			return li - ri
		case "*":
			return li * ri
		case "/":
			if ri != 0 {
				return li / ri
			}
		case "%":
			if ri != 0 {
				return li % ri
			}
		case ">=":
			return li >= ri
		case "<=":
			return li <= ri
		case ">":
			return li > ri
		case "<":
			return li < ri
		case "==":
			return li == ri
		case "<>":
			return li != ri
		}
	}
	// fallback: string concat
	return fmt.Sprintf("%v%v", left, right)
}

func (rt *Runtime) evalForStatement(forStmt *ast.ForStatement) {
	start := rt.evalExpr(forStmt.Start).(int)
	end := rt.evalExpr(forStmt.EndExpr).(int)
	step := 1
	if forStmt.Step != nil {
		step = rt.evalExpr(forStmt.Step).(int)
	}

	for i := start; i <= end; i += step {
		rt.Scope.Set(forStmt.Variable, i)
		rt.evalNode(forStmt.Body)
	}
}

func (rt *Runtime) evalWhileStatement(whileStmt *ast.WhileStatement) {
	for {
		condition := rt.evalExpr(whileStmt.Condition)
		if !isTrue(condition) {
			break
		}
		rt.evalNode(whileStmt.Body)
	}
}

func (rt *Runtime) evalRepeatStatement(repeatStmt *ast.RepeatStatement) {
	for {
		rt.evalNode(repeatStmt.Body)
		condition := rt.evalExpr(repeatStmt.Condition)
		if isTrue(condition) {
			break
		}
	}
}

func (rt *Runtime) evalIfStatement(ifStmt *ast.IfStatement) {
	condition := rt.evalExpr(ifStmt.Condition)
	if isTrue(condition) {
		rt.evalNode(ifStmt.Then)
	} else if ifStmt.Else != nil {
		rt.evalNode(ifStmt.Else)
	}
}

func (rt *Runtime) evalTerminateStatement(terminateStmt *ast.TerminateStatement) {
	if terminateStmt.Message != nil {
		message := rt.evalExpr(terminateStmt.Message)
		if str, ok := message.(string); ok {
			// Handle string escape sequences
			if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
				if unquoted, err := strconv.Unquote(str); err == nil {
					fmt.Fprint(rt.stdout, unquoted)
				} else {
					fmt.Fprint(rt.stdout, str)
				}
			} else {
				fmt.Fprint(rt.stdout, str)
			}
		} else {
			fmt.Fprint(rt.stdout, message)
		}
	}
	// Exit the program
	os.Exit(0)
}

func isTrue(val interface{}) bool {
	switch v := val.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case string:
		return v != ""
	default:
		return val != nil
	}
}
