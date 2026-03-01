package interpreter

import (
	"fmt"
	"math"

	"github.com/matiasinsaurralde/sl/pkg/ast"
	"github.com/matiasinsaurralde/sl/pkg/lexer"
)

// controlFlow signals for non-local control (return/break).
type controlFlow int

const (
	cfNone   controlFlow = iota
	cfReturn             // retorna
	cfBreak              // salir
)

// signal carries control-flow information out of stmt/block evaluation.
type signal struct {
	kind  controlFlow
	value *Value // for retorna
}

// RuntimeError is an interpreter runtime error.
type RuntimeError struct {
	Line int
	Msg  string
}

func (e *RuntimeError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("runtime error at line %d: %s", e.Line, e.Msg)
	}
	return "runtime error: " + e.Msg
}

// Interpreter runs SL programs.
type Interpreter struct {
	global  *Env
	subs    map[string]*ast.SubDecl
	stdlib  map[string]*builtinEntry
	types   map[string]ast.TypeNode // named type aliases from tipos sections
	stdin   *StdinReader
	stdout  *StdoutWriter
	cmdArgs []string
	ofs     string // output field separator
	ifs     string // input field separator
}

// BuiltinFn is the signature for built-in functions.
// args are the evaluated argument values; byRef[i] is the lvalue pointer for ref args.
type BuiltinFn func(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error)

// builtinEntry holds a builtin fn and its ref-parameter positions.
type builtinEntry struct {
	fn      BuiltinFn
	refMask []int // argument indices that are byRef
}

// New creates an Interpreter.
func New(cmdArgs []string) *Interpreter {
	interp := &Interpreter{
		global:  NewEnv(),
		subs:    make(map[string]*ast.SubDecl),
		stdlib:  make(map[string]*builtinEntry),
		types:   make(map[string]ast.TypeNode),
		cmdArgs: cmdArgs,
		ofs:     " ",
		ifs:     ",",
	}
	interp.stdin = newStdinReader(interp)
	interp.stdout = newStdoutWriter()
	RegisterBuiltins(interp)
	return interp
}

// Run executes a parsed program.
func (interp *Interpreter) Run(prog *ast.Program) error {
	// Register subroutines first (no forward declarations needed)
	for _, sub := range prog.Subs {
		interp.subs[sub.Name] = sub
	}

	// Evaluate global const/tipos/var declarations
	if err := interp.evalDecls(interp.global, prog.Consts, prog.Types, prog.Vars); err != nil {
		return err
	}

	// Execute main body
	sig, err := interp.execBlock(interp.global, prog.Body)
	if err != nil {
		return err
	}
	_ = sig // top-level return is fine, ignore
	return nil
}

// evalDecls evaluates const/tipos/var declarations into env.
func (interp *Interpreter) evalDecls(env *Env, consts []*ast.ConstDecl, types []*ast.TiposDecl, vars []*ast.VarDecl) error {
	for _, c := range consts {
		v, err := interp.evalExpr(env, c.Init)
		if err != nil {
			return err
		}
		env.Define(c.Name, v)
	}
	// Register tipo aliases so NamedType can be resolved in zeroForType.
	for _, td := range types {
		interp.types[td.Name] = td.Type
	}
	for _, v := range vars {
		if err := interp.evalVarDecl(env, v); err != nil {
			return err
		}
	}
	return nil
}

func (interp *Interpreter) evalVarDecl(env *Env, d *ast.VarDecl) error {
	if d.Init != nil {
		// Single name with initializer
		v, err := interp.evalExpr(env, d.Init)
		if err != nil {
			return err
		}
		// Apply fill if init is an ArrayLit with Fill=true and type has size
		if d.Type != nil {
			v = interp.applyTypeToValue(v, d.Type)
		}
		for _, name := range d.Names {
			env.Define(name, v.Copy())
		}
		return nil
	}
	// No initializer: create zero values based on type
	for _, name := range d.Names {
		var zero *Value
		if d.Type != nil {
			zero = interp.zeroForType(d.Type)
		} else {
			zero = zeroNil()
		}
		env.Define(name, zero)
	}
	return nil
}

// applyTypeToValue ensures v has the shape declared by typ (fills open dims).
func (interp *Interpreter) applyTypeToValue(v *Value, typ ast.TypeNode) *Value {
	switch t := typ.(type) {
	case *ast.VectorType:
		if t.Size > 0 && v.Kind == KindArr {
			// Fill to declared size
			v = fillArr(v, t.Size, interp.elemKind(t.ElemType))
		}
	case *ast.MatrixType:
		if v.Kind == KindArr && len(t.Dims) >= 2 && t.Dims[0] > 0 {
			v = fillMatrix(v, t.Dims, interp.elemKind(t.ElemType))
		}
	}
	return v
}

func (interp *Interpreter) elemKind(typ ast.TypeNode) Kind {
	switch t := typ.(type) {
	case *ast.SimpleType:
		switch t.Name {
		case "numerico":
			return KindNum
		case "cadena":
			return KindStr
		case "logico":
			return KindBool
		}
	}
	return KindNum
}

// fillArr ensures arr has exactly n elements, repeating the last one if shorter.
func fillArr(arr *Value, n int, elemKind Kind) *Value {
	if arr.Kind != KindArr || arr.Arr == nil {
		return MakeArr(n, elemKind)
	}
	cur := len(arr.Arr)
	if cur >= n {
		return arr
	}
	var last *Value
	if cur > 0 {
		last = arr.Arr[cur-1]
	} else {
		last = makeZeroByKind(elemKind)
	}
	for i := cur; i < n; i++ {
		arr.Arr = append(arr.Arr, last.Copy())
	}
	return arr
}

// fillMatrix fills a matrix (arr of arrs) to dims[0] rows, repeating the last row.
func fillMatrix(arr *Value, dims []int, elemKind Kind) *Value {
	rows := dims[0]
	if arr.Kind != KindArr || arr.Arr == nil {
		return MakeMatrix(dims, elemKind)
	}
	cur := len(arr.Arr)
	// Fill row count
	var lastRow *Value
	if cur > 0 {
		lastRow = arr.Arr[cur-1]
	} else {
		lastRow = MakeMatrix(dims[1:], elemKind)
	}
	for i := cur; i < rows; i++ {
		arr.Arr = append(arr.Arr, lastRow.Copy())
	}
	return arr
}

func (interp *Interpreter) zeroForType(typ ast.TypeNode) *Value {
	switch t := typ.(type) {
	case *ast.SimpleType:
		switch t.Name {
		case "numerico":
			return zeroNum()
		case "cadena":
			return zeroStr()
		case "logico":
			return zeroBool()
		}
	case *ast.VectorType:
		if t.Size > 0 {
			return MakeArr(t.Size, interp.elemKind(t.ElemType))
		}
		return zeroNil() // open array, uninitialized
	case *ast.MatrixType:
		// Check if any dim is fixed
		hasFixed := false
		for _, d := range t.Dims {
			if d > 0 {
				hasFixed = true
				break
			}
		}
		if hasFixed {
			return MakeMatrix(t.Dims, interp.elemKind(t.ElemType))
		}
		return zeroNil() // open matrix
	case *ast.RegistroType:
		rec := &Value{Kind: KindRecord}
		for _, fd := range t.Fields {
			for _, name := range fd.Names {
				rec.FNames = append(rec.FNames, name)
				rec.Fields = append(rec.Fields, interp.zeroForType(fd.Type))
			}
		}
		return rec
	case *ast.NamedType:
		if resolved, ok := interp.types[t.Name]; ok {
			return interp.zeroForType(resolved)
		}
		return zeroNil()
	}
	return zeroNil()
}

// execBlock executes a list of statements in env.
func (interp *Interpreter) execBlock(env *Env, stmts []ast.Stmt) (*signal, error) {
	for _, stmt := range stmts {
		sig, err := interp.execStmt(env, stmt)
		if err != nil {
			return nil, err
		}
		if sig != nil && sig.kind != cfNone {
			return sig, nil
		}
	}
	return nil, nil
}

// execStmt executes a single statement.
func (interp *Interpreter) execStmt(env *Env, stmt ast.Stmt) (*signal, error) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		return nil, interp.execAssign(env, s)
	case *ast.CallStmt:
		_, err := interp.evalCall(env, s.Call)
		return nil, err
	case *ast.ImprimirStmt:
		return nil, interp.execImprimir(env, s)
	case *ast.LeerStmt:
		return nil, interp.execLeer(env, s)
	case *ast.SiStmt:
		return interp.execSi(env, s)
	case *ast.DesdeStmt:
		return interp.execDesde(env, s)
	case *ast.MientrasStmt:
		return interp.execMientras(env, s)
	case *ast.RepetirStmt:
		return interp.execRepetir(env, s)
	case *ast.EvalStmt:
		return interp.execEval(env, s)
	case *ast.SalirStmt:
		return &signal{kind: cfBreak}, nil
	case *ast.RetornaStmt:
		return interp.execRetorna(env, s)
	case *ast.TerminarStmt:
		return interp.execTerminar(env, s)
	case nil:
		return nil, nil
	}
	return nil, fmt.Errorf("unknown statement type %T", stmt)
}

func (interp *Interpreter) execAssign(env *Env, s *ast.AssignStmt) error {
	rhs, err := interp.evalExpr(env, s.Value)
	if err != nil {
		return err
	}
	return interp.assignLValue(env, s.Target, rhs)
}

// assignLValue assigns rhs to the lvalue described by expr.
func (interp *Interpreter) assignLValue(env *Env, expr ast.Expr, rhs *Value) error {
	switch e := expr.(type) {
	case *ast.IdentExpr:
		ptr, ok := env.GetPtr(e.Name)
		if !ok {
			// Auto-define
			env.Define(e.Name, rhs.Copy())
			return nil
		}
		ptr.Assign(rhs)
		return nil

	case *ast.IndexExpr:
		arrVal, err := interp.evalExpr(env, e.Array)
		if err != nil {
			return err
		}
		if len(e.Indices) == 1 {
			// Check if it's a string character assignment
			if arrVal.Kind == KindStr {
				idxVal, err := interp.evalExpr(env, e.Indices[0])
				if err != nil {
					return err
				}
				idx := int(idxVal.ToNum())
				if idx >= 1 && idx <= len([]rune(arrVal.Str)) {
					rs := []rune(arrVal.Str)
					// Get the rhs as single char
					var ch rune
					if rhs.Kind == KindStr && len([]rune(rhs.Str)) > 0 {
						ch = []rune(rhs.Str)[0]
					}
					rs[idx-1] = ch
					// We need to mutate the parent variable
					// arrVal is a copy; we need the pointer
					return interp.assignStrChar(env, e.Array, idx, ch)
				}
				return nil
			}
			idxVal, err := interp.evalExpr(env, e.Indices[0])
			if err != nil {
				return err
			}
			idx := int(idxVal.ToNum())
			if idx < 1 || idx > len(arrVal.Arr) {
				return fmt.Errorf("array index %d out of bounds [1..%d]", idx, len(arrVal.Arr))
			}
			arrVal.Arr[idx-1].Assign(rhs)
			return nil
		}
		// Multi-dimensional
		ptr, err := interp.resolveMultiIndex(arrVal, e.Indices, env)
		if err != nil {
			return err
		}
		if ptr == nil {
			return nil // OOB outer dimension: silently ignore
		}
		ptr.Assign(rhs)
		return nil

	case *ast.FieldExpr:
		recVal, err := interp.evalExpr(env, e.Record)
		if err != nil {
			return err
		}
		if recVal.Kind != KindRecord {
			return fmt.Errorf("field access on non-record type")
		}
		for i, name := range recVal.FNames {
			if name == e.Field {
				recVal.Fields[i].Assign(rhs)
				return nil
			}
		}
		return fmt.Errorf("field %q not found", e.Field)
	}
	return fmt.Errorf("not an lvalue: %T", expr)
}

// assignStrChar assigns ch at index idx (1-based) in the string variable pointed to by expr.
func (interp *Interpreter) assignStrChar(env *Env, expr ast.Expr, idx int, ch rune) error {
	ptr, err := interp.evalLValuePtr(env, expr)
	if err != nil {
		return err
	}
	if ptr == nil || ptr.Kind != KindStr {
		return nil
	}
	rs := []rune(ptr.Str)
	if idx < 1 || idx > len(rs) {
		return nil // silently ignore out-of-bounds
	}
	rs[idx-1] = ch
	ptr.Str = string(rs)
	return nil
}

// evalLValuePtr returns a *Value pointer for the given lvalue expression.
func (interp *Interpreter) evalLValuePtr(env *Env, expr ast.Expr) (*Value, error) {
	switch e := expr.(type) {
	case *ast.IdentExpr:
		ptr, ok := env.GetPtr(e.Name)
		if !ok {
			return nil, fmt.Errorf("undefined variable %q", e.Name)
		}
		return ptr, nil
	case *ast.IndexExpr:
		arrVal, err := interp.evalExpr(env, e.Array)
		if err != nil {
			return nil, err
		}
		if len(e.Indices) == 1 {
			idxVal, err := interp.evalExpr(env, e.Indices[0])
			if err != nil {
				return nil, err
			}
			idx := int(idxVal.ToNum())
			if arrVal.Kind == KindArr && idx >= 1 && idx <= len(arrVal.Arr) {
				return arrVal.Arr[idx-1], nil
			}
			return nil, nil
		}
		return interp.resolveMultiIndex(arrVal, e.Indices, env)
	case *ast.FieldExpr:
		recVal, err := interp.evalExpr(env, e.Record)
		if err != nil {
			return nil, err
		}
		if recVal.Kind == KindRecord {
			for i, name := range recVal.FNames {
				if name == e.Field {
					return recVal.Fields[i], nil
				}
			}
		}
		return nil, fmt.Errorf("field %q not found", e.Field)
	}
	return nil, fmt.Errorf("not an lvalue: %T", expr)
}

func (interp *Interpreter) resolveMultiIndex(arr *Value, indices []ast.Expr, env *Env) (*Value, error) {
	cur := arr
	last := len(indices) - 1
	for i, idxExpr := range indices {
		if cur == nil || cur.Kind != KindArr {
			return nil, fmt.Errorf("index into non-array at dimension %d", i+1)
		}
		idxVal, err := interp.evalExpr(env, idxExpr)
		if err != nil {
			return nil, err
		}
		idx := int(idxVal.ToNum())
		if idx < 1 || idx > len(cur.Arr) {
			if i == last {
				// Auto-extend the innermost dimension only
				for len(cur.Arr) < idx {
					cur.Arr = append(cur.Arr, zeroNil())
				}
			} else {
				// Outer dimension OOB: silently skip (return nil, no error)
				return nil, nil
			}
		}
		cur = cur.Arr[idx-1]
	}
	return cur, nil
}

func (interp *Interpreter) execImprimir(env *Env, s *ast.ImprimirStmt) error {
	var sb []byte
	for _, arg := range s.Args {
		v, err := interp.evalExpr(env, arg)
		if err != nil {
			return err
		}
		sb = append(sb, v.ToStr(interp.ofs)...)
	}
	interp.stdout.Write(sb)
	return nil
}

func (interp *Interpreter) execLeer(env *Env, s *ast.LeerStmt) error {
	for _, varExpr := range s.Vars {
		if err := interp.leerOne(env, varExpr); err != nil {
			return err
		}
	}
	return nil
}

func (interp *Interpreter) leerOne(env *Env, expr ast.Expr) error {
	ptr, err := interp.evalLValuePtr(env, expr)
	if err != nil {
		// Try to evaluate as ident
		if id, ok := expr.(*ast.IdentExpr); ok {
			v, exists := env.Get(id.Name)
			if exists && v.Kind == KindArr {
				// Read entire array
				return interp.leerArr(v)
			}
		}
		return err
	}
	if ptr != nil && ptr.Kind == KindArr {
		return interp.leerArr(ptr)
	}
	// Read a single token
	tok := interp.stdin.NextToken(interp.ifs)
	if ptr == nil {
		return nil
	}
	switch ptr.Kind {
	case KindNum, KindNil:
		ptr.Kind = KindNum
		ptr.Num = parseNum(tok)
	case KindStr:
		ptr.Str = tok
	case KindBool:
		ptr.Bool = tok == "TRUE" || tok == "SI" || tok == "true" || tok == "1"
	default:
		ptr.Kind = KindNum
		ptr.Num = parseNum(tok)
	}
	return nil
}

func (interp *Interpreter) leerArr(arr *Value) error {
	for _, elem := range arr.Arr {
		if elem.Kind == KindArr {
			if err := interp.leerArr(elem); err != nil {
				return err
			}
		} else {
			tok := interp.stdin.NextToken(interp.ifs)
			switch elem.Kind {
			case KindNum, KindNil:
				elem.Kind = KindNum
				elem.Num = parseNum(tok)
			case KindStr:
				elem.Str = tok
			}
		}
	}
	return nil
}

func parseNum(s string) float64 {
	var f float64
	_, _ = fmt.Sscanf(s, "%g", &f)
	return f
}

func (interp *Interpreter) execSi(env *Env, s *ast.SiStmt) (*signal, error) {
	cond, err := interp.evalExpr(env, s.Cond)
	if err != nil {
		return nil, err
	}
	if cond.Truthy() {
		return interp.execBlock(env, s.Then)
	}
	for _, ei := range s.ElseIfs {
		econd, err := interp.evalExpr(env, ei.Cond)
		if err != nil {
			return nil, err
		}
		if econd.Truthy() {
			return interp.execBlock(env, ei.Body)
		}
	}
	if s.Else != nil {
		return interp.execBlock(env, s.Else)
	}
	return nil, nil
}

func (interp *Interpreter) execDesde(env *Env, s *ast.DesdeStmt) (*signal, error) {
	startVal, err := interp.evalExpr(env, s.Start)
	if err != nil {
		return nil, err
	}
	endVal, err := interp.evalExpr(env, s.End)
	if err != nil {
		return nil, err
	}
	var stepVal float64 = 1
	if s.Step != nil {
		sv, err := interp.evalExpr(env, s.Step)
		if err != nil {
			return nil, err
		}
		stepVal = math.Trunc(sv.ToNum())
	}

	start := math.Trunc(startVal.ToNum())
	end := endVal.ToNum()

	if stepVal == 0 {
		return nil, fmt.Errorf("desde loop: paso cannot be zero")
	}

	// Pre-compute count
	var count int
	if stepVal > 0 {
		count = int(math.Floor((end-start)/stepVal)) + 1
	} else {
		count = int(math.Floor((start-end)/(-stepVal))) + 1
	}

	// Set the loop variable
	ptr, ok := env.GetPtr(s.Var)
	if !ok {
		// Define if not exists
		env.Define(s.Var, zeroNum())
		ptr, _ = env.GetPtr(s.Var)
	}

	for i := 0; i < count; i++ {
		ptr.Kind = KindNum
		ptr.Num = start + float64(i)*stepVal
		sig, err := interp.execBlock(env, s.Body)
		if err != nil {
			return nil, err
		}
		if sig != nil {
			if sig.kind == cfBreak {
				return nil, nil
			}
			return sig, nil
		}
	}
	return nil, nil
}

func (interp *Interpreter) execMientras(env *Env, s *ast.MientrasStmt) (*signal, error) {
	for {
		cond, err := interp.evalExpr(env, s.Cond)
		if err != nil {
			return nil, err
		}
		if !cond.Truthy() {
			break
		}
		sig, err := interp.execBlock(env, s.Body)
		if err != nil {
			return nil, err
		}
		if sig != nil {
			if sig.kind == cfBreak {
				return nil, nil
			}
			return sig, nil
		}
	}
	return nil, nil
}

func (interp *Interpreter) execRepetir(env *Env, s *ast.RepetirStmt) (*signal, error) {
	for {
		sig, err := interp.execBlock(env, s.Body)
		if err != nil {
			return nil, err
		}
		if sig != nil {
			if sig.kind == cfBreak {
				return nil, nil
			}
			return sig, nil
		}
		cond, err := interp.evalExpr(env, s.Cond)
		if err != nil {
			return nil, err
		}
		if cond.Truthy() {
			break
		}
	}
	return nil, nil
}

func (interp *Interpreter) execEval(env *Env, s *ast.EvalStmt) (*signal, error) {
	for _, c := range s.Cases {
		cond, err := interp.evalExpr(env, c.Cond)
		if err != nil {
			return nil, err
		}
		if cond.Truthy() {
			return interp.execBlock(env, c.Body)
		}
	}
	if s.Else != nil {
		return interp.execBlock(env, s.Else)
	}
	return nil, nil
}

func (interp *Interpreter) execRetorna(env *Env, s *ast.RetornaStmt) (*signal, error) {
	var v *Value
	if s.Value != nil {
		var err error
		v, err = interp.evalExpr(env, s.Value)
		if err != nil {
			return nil, err
		}
	}
	return &signal{kind: cfReturn, value: v}, nil
}

func (interp *Interpreter) execTerminar(env *Env, s *ast.TerminarStmt) (*signal, error) {
	msg := ""
	if s.Msg != nil {
		v, err := interp.evalExpr(env, s.Msg)
		if err != nil {
			return nil, err
		}
		msg = v.ToStr(interp.ofs)
	}
	if msg != "" {
		interp.stdout.WriteString(msg + "\n")
	}
	return nil, &TerminateError{Msg: msg}
}

// TerminateError is returned when the program calls terminar().
type TerminateError struct {
	Msg string
}

func (e *TerminateError) Error() string {
	return "program terminated: " + e.Msg
}

// ---- Expression evaluation ----

func (interp *Interpreter) evalExpr(env *Env, expr ast.Expr) (*Value, error) {
	switch e := expr.(type) {
	case *ast.NumberLit:
		return &Value{Kind: KindNum, Num: e.Value}, nil
	case *ast.StringLit:
		return &Value{Kind: KindStr, Str: e.Value}, nil
	case *ast.BoolLit:
		return &Value{Kind: KindBool, Bool: e.Value}, nil
	case *ast.IdentExpr:
		return interp.evalIdent(env, e)
	case *ast.BinaryExpr:
		return interp.evalBinary(env, e)
	case *ast.UnaryExpr:
		return interp.evalUnary(env, e)
	case *ast.CallExpr:
		return interp.evalCall(env, e)
	case *ast.IndexExpr:
		return interp.evalIndex(env, e)
	case *ast.FieldExpr:
		return interp.evalField(env, e)
	case *ast.ArrayLit:
		return interp.evalArrayLit(env, e)
	case nil:
		return zeroNil(), nil
	}
	return zeroNil(), fmt.Errorf("unknown expression type %T", expr)
}

func (interp *Interpreter) evalIdent(env *Env, e *ast.IdentExpr) (*Value, error) {
	v, ok := env.Get(e.Name)
	if !ok {
		// Check stdlib
		return nil, fmt.Errorf("line %d: undefined identifier %q", e.Line, e.Name)
	}
	return v, nil
}

func (interp *Interpreter) evalBinary(env *Env, e *ast.BinaryExpr) (*Value, error) {
	left, err := interp.evalExpr(env, e.Left)
	if err != nil {
		return nil, err
	}
	right, err := interp.evalExpr(env, e.Right)
	if err != nil {
		return nil, err
	}

	switch e.Op {
	// Arithmetic
	case lexer.PLUS:
		if left.Kind == KindStr || right.Kind == KindStr {
			return &Value{Kind: KindStr, Str: left.ToStr("") + right.ToStr("")}, nil
		}
		return &Value{Kind: KindNum, Num: left.ToNum() + right.ToNum()}, nil
	case lexer.MINUS:
		return &Value{Kind: KindNum, Num: left.ToNum() - right.ToNum()}, nil
	case lexer.STAR:
		return &Value{Kind: KindNum, Num: left.ToNum() * right.ToNum()}, nil
	case lexer.SLASH:
		r := right.ToNum()
		if r == 0 {
			return nil, &RuntimeError{Line: e.Line, Msg: "division by zero"}
		}
		return &Value{Kind: KindNum, Num: left.ToNum() / r}, nil
	case lexer.PERCENT:
		// Truncate both to int before modulo
		l := int64(math.Trunc(left.ToNum()))
		r := int64(math.Trunc(right.ToNum()))
		if r == 0 {
			return nil, &RuntimeError{Line: e.Line, Msg: "modulo by zero"}
		}
		return &Value{Kind: KindNum, Num: float64(l % r)}, nil
	case lexer.CARET:
		return &Value{Kind: KindNum, Num: math.Pow(left.ToNum(), right.ToNum())}, nil

	// Relational
	case lexer.EQ:
		return &Value{Kind: KindBool, Bool: valEqual(left, right)}, nil
	case lexer.NEQ:
		return &Value{Kind: KindBool, Bool: !valEqual(left, right)}, nil
	case lexer.LT:
		return &Value{Kind: KindBool, Bool: valLess(left, right)}, nil
	case lexer.LE:
		return &Value{Kind: KindBool, Bool: valLess(left, right) || valEqual(left, right)}, nil
	case lexer.GT:
		return &Value{Kind: KindBool, Bool: valLess(right, left)}, nil
	case lexer.GE:
		return &Value{Kind: KindBool, Bool: valLess(right, left) || valEqual(left, right)}, nil

	// Logical
	case lexer.AND:
		return &Value{Kind: KindBool, Bool: left.Truthy() && right.Truthy()}, nil
	case lexer.OR:
		return &Value{Kind: KindBool, Bool: left.Truthy() || right.Truthy()}, nil
	}
	return zeroNil(), fmt.Errorf("unknown binary operator %s", e.Op)
}

func valEqual(a, b *Value) bool {
	if a.Kind == KindNum && b.Kind == KindNum {
		return a.Num == b.Num
	}
	if a.Kind == KindStr && b.Kind == KindStr {
		return a.Str == b.Str
	}
	if a.Kind == KindBool && b.Kind == KindBool {
		return a.Bool == b.Bool
	}
	if a.Kind == KindNum && b.Kind == KindBool {
		return (a.Num != 0) == b.Bool
	}
	if a.Kind == KindBool && b.Kind == KindNum {
		return a.Bool == (b.Num != 0)
	}
	return false
}

func valLess(a, b *Value) bool {
	if a.Kind == KindNum && b.Kind == KindNum {
		return a.Num < b.Num
	}
	if a.Kind == KindStr && b.Kind == KindStr {
		return a.Str < b.Str
	}
	if a.Kind == KindNum {
		return a.Num < b.ToNum()
	}
	if b.Kind == KindNum {
		return a.ToNum() < b.Num
	}
	return false
}

func (interp *Interpreter) evalUnary(env *Env, e *ast.UnaryExpr) (*Value, error) {
	operand, err := interp.evalExpr(env, e.Operand)
	if err != nil {
		return nil, err
	}
	switch e.Op {
	case lexer.MINUS:
		return &Value{Kind: KindNum, Num: -operand.ToNum()}, nil
	case lexer.PLUS:
		return &Value{Kind: KindNum, Num: operand.ToNum()}, nil
	case lexer.NOT:
		return &Value{Kind: KindBool, Bool: !operand.Truthy()}, nil
	}
	return zeroNil(), nil
}

func (interp *Interpreter) evalCall(env *Env, e *ast.CallExpr) (*Value, error) {
	// Check user-defined subs first
	if sub, ok := interp.subs[e.Name]; ok {
		return interp.callSub(env, sub, e.Args, e.Line)
	}
	// Check stdlib
	if entry, ok := interp.stdlib[e.Name]; ok {
		return interp.callBuiltin(env, entry, e.Args)
	}
	return nil, fmt.Errorf("line %d: undefined function %q", e.Line, e.Name)
}

// callBuiltin evaluates args and calls the builtin fn.
// entry.refMask lists which argument indices should be passed by reference.
func (interp *Interpreter) callBuiltin(env *Env, entry *builtinEntry, argExprs []ast.Expr) (*Value, error) {
	args := make([]*Value, len(argExprs))
	byRef := make([]*Value, len(argExprs))

	// Build a set of ref positions from the entry's refMask
	refSet := make(map[int]bool, len(entry.refMask))
	for _, idx := range entry.refMask {
		refSet[idx] = true
	}

	for i, ae := range argExprs {
		if refSet[i] {
			ptr, err := interp.evalLValuePtr(env, ae)
			if err == nil && ptr != nil {
				byRef[i] = ptr
				args[i] = ptr
				continue
			}
		}
		v, err := interp.evalExpr(env, ae)
		if err != nil {
			return nil, err
		}
		args[i] = v
	}
	return entry.fn(interp, args, byRef)
}

// callSub calls a user-defined subroutine.
func (interp *Interpreter) callSub(callerEnv *Env, sub *ast.SubDecl, argExprs []ast.Expr, line int) (*Value, error) {
	// Create new environment for the call
	callEnv := interp.global.Child()

	// Bind parameters
	paramIdx := 0
	for _, g := range sub.Params {
		for _, name := range g.Names {
			if paramIdx >= len(argExprs) {
				callEnv.Define(name, makeZeroByKind(KindNum))
				paramIdx++
				continue
			}
			ae := argExprs[paramIdx]
			if g.ByRef {
				// Pass by reference: get the lvalue pointer
				ptr, err := interp.evalLValuePtr(callerEnv, ae)
				if err != nil || ptr == nil {
					// Fall back to value
					v, err2 := interp.evalExpr(callerEnv, ae)
					if err2 != nil {
						return nil, err2
					}
					callEnv.Define(name, v)
				} else {
					// Share the same pointer
					callEnv.vars[name] = ptr
				}
			} else {
				v, err := interp.evalExpr(callerEnv, ae)
				if err != nil {
					return nil, err
				}
				callEnv.Define(name, v.Copy())
			}
			paramIdx++
		}
	}

	// Local declarations (evaluated in the call env, so params are accessible)
	if err := interp.evalDecls(callEnv, sub.Consts, sub.Types, sub.Vars); err != nil {
		return nil, err
	}

	// Execute body
	sig, err := interp.execBlock(callEnv, sub.Body)
	if err != nil {
		return nil, err
	}
	if sig != nil && sig.kind == cfReturn && sig.value != nil {
		return sig.value, nil
	}
	return zeroNil(), nil
}

func (interp *Interpreter) evalIndex(env *Env, e *ast.IndexExpr) (*Value, error) {
	arrVal, err := interp.evalExpr(env, e.Array)
	if err != nil {
		return nil, err
	}

	if len(e.Indices) == 1 {
		idxVal, err := interp.evalExpr(env, e.Indices[0])
		if err != nil {
			return nil, err
		}
		idx := int(idxVal.ToNum())

		// String character indexing
		if arrVal.Kind == KindStr {
			rs := []rune(arrVal.Str)
			if idx < 1 || idx > len(rs) {
				return &Value{Kind: KindStr, Str: ""}, nil
			}
			return &Value{Kind: KindStr, Str: string(rs[idx-1])}, nil
		}

		// Array indexing
		if arrVal.Kind == KindArr {
			if idx < 1 || idx > len(arrVal.Arr) {
				return zeroNil(), nil
			}
			return arrVal.Arr[idx-1], nil
		}
		return zeroNil(), nil
	}

	// Multi-dimensional
	cur := arrVal
	for i, idxExpr := range e.Indices {
		if cur == nil || cur.Kind != KindArr {
			return zeroNil(), fmt.Errorf("index into non-array at dimension %d", i+1)
		}
		idxVal, err := interp.evalExpr(env, idxExpr)
		if err != nil {
			return nil, err
		}
		idx := int(idxVal.ToNum())
		if idx < 1 || idx > len(cur.Arr) {
			return zeroNil(), nil
		}
		cur = cur.Arr[idx-1]
	}
	return cur, nil
}

func (interp *Interpreter) evalField(env *Env, e *ast.FieldExpr) (*Value, error) {
	rec, err := interp.evalExpr(env, e.Record)
	if err != nil {
		return nil, err
	}
	if rec.Kind != KindRecord {
		return zeroNil(), nil
	}
	for i, name := range rec.FNames {
		if name == e.Field {
			return rec.Fields[i], nil
		}
	}
	return zeroNil(), fmt.Errorf("field %q not found in record", e.Field)
}

func (interp *Interpreter) evalArrayLit(env *Env, e *ast.ArrayLit) (*Value, error) {
	if len(e.Elems) == 0 && !e.Fill {
		// {} = clear/empty array
		return zeroNil(), nil
	}

	arr := make([]*Value, len(e.Elems))
	for i, elem := range e.Elems {
		v, err := interp.evalExpr(env, elem)
		if err != nil {
			return nil, err
		}
		arr[i] = v
	}

	v := &Value{Kind: KindArr, Arr: arr}
	// Fill flag is applied when the variable has a known size (handled in applyTypeToValue).
	// Store the fill flag in a special way: we re-use the concept at assignment time.
	// For now just return the partial array with Fill semantics handled by context.
	_ = e.Fill // fill is applied by applyTypeToValue or at assignment
	return v, nil
}

// RegisterBuiltin adds a built-in function with optional ref-parameter positions.
func (interp *Interpreter) RegisterBuiltin(name string, fn BuiltinFn, refPositions ...int) {
	interp.stdlib[name] = &builtinEntry{fn: fn, refMask: refPositions}
}
