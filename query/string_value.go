package query

type CheckString func(v string) bool

func CheckStringNotEmpty(v string) bool {
	if v == "" {
		return false
	}
	return true
}

type stringValue struct {
	*baseValue

	strPtr    *string
	checkFunc CheckString
}

func NewStringValue(strPtr *string, required bool, errno string, msg string, checkFunc CheckString) *stringValue {
	s := &stringValue{
		baseValue: newBaseValue(required, errno, msg),

		strPtr:    strPtr,
		checkFunc: checkFunc,
	}

	return s
}

func (sv *stringValue) Set(str string) error {
	*(sv.strPtr) = str

	return nil
}

func (sv *stringValue) Check() bool {
	if sv.checkFunc == nil {
		return true
	}

	return sv.checkFunc(*(sv.strPtr))
}
