package query

import "strconv"

type CheckInt func(v int) bool

func CheckIntIsPositive(v int) bool {
	if v > 0 {
		return true
	}
	return false
}

func CheckIntGreaterEqual0(v int) bool {
	if v >= 0 {
		return true
	}
	return false
}

type intValue struct {
	*baseValue

	intPtr    *int
	checkFunc CheckInt
}

func NewIntValue(intPtr *int, required bool, errno string, msg string, checkFunc CheckInt) *intValue {
	iv := &intValue{
		baseValue: newBaseValue(required, errno, msg),

		intPtr:    intPtr,
		checkFunc: checkFunc,
	}

	return iv
}

func (iv *intValue) Set(str string) error {
	var v int = 0
	var e error = nil

	if str != "" {
		v, e = strconv.Atoi(str)
	}

	if e != nil {
		return e
	}

	*(iv.intPtr) = v

	return nil
}

func (iv *intValue) Check() bool {
	if iv.checkFunc == nil {
		return true
	}

	return iv.checkFunc(*(iv.intPtr))
}
