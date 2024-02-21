package vm

type VMFunction struct {
	Index    int
	Parent   *Context
	CallSite any
}

func NewFunction(index int, parent *Context) *VMFunction {
	f := new(VMFunction)
	f.Index = index
	f.Parent = context
	f.CallSite = nil
	return f
}

func (f *VMFunction) String() string {
	return "<Function>"
}

type NativeLambda func(this any, argv []any) (any, error)

type NativeFunction struct {
	Callee   NativeLambda
	Arity    int8
	CallSite any
}

func NewNativeFunction(callee NativeLambda, arity int8) *NativeFunction {
	f := new(NativeFunction)
	f.Callee = callee
	f.Arity = arity
	f.CallSite = nil
	return f
}

func (f *NativeFunction) String() string {
	return "<Function>"
}
