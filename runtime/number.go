package runtime

import "math"

type KulaNumber float64

var NumberProto *KulaObject = NewObject()

func (n KulaNumber) Floor() KulaNumber {
	return FromFloat64(math.Floor(float64(n)))
}

func (n KulaNumber) Round() KulaNumber {
	return FromFloat64(math.Round(float64(n)))
}

func FromInt(i int) KulaNumber {
	return KulaNumber(i)
}

func FromFloat64(f float64) KulaNumber {
	return KulaNumber(f)
}
