package query

import "strconv"

type CheckInt64 func(v int64) bool

func CheckInt64IsPositive(v int64) bool {
	if v > 0 {
		return true
	}
	return false
}

type int64Value struct {
	*baseValue

	int64Ptr  *int64
	checkFunc CheckInt64
}

func NewInt64Value(int64Ptr *int64, required bool, errno string, msg string, checkFunc CheckInt64) *int64Value {
	iv := &int64Value{
		baseValue: newBaseValue(required, errno, msg),

		int64Ptr:  int64Ptr,
		checkFunc: checkFunc,
	}

	return iv
}

func (iv *int64Value) Set(str string) error {
	var v int64 = 0
	var e error = nil

	if str != "" {
		v, e = strconv.ParseInt(str, 10, 64)
	}

	if e != nil {
		return e
	}

	*(iv.int64Ptr) = v

	return nil
}

func (iv *int64Value) Check() bool {
	if iv.checkFunc == nil {
		return true
	}

	return iv.checkFunc(*(iv.int64Ptr))
}
