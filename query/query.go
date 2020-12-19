package query

import (
	"github.com/go-jar/goerror"

	"net/url"
	"strings"
)

type QuerySet struct {
	formal map[string]Value
	exists map[string]bool
}

func NewQuerySet() *QuerySet {
	qs := &QuerySet{
		formal: make(map[string]Value),
		exists: make(map[string]bool),
	}

	return qs
}

func (qs *QuerySet) ExistsInfo() map[string]bool {
	return qs.exists
}

func (qs *QuerySet) Exist(name string) bool {
	if v, ok := qs.exists[name]; ok && v {
		return true
	}

	return false
}

func (qs *QuerySet) Var(name string, v Value) *QuerySet {
	qs.formal[name] = v

	return qs
}

func (qs *QuerySet) IntVar(intPtr *int, name string, required bool, errno string, msg string, checkFunc CheckInt) *QuerySet {
	qs.Var(name, NewIntValue(intPtr, required, errno, msg, checkFunc))

	return qs
}

func (qs *QuerySet) StringVar(strPtr *string, name string, required bool, errno string, msg string, checkFunc CheckString) *QuerySet {
	qs.Var(name, NewStringValue(strPtr, required, errno, msg, checkFunc))

	return qs
}

func (qs *QuerySet) Int64Var(int64Ptr *int64, name string, required bool, errno string, msg string, checkFunc CheckInt64) *QuerySet {
	qs.Var(name, NewInt64Value(int64Ptr, required, errno, msg, checkFunc))

	return qs
}

func (qs *QuerySet) Parse(actual url.Values) *goerror.Error {
	for name, v := range qs.formal {
		if len(actual[name]) == 0 {
			if v.Required() {
				return v.Error()
			}
			continue
		}

		qs.exists[name] = true
		str := strings.TrimSpace(actual.Get(name))
		err := v.Set(str)
		if err != nil {
			return v.Error()
		}
		if v.Check() == false {
			return v.Error()
		}
	}

	return nil
}
