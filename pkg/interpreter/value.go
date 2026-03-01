// Package interpreter implements the SL tree-walking interpreter.
package interpreter

import (
	"fmt"
	"math"
	"strings"
)

// Kind identifies the runtime type of a Value.
type Kind int8

const (
	KindNil    Kind = iota // uninitialized
	KindNum                // numerico (float64)
	KindStr                // cadena (string)
	KindBool               // logico (bool)
	KindArr                // vector or matrix row — slice of *Value (1-indexed)
	KindRecord             // registro — field slice + names
)

// Value is the universal SL runtime value.
type Value struct {
	Kind   Kind
	Num    float64
	Str    string
	Bool   bool
	Arr    []*Value // KindArr: elements (1-indexed, so Arr[0] is element 1)
	FNames []string // KindRecord: field names in order
	Fields []*Value // KindRecord: field values parallel to FNames
}

// zeroNum returns the default numeric value.
func zeroNum() *Value { return &Value{Kind: KindNum, Num: 0} }

// zeroStr returns the default string value.
func zeroStr() *Value { return &Value{Kind: KindStr, Str: ""} }

// zeroBool returns the default bool value.
func zeroBool() *Value { return &Value{Kind: KindBool, Bool: false} }

// zeroNil returns an uninitialized value.
func zeroNil() *Value { return &Value{Kind: KindNil} }

// Copy returns a deep copy of v.
func (v *Value) Copy() *Value {
	if v == nil {
		return zeroNil()
	}
	out := &Value{Kind: v.Kind, Num: v.Num, Str: v.Str, Bool: v.Bool}
	if v.Arr != nil {
		out.Arr = make([]*Value, len(v.Arr))
		for i, e := range v.Arr {
			out.Arr[i] = e.Copy()
		}
	}
	if v.FNames != nil {
		out.FNames = append([]string{}, v.FNames...)
		out.Fields = make([]*Value, len(v.Fields))
		for i, f := range v.Fields {
			out.Fields[i] = f.Copy()
		}
	}
	return out
}

// Assign copies the contents of src into dst (in-place, preserving pointer).
func (dst *Value) Assign(src *Value) {
	if src == nil {
		src = zeroNil()
	}
	dst.Kind = src.Kind
	dst.Num = src.Num
	dst.Str = src.Str
	dst.Bool = src.Bool
	// deep copy arrays/records
	if src.Arr != nil {
		dst.Arr = make([]*Value, len(src.Arr))
		for i, e := range src.Arr {
			dst.Arr[i] = e.Copy()
		}
	} else {
		dst.Arr = nil
	}
	if src.FNames != nil {
		dst.FNames = append([]string{}, src.FNames...)
		dst.Fields = make([]*Value, len(src.Fields))
		for i, f := range src.Fields {
			dst.Fields[i] = f.Copy()
		}
	} else {
		dst.FNames = nil
		dst.Fields = nil
	}
}

// Truthy returns the boolean interpretation of v.
func (v *Value) Truthy() bool {
	switch v.Kind {
	case KindBool:
		return v.Bool
	case KindNum:
		return v.Num != 0
	case KindStr:
		return v.Str != ""
	}
	return false
}

// ToNum coerces v to float64.
func (v *Value) ToNum() float64 {
	switch v.Kind {
	case KindNum:
		return v.Num
	case KindBool:
		if v.Bool {
			return 1
		}
		return 0
	}
	return 0
}

// ToStr returns the string representation of v for imprimir.
func (v *Value) ToStr(ofs string) string {
	switch v.Kind {
	case KindNum:
		return formatNum(v.Num)
	case KindStr:
		return v.Str
	case KindBool:
		if v.Bool {
			return "TRUE"
		}
		return "FALSE"
	case KindArr:
		parts := make([]string, len(v.Arr))
		for i, e := range v.Arr {
			parts[i] = e.ToStr(ofs)
		}
		return strings.Join(parts, ofs)
	case KindRecord:
		parts := make([]string, len(v.Fields))
		for i, f := range v.Fields {
			parts[i] = f.ToStr(ofs)
		}
		return strings.Join(parts, ofs)
	case KindNil:
		return ""
	}
	return ""
}

// formatNum formats a float64 as SL would: integer if whole, else decimal.
func formatNum(f float64) string {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return fmt.Sprintf("%v", f)
	}
	// If it's a whole number, print without decimals
	if f == math.Trunc(f) && math.Abs(f) < 1e15 {
		return fmt.Sprintf("%d", int64(f))
	}
	s := fmt.Sprintf("%g", f)
	return s
}

// ArrLen returns the element count of an array (0 for nil/uninitialized).
func (v *Value) ArrLen() int {
	if v == nil || v.Kind != KindArr || v.Arr == nil {
		return 0
	}
	return len(v.Arr)
}

// Index returns a pointer to the i-th element (1-indexed).
// Returns nil if out of bounds.
func (v *Value) Index(i int) *Value {
	if v == nil || v.Kind != KindArr {
		return nil
	}
	if i < 1 || i > len(v.Arr) {
		return nil
	}
	return v.Arr[i-1]
}

// MakeArr creates a 1-D array of n zero-value elements (type-inited by kind).
func MakeArr(n int, elemKind Kind) *Value {
	arr := make([]*Value, n)
	for i := range arr {
		arr[i] = makeZeroByKind(elemKind)
	}
	return &Value{Kind: KindArr, Arr: arr}
}

func makeZeroByKind(k Kind) *Value {
	switch k {
	case KindNum:
		return zeroNum()
	case KindStr:
		return zeroStr()
	case KindBool:
		return zeroBool()
	}
	return zeroNil()
}

// MakeMatrix creates a matrix with dims[0] rows, dims[1] cols (2-D only for now).
func MakeMatrix(dims []int, elemKind Kind) *Value {
	if len(dims) == 0 {
		return zeroNil()
	}
	rows := dims[0]
	if len(dims) == 1 {
		return MakeArr(rows, elemKind)
	}
	arr := make([]*Value, rows)
	for i := range arr {
		arr[i] = MakeMatrix(dims[1:], elemKind)
	}
	return &Value{Kind: KindArr, Arr: arr}
}
