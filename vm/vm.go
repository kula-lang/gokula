package vm

import (
	"fmt"
	"gokula/runtime"
	"gokula/utils"
	"strings"
)

var global *Context
var context *Context
var vmStack utils.Stack[*utils.Stack[any]]
var currentStack *utils.Stack[any]
var callStack utils.Stack[CallInfo]

var ip int
var fp int

type CallInfo struct {
	Ip, Fp  int
	Context *Context
}

func initVM() {
	global = NewContext(nil)
	context = global
	vmStack = utils.NewStack[*utils.Stack[any]]()
	innerStack := utils.NewStack[any]()
	currentStack = &innerStack
	vmStack.Push(currentStack)
	callStack = utils.NewStack[CallInfo]()
	ip = 0
	fp = -1
}

func (cf *CompiledFile) Run() error {
	initVM()

	// Standard Library
	initStdlib()

	for {
		var ins *Instruction
		if fp >= 0 {
			ins = &CompiledFileInstance.Functions[fp].Instructions[ip]
		} else {
			ins = &CompiledFileInstance.Chunk[ip]
		}
		// fmt.Println("Do", ins, "in [F", fp, "]")
		err := ins.run()
		if err != nil {
			return err
		}
		ip++
		if fp >= 0 {
			if ip >= len(CompiledFileInstance.Functions[fp].Instructions) {
				vmStack.Pop().Clear()
				currentStack = vmStack.Peek()
				currentStack.Push(nil)
				callInfo := callStack.Pop()
				ip = callInfo.Ip
				fp = callInfo.Fp
				context = callInfo.Context
				ip++
			}
		} else {
			if ip >= len(CompiledFileInstance.Chunk) {
				break
			}
		}
	}

	return nil
}

func (ins *Instruction) run() error {
	switch ins.Op {
	case LOADC:
		currentStack.Push(CompiledFileInstance.Literals[ins.Val])
	case LOAD:
		v, err := context.Get(CompiledFileInstance.SymbolArray[ins.Val])
		if err != nil {
			return err
		}
		currentStack.Push(v)
	case DECL:
		top := currentStack.Peek()
		context.Define(CompiledFileInstance.SymbolArray[ins.Val], top)
	case ASGN:
		top := currentStack.Peek()
		context.Assgin(CompiledFileInstance.SymbolArray[ins.Val], top)
	case POP:
		currentStack.Pop()
	case DUP:
		currentStack.Push(currentStack.Peek())
	case FUNC:
		f := NewFunction(ins.Val, context)
		currentStack.Push(f)
	case RET:
		vmStack.Pop().Clear()
		currentStack = vmStack.Peek()
		currentStack.Push(nil)
		callInfo := callStack.Pop()
		ip = callInfo.Ip
		fp = callInfo.Fp
		context = callInfo.Context
	case RETV:
		top := currentStack.Pop()
		vmStack.Pop().Clear()
		currentStack = vmStack.Peek()
		currentStack.Push(top)
		callInfo := callStack.Pop()
		ip = callInfo.Ip
		fp = callInfo.Fp
		context = callInfo.Context
	case ENVST:
		context = NewContext(context)
	case ENVED:
		context = context.enclosing
	case GET:
		key := currentStack.Pop()
		container := currentStack.Pop()
		value, err := evalGet(container, key, ins)
		if err != nil {
			return err
		}
		currentStack.Push(value)
	case GETWT:
		key := currentStack.Pop()
		container := currentStack.Pop()
		value, err := evalGet(container, key, ins)
		if err != nil {
			return err
		}
		currentStack.Push(container)
		currentStack.Push(value)
	case SET:
		value := currentStack.Pop()
		key := currentStack.Pop()
		container := currentStack.Pop()
		err := evalSet(container, key, value)
		if err != nil {
			return err
		}
		currentStack.Push(value)
	case CALL:
		argc := ins.Val
		argv := make([]any, argc)
		for c := argc - 1; c >= 0; c -= 1 {
			argv[c] = currentStack.Pop()
		}
		function := currentStack.Pop()

		if vmf, ok := function.(*VMFunction); ok {
			vmf.calcVMFunction(argv)
		} else if nf, ok := function.(*NativeFunction); ok {
			val, err := nf.calcNativeFunction(argv)
			if err != nil {
				return err
			}
			currentStack.Push(val)
		} else if object, ok := function.(*runtime.KulaObject); ok {
			key := runtime.FUNC__
			functionSugar := object.Get((*runtime.KulaString)(&key))
			if vmf, ok := functionSugar.(*VMFunction); ok {
				vmf.calcVMFunction(argv)
			} else {
				return fmt.Errorf("object has no such function")
			}
		} else {
			return fmt.Errorf("can only call functions")
		}
	case CALWT:
		argc := ins.Val
		argv := make([]any, argc)
		for c := argc - 1; c >= 0; c -= 1 {
			argv[c] = currentStack.Pop()
		}
		function := currentStack.Pop()
		callSite := currentStack.Pop()

		if vmf, ok := function.(*VMFunction); ok {
			vmf.CallSite = callSite
			vmf.calcVMFunction(argv)
		} else if nf, ok := function.(*NativeFunction); ok {
			nf.CallSite = callSite
			val, err := nf.calcNativeFunction(argv)
			if err != nil {
				return err
			}
			currentStack.Push(val)
		} else if object, ok := function.(*runtime.KulaObject); ok {
			key := runtime.FUNC__
			functionSugar := object.Get((*runtime.KulaString)(&key))
			if vmf, ok := functionSugar.(*VMFunction); ok {
				vmf.CallSite = &callSite
				vmf.calcVMFunction(argv)
			} else {
				return fmt.Errorf("object has no such function")
			}
		} else {
			return fmt.Errorf("can only call functions")
		}
	case PRINT:
		ls := make([]string, ins.Val)
		for t := ins.Val - 1; t >= 0; t-- {
			ls[t] = string(*runtime.Stringify(currentStack.Pop()))
		}
		fmt.Println(strings.Join(ls, " "))
	case JMP:
		ip = ins.Val - 1
	case JMPT:
		if runtime.Booleanify(currentStack.Pop()) {
			ip = ins.Val - 1
		}
	case JMPF:
		if !runtime.Booleanify(currentStack.Pop()) {
			ip = ins.Val - 1
		}
	// calculating
	case ADD:
		v2 := currentStack.Pop()
		v1 := currentStack.Pop()
		if n1, ok := v1.(runtime.KulaNumber); ok {
			if n2, ok := v2.(runtime.KulaNumber); ok {
				currentStack.Push(n1 + n2)
				break
			}
		}
		if s1, ok := v1.(*runtime.KulaString); ok {
			if s2, ok := v2.(*runtime.KulaString); ok {
				str := *s1 + *s2
				currentStack.Push(&str)
				break
			}
		}
		return fmt.Errorf("operands must be 2 numbers or 2 strings")
	case SUB:
		v2 := currentStack.Pop()
		v1 := currentStack.Pop()
		if n1, ok := v1.(runtime.KulaNumber); ok {
			if n2, ok := v2.(runtime.KulaNumber); ok {
				currentStack.Push(n1 - n2)
				break
			}
		}
		return fmt.Errorf("operands must be 2 numbers")
	case MUL:
		v2 := currentStack.Pop()
		v1 := currentStack.Pop()
		if n1, ok := v1.(runtime.KulaNumber); ok {
			if n2, ok := v2.(runtime.KulaNumber); ok {
				currentStack.Push(n1 * n2)
				break
			}
		}
		return fmt.Errorf("operands must be 2 numbers")
	case DIV:
		v2 := currentStack.Pop()
		v1 := currentStack.Pop()
		if n1, ok := v1.(runtime.KulaNumber); ok {
			if n2, ok := v2.(runtime.KulaNumber); ok {
				currentStack.Push(n1 / n2)
				break
			}
		}
		return fmt.Errorf("operands must be 2 numbers")
	case MOD:
		v2 := currentStack.Pop()
		v1 := currentStack.Pop()
		if n1, ok := v1.(runtime.KulaNumber); ok {
			if n2, ok := v2.(runtime.KulaNumber); ok {
				currentStack.Push(runtime.KulaNumber(int(n1) % int(n2)))
				break
			}
		}
		return fmt.Errorf("operands must be 2 numbers")
	case GT:
		v2 := currentStack.Pop()
		v1 := currentStack.Pop()
		if n1, ok := v1.(runtime.KulaNumber); ok {
			if n2, ok := v2.(runtime.KulaNumber); ok {
				currentStack.Push(runtime.KulaBool(n1 > n2))
				break
			}
		}
		return fmt.Errorf("operands must be 2 numbers")
	case GE:
		v2 := currentStack.Pop()
		v1 := currentStack.Pop()
		if n1, ok := v1.(runtime.KulaNumber); ok {
			if n2, ok := v2.(runtime.KulaNumber); ok {
				currentStack.Push(runtime.KulaBool(n1 >= n2))
				break
			}
		}
		return fmt.Errorf("operands must be 2 numbers")
	case LT:
		v2 := currentStack.Pop()
		v1 := currentStack.Pop()
		if n1, ok := v1.(runtime.KulaNumber); ok {
			if n2, ok := v2.(runtime.KulaNumber); ok {
				currentStack.Push(runtime.KulaBool(n1 < n2))
				break
			}
		}
		return fmt.Errorf("operands must be 2 numbers")
	case LE:
		v2 := currentStack.Pop()
		v1 := currentStack.Pop()
		if n1, ok := v1.(runtime.KulaNumber); ok {
			if n2, ok := v2.(runtime.KulaNumber); ok {
				currentStack.Push(runtime.KulaBool(n1 <= n2))
				break
			}
		}
		return fmt.Errorf("operands must be 2 numbers")
	case EQ:
		v2 := currentStack.Pop()
		v1 := currentStack.Pop()
		currentStack.Push(runtime.KulaBool(v1 == v2))
	case NEQ:
		v2 := currentStack.Pop()
		v1 := currentStack.Pop()
		currentStack.Push(runtime.KulaBool(v1 != v2))
	case NEG:
		top := currentStack.Pop()
		currentStack.Push(runtime.Booleanify(top))
	default:
		fmt.Println("err: ", ins.Op.String())
	}

	return nil
}

func evalGet(container any, key any, ins *Instruction) (any, error) {
	if object, ok := container.(*runtime.KulaObject); ok {
		if keyString, ok := key.(*runtime.KulaString); ok {
			return object.Get(keyString), nil
		}
		return nil, fmt.Errorf("index of 'Object' can only be 'String'")
	} else if array, ok := container.(*runtime.KulaArray); ok {
		if keyNumber, ok := key.(runtime.KulaNumber); ok {
			return array.Get(keyNumber), nil
		} else if keyString, ok := key.(*runtime.KulaString); ok {
			return runtime.ArrayProto.Get(keyString), nil
		}
		return nil, fmt.Errorf("index of 'Array' can only be 'Number'")
	}

	if keyString, ok := key.(*runtime.KulaString); ok {
		if _, ok := container.(*runtime.KulaString); ok {
			return runtime.StringProto.Get(keyString), nil
		} else if _, ok := container.(runtime.KulaNumber); ok {
			return runtime.NumberProto.Get(keyString), nil
		} else if _, ok := container.(runtime.KulaBool); ok {
			return runtime.BoolProto.Get(keyString), nil
		}
	}
	return nil, fmt.Errorf("what do you want to get?")
}

func evalSet(container any, key, value any) error {
	if object, ok := container.(*runtime.KulaObject); ok {
		if keyString, ok := key.(*runtime.KulaString); ok {
			object.Set(keyString, value)
			return nil
		}
	}
	if array, ok := container.(*runtime.KulaArray); ok {
		if keyNumber, ok := key.(runtime.KulaNumber); ok {
			array.Set(keyNumber, value)
			return nil
		}
	}
	return fmt.Errorf("cannot set key '%s' to container '%s'", key, container)
}

func (fn *VMFunction) calcVMFunction(argv []any) {
	callStack.Push(CallInfo{
		Ip:      ip,
		Fp:      fp,
		Context: context,
	})
	ip = -1
	fp = fn.Index
	context = NewContext(fn.Parent)

	fc := CompiledFileInstance.Functions[fn.Index]
	innerStack := utils.NewStack[any]()
	currentStack = &innerStack
	vmStack.Push(currentStack)
	for i := 0; i < len(argv); i++ {
		vIndex := fc.Params[i]
		vName := CompiledFileInstance.SymbolArray[vIndex]
		context.Define(vName, argv[i])
	}
	context.Define("self", fn)
	if fn.CallSite != nil {
		context.Define("this", fn.CallSite)
		fn.CallSite = nil
	}
}

func (nf *NativeFunction) calcNativeFunction(argv []any) (val any, err error) {
	return nf.Callee(nf.CallSite, argv)
}
