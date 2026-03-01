package interpreter

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// RegisterBuiltins registers all SL standard library functions.
func RegisterBuiltins(interp *Interpreter) {
	// Math
	interp.RegisterBuiltin("abs", builtinAbs)
	interp.RegisterBuiltin("arctan", builtinArctan)
	interp.RegisterBuiltin("cos", builtinCos)
	interp.RegisterBuiltin("exp", builtinExp)
	interp.RegisterBuiltin("int", builtinInt)
	interp.RegisterBuiltin("log", builtinLog)
	interp.RegisterBuiltin("sin", builtinSin)
	interp.RegisterBuiltin("sqrt", builtinSqrt)
	interp.RegisterBuiltin("tan", builtinTan)

	// String
	interp.RegisterBuiltin("ascii", builtinAscii)
	interp.RegisterBuiltin("lower", builtinLower)
	interp.RegisterBuiltin("ord", builtinOrd)
	interp.RegisterBuiltin("pos", builtinPos)
	interp.RegisterBuiltin("strdup", builtinStrdup)
	interp.RegisterBuiltin("strlen", builtinStrlen)
	interp.RegisterBuiltin("substr", builtinSubstr)
	interp.RegisterBuiltin("upper", builtinUpper)

	// I/O
	interp.RegisterBuiltin("imprimir", builtinImprimir)
	// leer: all args are byRef; we use a variadic approach — handled in execLeer directly
	interp.RegisterBuiltin("beep", builtinBeep)
	interp.RegisterBuiltin("cls", builtinCls)
	interp.RegisterBuiltin("eof", builtinEof)
	interp.RegisterBuiltin("get_color", builtinNoop)
	interp.RegisterBuiltin("get_curpos", builtinNoop)
	interp.RegisterBuiltin("get_ifs", builtinGetIfs)
	interp.RegisterBuiltin("get_ofs", builtinGetOfs)
	interp.RegisterBuiltin("get_scrsize", builtinNoop)
	interp.RegisterBuiltin("readkey", builtinNoop)
	interp.RegisterBuiltin("set_color", builtinNoop)
	interp.RegisterBuiltin("set_curpos", builtinNoop)
	interp.RegisterBuiltin("set_ifs", builtinSetIfs)
	interp.RegisterBuiltin("set_ofs", builtinSetOfs)
	interp.RegisterBuiltin("set_stdin", builtinSetStdin)
	interp.RegisterBuiltin("set_stdout", builtinSetStdout)

	// Type conversion
	interp.RegisterBuiltin("str", builtinStr)
	interp.RegisterBuiltin("val", builtinVal)

	// Arrays
	interp.RegisterBuiltin("alen", builtinAlen)
	interp.RegisterBuiltin("dim", builtinDim, 0) // arg 0 is byRef

	// Other
	interp.RegisterBuiltin("dec", builtinDec, 0) // arg 0 is byRef
	interp.RegisterBuiltin("ifval", builtinIfval)
	interp.RegisterBuiltin("inc", builtinInc, 0)                      // arg 0 is byRef
	interp.RegisterBuiltin("intercambiar", builtinIntercambiar, 0, 1) // args 0,1 are byRef
	interp.RegisterBuiltin("swap", builtinIntercambiar, 0, 1)
	interp.RegisterBuiltin("max", builtinMax)
	interp.RegisterBuiltin("min", builtinMin)
	interp.RegisterBuiltin("paramval", builtinParamval)
	interp.RegisterBuiltin("pcount", builtinPcount)
	interp.RegisterBuiltin("random", builtinRandom)
	interp.RegisterBuiltin("runcmd", builtinRuncmd)
	interp.RegisterBuiltin("sec", builtinSec)
	interp.RegisterBuiltin("terminar", builtinTerminar)
}

var startTime = time.Now()

// ---- Math ----

func builtinAbs(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: math.Abs(args[0].ToNum())}, nil
}

func builtinArctan(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: math.Atan(args[0].ToNum())}, nil
}

func builtinCos(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: math.Cos(args[0].ToNum())}, nil
}

func builtinExp(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: math.Exp(args[0].ToNum())}, nil
}

func builtinInt(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: math.Trunc(args[0].ToNum())}, nil
}

func builtinLog(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: math.Log(args[0].ToNum())}, nil
}

func builtinSin(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: math.Sin(args[0].ToNum())}, nil
}

func builtinSqrt(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: math.Sqrt(args[0].ToNum())}, nil
}

func builtinTan(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: math.Tan(args[0].ToNum())}, nil
}

// ---- String ----

func builtinAscii(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return &Value{Kind: KindStr, Str: ""}, nil
	}
	n := int(args[0].ToNum())
	return &Value{Kind: KindStr, Str: string(rune(n))}, nil
}

func builtinLower(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return &Value{Kind: KindStr, Str: ""}, nil
	}
	return &Value{Kind: KindStr, Str: strings.ToLower(args[0].ToStr(""))}, nil
}

func builtinOrd(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	s := args[0].ToStr("")
	if len(s) == 0 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: float64([]rune(s)[0])}, nil
}

func builtinPos(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 2 {
		return zeroNum(), nil
	}
	haystack := args[0].ToStr("")
	needle := args[1].ToStr("")
	start := 1
	if len(args) >= 3 {
		start = int(args[2].ToNum())
	}
	if start < 1 {
		start = 1
	}
	rs := []rune(haystack)
	if start > len(rs) {
		return zeroNum(), nil
	}
	sub := string(rs[start-1:])
	idx := strings.Index(sub, needle)
	if idx < 0 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: float64(start + len([]rune(sub[:idx])))}, nil
}

func builtinStrdup(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 2 {
		return &Value{Kind: KindStr, Str: ""}, nil
	}
	s := args[0].ToStr("")
	n := int(args[1].ToNum())
	return &Value{Kind: KindStr, Str: strings.Repeat(s, n)}, nil
}

func builtinStrlen(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: float64(len([]rune(args[0].ToStr(""))))}, nil
}

func builtinSubstr(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 2 {
		return &Value{Kind: KindStr, Str: ""}, nil
	}
	s := []rune(args[0].ToStr(""))
	pos := int(args[1].ToNum())
	if pos < 1 {
		pos = 1
	}
	if pos > len(s) {
		return &Value{Kind: KindStr, Str: ""}, nil
	}
	start := pos - 1
	if len(args) >= 3 {
		n := int(args[2].ToNum())
		end := start + n
		if end > len(s) {
			end = len(s)
		}
		return &Value{Kind: KindStr, Str: string(s[start:end])}, nil
	}
	return &Value{Kind: KindStr, Str: string(s[start:])}, nil
}

func builtinUpper(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return &Value{Kind: KindStr, Str: ""}, nil
	}
	return &Value{Kind: KindStr, Str: strings.ToUpper(args[0].ToStr(""))}, nil
}

// ---- I/O (as builtins, also handled as statements) ----

func builtinImprimir(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	var sb strings.Builder
	for _, a := range args {
		sb.WriteString(a.ToStr(interp.ofs))
	}
	interp.stdout.WriteString(sb.String())
	return zeroNil(), nil
}

// leer is handled as a statement (LeerStmt) and does not need a builtin entry.

func builtinBeep(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	// No-op in CLI mode
	return zeroNil(), nil
}

func builtinCls(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	fmt.Print("\033[H\033[2J")
	return zeroNil(), nil
}

func builtinEof(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	eof := interp.stdin.EOF()
	return &Value{Kind: KindBool, Bool: eof}, nil
}

func builtinGetIfs(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	return &Value{Kind: KindStr, Str: interp.ifs}, nil
}

func builtinGetOfs(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	return &Value{Kind: KindStr, Str: interp.ofs}, nil
}

func builtinSetIfs(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) > 0 {
		interp.ifs = args[0].ToStr("")
	}
	return zeroNil(), nil
}

func builtinSetOfs(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) > 0 {
		interp.ofs = args[0].ToStr("")
	}
	return zeroNil(), nil
}

func builtinSetStdin(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return &Value{Kind: KindBool, Bool: false}, nil
	}
	path := args[0].ToStr("")
	err := interp.stdin.SetFile(path)
	// Reset IFS to default after set_stdin
	interp.ifs = ","
	return &Value{Kind: KindBool, Bool: err == nil}, nil
}

func builtinSetStdout(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return &Value{Kind: KindBool, Bool: false}, nil
	}
	path := args[0].ToStr("")
	mode := "wt"
	if len(args) >= 2 {
		mode = args[1].ToStr("")
	}
	err := interp.stdout.SetFile(path, mode)
	return &Value{Kind: KindBool, Bool: err == nil}, nil
}

// ---- Type conversion ----

func builtinStr(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return &Value{Kind: KindStr, Str: ""}, nil
	}
	n := args[0].ToNum()
	var s string
	if len(args) >= 3 {
		width := int(args[1].ToNum())
		decimals := int(args[2].ToNum())
		fill := ' '
		if len(args) >= 4 {
			fs := args[3].ToStr("")
			if len(fs) > 0 {
				fill = rune(fs[0])
			}
		}
		format := fmt.Sprintf("%%%c%d.%df", fill, width, decimals)
		s = fmt.Sprintf(format, n)
	} else if len(args) >= 2 {
		width := int(args[1].ToNum())
		s = fmt.Sprintf("%*g", width, n)
	} else {
		s = formatNum(n)
	}
	return &Value{Kind: KindStr, Str: s}, nil
}

func builtinVal(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: parseNum(args[0].ToStr(""))}, nil
}

// ---- Arrays ----

func builtinAlen(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	v := args[0]
	if v.Kind != KindArr || v.Arr == nil {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: float64(len(v.Arr))}, nil
}

func builtinDim(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	// dim(arr, d1 [, d2, ...])
	// Natural order: dim(mat, rows [, cols]) — d1 is the outer (row) count.
	// For 2-D open arrays the comment in ejemplo_5 explains that the programmer
	// passes (cols, rows) — i.e. in reverse — so that dim(R, cant_cols, cant_filas)
	// actually creates cant_filas rows (outer) × cant_cols cols (inner).
	// We therefore do NOT reverse here; the caller is responsible for argument order.
	// Inner rows are created as *empty* arrays so that auto-extend during assignment
	// gives them exactly as many columns as the loops write.
	if len(args) < 2 {
		return zeroNil(), nil
	}

	// Get the array pointer (byRef[0])
	arrPtr := byRef[0]
	if arrPtr == nil {
		arrPtr = args[0]
	}

	dims := make([]int, len(args)-1)
	for i := 1; i < len(args); i++ {
		dims[i-1] = int(args[i].ToNum())
	}

	var newArr *Value
	if len(dims) == 1 {
		newArr = MakeArr(dims[0], KindNum)
	} else {
		// Create outer (rows) dimension; each row starts as an empty array.
		// Inner dimensions are grown on demand by resolveMultiIndex auto-extend.
		rows := dims[0]
		inner := make([]*Value, rows)
		for i := range inner {
			inner[i] = &Value{Kind: KindArr, Arr: []*Value{}}
		}
		newArr = &Value{Kind: KindArr, Arr: inner}
	}
	arrPtr.Assign(newArr)
	return zeroNil(), nil
}

// ---- Other ----

func builtinDec(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	ptr := byRef[0]
	if ptr == nil {
		ptr = args[0]
	}
	decr := 1.0
	if len(args) >= 2 {
		decr = args[1].ToNum()
	}
	ptr.Num -= decr
	return &Value{Kind: KindNum, Num: ptr.Num}, nil
}

func builtinInc(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	ptr := byRef[0]
	if ptr == nil {
		ptr = args[0]
	}
	incr := 1.0
	if len(args) >= 2 {
		incr = args[1].ToNum()
	}
	ptr.Num += incr
	return &Value{Kind: KindNum, Num: ptr.Num}, nil
}

func builtinIfval(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	// ifval is special: it should do lazy evaluation.
	// Here we receive pre-evaluated args (not fully lazy).
	// The parser/interpreter would need special handling for true laziness.
	// For now, treat as: if args[0] truthy, return args[1], else args[2].
	if len(args) < 3 {
		return zeroNil(), nil
	}
	if args[0].Truthy() {
		return args[1], nil
	}
	return args[2], nil
}

func builtinIntercambiar(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	a := byRef[0]
	b := byRef[1]
	if a == nil {
		a = args[0]
	}
	if b == nil {
		b = args[1]
	}
	if a == nil || b == nil {
		return zeroNil(), nil
	}
	tmp := a.Copy()
	a.Assign(b)
	b.Assign(tmp)
	return zeroNil(), nil
}

func builtinMax(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 2 {
		return zeroNum(), nil
	}
	if valLess(args[0], args[1]) {
		return args[1], nil
	}
	return args[0], nil
}

func builtinMin(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 2 {
		return zeroNum(), nil
	}
	if valLess(args[0], args[1]) {
		return args[0], nil
	}
	return args[1], nil
}

func builtinParamval(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return &Value{Kind: KindStr, Str: ""}, nil
	}
	pos := int(args[0].ToNum())
	if pos < 1 || pos > len(interp.cmdArgs) {
		return &Value{Kind: KindStr, Str: ""}, nil
	}
	return &Value{Kind: KindStr, Str: interp.cmdArgs[pos-1]}, nil
}

func builtinPcount(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	return &Value{Kind: KindNum, Num: float64(len(interp.cmdArgs))}, nil
}

func builtinRandom(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNum(), nil
	}
	n := int(args[0].ToNum())
	if n <= 0 {
		return zeroNum(), nil
	}
	return &Value{Kind: KindNum, Num: float64(rand.Intn(n))}, nil
}

func builtinRuncmd(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	if len(args) < 1 {
		return zeroNil(), nil
	}
	cmd := args[0].ToStr("")
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run() //nolint:errcheck
	return zeroNil(), nil
}

func builtinSec(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	elapsed := time.Since(startTime).Seconds()
	return &Value{Kind: KindNum, Num: elapsed}, nil
}

func builtinTerminar(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	msg := ""
	if len(args) > 0 {
		msg = args[0].ToStr("")
	}
	if msg != "" {
		interp.stdout.WriteString(msg + "\n")
	}
	return nil, &TerminateError{Msg: msg}
}

func builtinNoop(interp *Interpreter, args []*Value, byRef []*Value) (*Value, error) {
	return zeroNil(), nil
}
