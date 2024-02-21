package objects

import "strings"

type KulaArray []any

var ArrayProto *KulaObject = NewObject()

func NewArray() *KulaArray {
	a := make([]any, 0)
	return (*KulaArray)(&a)
}

func FromSlice(slice []any) *KulaArray {
	return (*KulaArray)(&slice)
}

func (a *KulaArray) Get(index KulaNumber) any {
	i := int(index)
	if i >= 0 && i < len(*a) {
		return (*a)[i]
	}
	return nil
}

func (a *KulaArray) Set(index KulaNumber, value any) {
	i := int(index)
	if i >= 0 && i < len(*a) {
		(*a)[i] = value
	}
}

func (a *KulaArray) Insert(index KulaNumber, value any) {
	i := int(index)
	if i >= 0 && i <= len(*a) {
		(*a) = append((*a)[:i], append([]any{value}, (*a)[i:]...)...)
	}
}

func (a *KulaArray) Remove(index KulaNumber) {
	i := int(index)
	if i >= 0 && i < len(*a) {
		(*a) = append((*a)[:i], (*a)[i+1:]...)
	}
}

func (a *KulaArray) Length() KulaNumber {
	return FromInt(len(*a))
}

func (a *KulaArray) String() string {
	var sb strings.Builder
	sb.Grow(64)
	sb.WriteByte('[')
	slice := make([]string, len(*a))
	for index, item := range *a {
		str := string(*Stringify(item))
		if _, ok := item.(*KulaString); ok {
			str = "\"" + str + "\""
		}
		slice[index] = str
	}
	sb.Write([]byte(strings.Join(slice, ",")))
	sb.WriteByte(']')
	return sb.String()
}
