package main

import (
	"fmt"
	"reflect"
)

func swap[T any](a, b T) {

	a_reflect := reflect.ValueOf(a)
	b_reflect := reflect.ValueOf(b)

	if a_reflect.Kind() != reflect.Ptr || b_reflect.Kind() != reflect.Ptr {
		panic("swap: both arguments must be pointers (call swap(&a, &b))")
	}

	if a_reflect.Type() != b_reflect.Type() {
		panic("swap: pointer types must be identical")
	}

	tmp_reflect := reflect.New(a_reflect.Elem().Type()).Elem()
	tmp_reflect.Set(a_reflect.Elem())
	a_reflect.Elem().Set(b_reflect.Elem())
	b_reflect.Elem().Set(tmp_reflect)
}

func main() {
	a := 10
	b := 20

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)

	swap(&a, &b)

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)
}
