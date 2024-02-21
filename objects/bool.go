package objects

type KulaBool bool

var BoolProto *KulaObject = NewObject()

func Booleanify(v any) KulaBool {
	if v == nil {
		return false
	} else if b, ok := v.(KulaBool); ok {
		return b
	}
	return true
}
