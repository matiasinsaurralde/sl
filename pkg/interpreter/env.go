package interpreter

// Env is a lexical scope / environment mapping names to value pointers.
type Env struct {
	vars   map[string]*Value
	parent *Env
}

// NewEnv creates a root environment.
func NewEnv() *Env {
	return &Env{vars: make(map[string]*Value)}
}

// Child creates a child scope inheriting from e.
func (e *Env) Child() *Env {
	return &Env{vars: make(map[string]*Value), parent: e}
}

// Define creates a new binding in the current scope.
func (e *Env) Define(name string, v *Value) {
	e.vars[name] = v
}

// Set assigns to the nearest scope that has name defined.
// If not found, defines in current scope.
func (e *Env) Set(name string, v *Value) {
	if ptr, ok := e.lookup(name); ok {
		ptr.Assign(v)
		return
	}
	e.vars[name] = v
}

// Get returns the value for name (searching up the chain).
func (e *Env) Get(name string) (*Value, bool) {
	ptr, ok := e.lookup(name)
	return ptr, ok
}

// GetPtr returns a pointer to the value cell for name (for ref params, lvalue).
func (e *Env) GetPtr(name string) (*Value, bool) {
	return e.lookup(name)
}

func (e *Env) lookup(name string) (*Value, bool) {
	for env := e; env != nil; env = env.parent {
		if v, ok := env.vars[name]; ok {
			return v, true
		}
	}
	return nil, false
}
