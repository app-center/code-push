package flagkit

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

func MustBind(f interface{}, fs *flag.FlagSet) {
	if err := bindFlags(f, fs); err != nil {
		panic(err)
	}
}

func Bind(f interface{}, fs *flag.FlagSet) error {
	return bindFlags(f, fs)
}

func bindFlags(f interface{}, fs *flag.FlagSet) error {
	if fs == nil {
		fs = flag.CommandLine
	}

	v := reflect.ValueOf(f).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		name := field.Tag.Get("flag")
		dValue := field.Tag.Get("value")
		usage := field.Tag.Get("usage")

		if name == "" {
			continue
		}

		fieldAddr := v.Field(i).UnsafeAddr()

		switch field.Type.Kind() {
		case reflect.Bool:
			value, _ := strconv.ParseBool(dValue)
			fs.BoolVar((*bool)(unsafe.Pointer(fieldAddr)), name, value, usage)
		case reflect.String:
			fs.StringVar((*string)(unsafe.Pointer(fieldAddr)), name, dValue, usage)
		case reflect.Int:
			value, _ := strconv.Atoi(dValue)
			fs.IntVar((*int)(unsafe.Pointer(fieldAddr)), name, value, usage)
		default:
			return fmt.Errorf("type of field %s:%s do not support", field.Name, field.Type.String())
		}
	}

	return nil
}
