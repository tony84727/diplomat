package reflecthelper

import "reflect"

func Actual(v reflect.Value) reflect.Value {
	for {
		switch v.Type().Kind() {
		case reflect.Interface, reflect.Ptr:
			v = v.Elem()
		default:
			return v
		}
	}
}
