package vm

import (
	"fmt"
	"gokula/runtime"
	"time"
)

func assert[T any](v any) (zero T, err error) {
	if val, ok := v.(T); ok {
		return val, nil
	}
	err = fmt.Errorf("wrong argument '%s' type", *runtime.Stringify(v))
	return
}

func TypeOf(val any) *runtime.KulaString {
	var str runtime.KulaString
	if val == nil {
		str = "None"
	} else if _, ok := val.(runtime.KulaBool); ok {
		str = "Bool"
	} else if _, ok := val.(runtime.KulaNumber); ok {
		str = "Number"
	} else if _, ok := val.(*runtime.KulaString); ok {
		str = "String"
	} else if _, ok := val.(*runtime.KulaArray); ok {
		str = "Array"
	} else if _, ok := val.(*runtime.KulaObject); ok {
		str = "Object"
	} else if _, ok := val.(*VMFunction); ok {
		str = "Function"
	} else if _, ok := val.(*NativeFunction); ok {
		str = "Function"
	}
	return &str
}

func initStdlib() error {
	startTime := time.Now().UnixNano()
	global.Define("clock", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return runtime.KulaNumber(float64(time.Now().UnixNano()-startTime) / 1000000000.0), nil
		}, 0,
	))
	global.Define("String", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return runtime.Stringify(argv[0]), nil
		}, 1,
	))
	global.Define("Bool", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return runtime.Booleanify(argv[0]), nil
		}, 1,
	))
	global.Define("Object", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return runtime.NewObject(), nil
		}, 0,
	))
	global.Define("Array", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return runtime.NewArray(), nil
		}, 0,
	))
	global.Define("asArray", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return runtime.FromSlice(argv), nil
		}, -1,
	))
	global.Define("asObject", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			if len(argv)%2 == 1 {
				return nil, fmt.Errorf("need odd arguments but even is given")
			}
			obj := runtime.NewObject()
			for i := 0; i+1 < len(argv); i += 2 {
				key := argv[i].(*runtime.KulaString)
				obj.Set(key, argv[i+1])
			}
			return obj, nil
		}, -1,
	))
	global.Define("typeof", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return TypeOf(this), nil
		}, 1,
	))
	global.Define("__string_proto__", runtime.StringProto)
	global.Define("__array_proto__", runtime.ArrayProto)
	global.Define("__number_proto__", runtime.NumberProto)
	global.Define("__object_proto__", runtime.ObjectProto)
	global.Define("__string_proto__", runtime.StringProto)

	runtime.StringProto.SetNative("at", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			str, err := assert[runtime.KulaString](this)
			if err != nil {
				return nil, err
			}
			index, err := assert[runtime.KulaNumber](argv[0])
			if err != nil {
				return nil, err
			}
			return str.At(index), nil
		}, 1,
	))
	runtime.StringProto.SetNative("cut", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			str, err := assert[runtime.KulaString](this)
			if err != nil {
				return nil, err
			}
			index, err := assert[runtime.KulaNumber](argv[0])
			if err != nil {
				return nil, err
			}
			len, err := assert[runtime.KulaNumber](argv[1])
			if err != nil {
				return nil, err
			}
			return str.Cut(index, len), nil
		}, 2,
	))

	runtime.ArrayProto.SetNative("insert", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			arr := this.(*runtime.KulaArray)
			index, err := assert[runtime.KulaNumber](argv[0])
			if err != nil {
				return nil, err
			}
			value := argv[1]
			arr.Insert(index, value)
			return nil, nil
		}, 2,
	))
	runtime.ArrayProto.SetNative("remove", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			arr := this.(*runtime.KulaArray)
			index, err := assert[runtime.KulaNumber](argv[0])
			if err != nil {
				return nil, err
			}
			arr.Remove(index)
			return nil, nil
		}, 2,
	))

	runtime.ObjectProto.SetNative("copy", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			obj := this.(*runtime.KulaObject)
			copied := runtime.NewObject()
			for k, v := range *obj {
				copied.Set((*runtime.KulaString)(&k), v)
			}
			return copied, nil
		}, 0,
	))
	runtime.ObjectProto.SetNative("keys", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			obj := this.(*runtime.KulaObject)
			arr := runtime.NewArray()
			for k := range *obj {
				arr.Insert(arr.Length(), (*runtime.KulaString)(&k))
			}
			return arr, nil
		}, 0,
	))
	runtime.ObjectProto.SetNative("values", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			obj := this.(*runtime.KulaObject)
			arr := runtime.NewArray()
			for _, v := range *obj {
				arr.Insert(arr.Length(), v)
			}
			return arr, nil
		}, 0,
	))

	runtime.NumberProto.SetNative("floor", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return this.(runtime.KulaNumber).Floor(), nil
		}, 0,
	))
	runtime.NumberProto.SetNative("round", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return this.(runtime.KulaNumber).Round(), nil
		}, 0,
	))

	return nil
}
