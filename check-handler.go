package justrpc

import "reflect"

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

func checkHandlerType(hType reflect.Type) {
	if hType.Kind() != reflect.Func {
		panic("handler must be a function, " + hType.Kind().String() + " provided")
	}

	if hType.NumIn() > 2 {
		panic("handler has invalid number of arguments")
	}

	if hType.NumOut() > 3 {
		panic("handler has invalid number of return values")
	}

	if hType.NumOut() == 3 && !hType.Out(2).Implements(errorInterface) {
		panic("the last of three return values must be error")
	}
}
