package objects

import "strings"

type KulaObject map[string]any

var ObjectProto = newObject()

func newObject() *KulaObject {
	data := make(map[string]any)
	return (*KulaObject)(&data)
}

func NewObject() *KulaObject {
	obj := newObject()
	obj.SetNative(PROTO__, ObjectProto)
	return obj
}

func (obj *KulaObject) Get(key *KulaString) any {
	val, ok := (*obj)[string(*key)]
	if ok {
		return val
	}
	val, ok = (*obj)[PROTO__]
	if !ok {
		return nil
	}
	proto, ok := val.(*KulaObject)
	if ok {
		return proto.Get(key)
	}
	return nil
}

func (obj *KulaObject) Set(key *KulaString, value any) {
	(*obj)[string(*key)] = value
}

func (obj *KulaObject) SetNative(key string, value any) {
	(*obj)[key] = value
}

func (obj *KulaObject) String() string {
	var sb strings.Builder
	sb.Grow(64)
	sb.WriteByte('{')
	slice := make([]string, 0)
	for key, value := range *obj {
		str := string(*Stringify(value))
		if _, ok := value.(*KulaString); ok {
			str = "\"" + str + "\""
		}
		str = "\"" + key + "\":" + str
		slice = append(slice, str)
	}
	sb.Write([]byte(strings.Join(slice, ",")))
	sb.WriteByte('}')
	return sb.String()
}
