package vm

import (
	"fmt"
	"gokula/objects"
	"time"
)

func assert[T any](v any) (zero T, err error) {
	if val, ok := v.(T); ok {
		return val, nil
	}
	err = fmt.Errorf("wrong argument '%s' type", *objects.Stringify(v))
	return
}

func TypeOf(val any) *objects.KulaString {
	var str objects.KulaString
	if val == nil {
		str = "None"
	} else if _, ok := val.(objects.KulaBool); ok {
		str = "Bool"
	} else if _, ok := val.(objects.KulaNumber); ok {
		str = "Number"
	} else if _, ok := val.(*objects.KulaString); ok {
		str = "String"
	} else if _, ok := val.(*objects.KulaArray); ok {
		str = "Array"
	} else if _, ok := val.(*objects.KulaObject); ok {
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
			return objects.KulaNumber(float64(time.Now().UnixNano()-startTime) / 1000000000.0), nil
		}, 0,
	))
	global.Define("String", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return objects.Stringify(argv[0]), nil
		}, 1,
	))
	global.Define("Bool", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return objects.Booleanify(argv[0]), nil
		}, 1,
	))
	global.Define("Object", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return objects.NewObject(), nil
		}, 0,
	))
	global.Define("Array", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return objects.NewArray(), nil
		}, 0,
	))
	global.Define("asArray", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return objects.FromSlice(argv), nil
		}, -1,
	))
	global.Define("asObject", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			if len(argv)%2 == 1 {
				return nil, fmt.Errorf("need odd arguments but even is given")
			}
			obj := objects.NewObject()
			for i := 0; i+1 < len(argv); i += 2 {
				key := argv[i].(*objects.KulaString)
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
	global.Define("__string_proto__", objects.StringProto)
	global.Define("__array_proto__", objects.ArrayProto)
	global.Define("__number_proto__", objects.NumberProto)
	global.Define("__object_proto__", objects.ObjectProto)
	global.Define("__string_proto__", objects.StringProto)

	objects.StringProto.SetNative("at", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			str, err := assert[objects.KulaString](this)
			if err != nil {
				return nil, err
			}
			index, err := assert[objects.KulaNumber](argv[0])
			if err != nil {
				return nil, err
			}
			return str.At(index), nil
		}, 1,
	))
	objects.StringProto.SetNative("cut", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			str, err := assert[objects.KulaString](this)
			if err != nil {
				return nil, err
			}
			index, err := assert[objects.KulaNumber](argv[0])
			if err != nil {
				return nil, err
			}
			len, err := assert[objects.KulaNumber](argv[1])
			if err != nil {
				return nil, err
			}
			return str.Cut(index, len), nil
		}, 2,
	))

	objects.ArrayProto.SetNative("insert", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			arr := this.(*objects.KulaArray)
			index, err := assert[objects.KulaNumber](argv[0])
			if err != nil {
				return nil, err
			}
			value := argv[1]
			arr.Insert(index, value)
			return nil, nil
		}, 2,
	))
	objects.ArrayProto.SetNative("remove", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			arr := this.(*objects.KulaArray)
			index, err := assert[objects.KulaNumber](argv[0])
			if err != nil {
				return nil, err
			}
			arr.Remove(index)
			return nil, nil
		}, 2,
	))

	objects.ObjectProto.SetNative("copy", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			obj := this.(*objects.KulaObject)
			copied := objects.NewObject()
			for k, v := range *obj {
				copied.Set((*objects.KulaString)(&k), v)
			}
			return copied, nil
		}, 0,
	))
	objects.ObjectProto.SetNative("keys", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			obj := this.(*objects.KulaObject)
			arr := objects.NewArray()
			for k := range *obj {
				arr.Insert(arr.Length(), (*objects.KulaString)(&k))
			}
			return arr, nil
		}, 0,
	))
	objects.ObjectProto.SetNative("values", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			obj := this.(*objects.KulaObject)
			arr := objects.NewArray()
			for _, v := range *obj {
				arr.Insert(arr.Length(), v)
			}
			return arr, nil
		}, 0,
	))

	objects.NumberProto.SetNative("floor", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return this.(objects.KulaNumber).Floor(), nil
		}, 0,
	))
	objects.NumberProto.SetNative("round", NewNativeFunction(
		func(this any, argv []any) (any, error) {
			return this.(objects.KulaNumber).Round(), nil
		}, 0,
	))

	return nil
}
