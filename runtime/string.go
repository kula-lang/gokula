package runtime

import (
	"strconv"
	"strings"
)

type KulaString string

var StringProto *KulaObject = NewObject()

func Stringify(v any) *KulaString {
	var str KulaString
	if v == nil {
		str = "null"
		return &str
	} else if b, ok := v.(KulaBool); ok {
		str = KulaString(strconv.FormatBool(bool(b)))
		return &str
	} else if number, ok := v.(KulaNumber); ok {
		s1 := strconv.FormatFloat(float64(number), 'f', 8, 64)
		s1 = strings.TrimSuffix(strings.TrimRight(s1, "0"), ".")
		str = KulaString(s1)
		return &str
	} else if s, ok := v.(*KulaString); ok {
		return s
	} else if arr, ok := v.(*KulaArray); ok {
		str = KulaString(arr.String())
		return &str
	} else if obj, ok := v.(*KulaObject); ok {
		str = KulaString(obj.String())
		return &str
	}
	str = "<UnknownValue>"
	return &str
}

func (s *KulaString) At(index KulaNumber) *KulaString {
	i := int(index)
	str := (string)(*s)
	runes := []rune(str)
	r := string(runes[i])
	return (*KulaString)(&r)
}

func (s *KulaString) Cut(index, length KulaNumber) *KulaString {
	i := int(index)
	l := int(length)
	str := (string)(*s)
	runes := []rune(str)
	r := string(runes[i : i+l])
	return (*KulaString)(&r)
}

func (s *KulaString) Parse() (n KulaNumber) {
	num, err := strconv.ParseFloat(string(*s), 64)
	if err != nil {
		return
	}
	return FromFloat64(num)
}

func (s *KulaString) Split(seprator *KulaString) *KulaArray {
	strArr := strings.Split(string(*s), string(*seprator))
	arr := make([]any, len(strArr))
	for index, item := range strArr {
		arr[index] = item
	}
	return FromSlice(arr)
}

func (s *KulaString) Length() KulaNumber {
	return FromInt(len(string(*s)))
}

func (s *KulaString) CharCode(index KulaNumber) KulaNumber {
	runes := []rune(string(*s))
	r := runes[int(index)]
	return FromInt(int(r))
}
