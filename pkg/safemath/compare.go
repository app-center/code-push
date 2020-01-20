package safemath

import "reflect"

type compareInputType interface{}
type compareResultType int8

const (
	CompareLessFlag  = -1
	CompareEqualFlag = 0
	CompareLargeFlag = 1
)

func Compare(v1, v2 compareInputType) (result compareResultType) {
	defer func() {
		if r := recover(); r != nil {
			result = CompareLessFlag
		}
	}()

	if v1 == v2 {
		result = CompareEqualFlag
		return
	}

	kind1 := reflect.TypeOf(v1).Kind()
	kind2 := reflect.TypeOf(v2).Kind()

	rv1 := reflect.ValueOf(v1)
	rv2 := reflect.ValueOf(v2)

	if kind1 < reflect.Int || kind1 > reflect.Complex128 {
		result = CompareLessFlag
		return
	}

	if kind2 < reflect.Int || kind2 > reflect.Complex128 {
		result = CompareLargeFlag
		return
	}

	var majorKind reflect.Kind

	if kind1 >= kind2 {
		majorKind = kind1
	} else {
		majorKind = kind2
	}

	switch {
	case majorKind > reflect.Float64:
		result = CompareLessFlag
	case majorKind > reflect.Uintptr:
		vv1 := rv1.Float()
		vv2 := rv2.Float()

		if vv1 < vv2 {
			result = CompareLessFlag
		} else if vv1 > vv2 {
			result = CompareLargeFlag
		} else {
			result = CompareEqualFlag
		}
	case majorKind > reflect.Int64:
		vv1 := rv1.Uint()
		vv2 := rv2.Uint()

		if vv1 < vv2 {
			result = CompareLessFlag
		} else if vv1 > vv2 {
			result = CompareLargeFlag
		} else {
			result = CompareEqualFlag
		}
	case majorKind > reflect.Invalid:
		vv1 := rv1.Int()
		vv2 := rv2.Int()

		if vv1 < vv2 {
			result = CompareLessFlag
		} else if vv1 > vv2 {
			result = CompareLargeFlag
		} else {
			result = CompareEqualFlag
		}
	default:
		result = CompareLessFlag
	}

	return
}
