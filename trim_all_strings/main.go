package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func TrimAllStrings(a any) {
	visited := make(map[uintptr]bool)

	var trim func(v reflect.Value)
	trim = func(v reflect.Value) {
		if !v.IsValid() {
			return
		}

		switch v.Kind() {
		case reflect.Ptr:
			ptr := v.Pointer()
			if ptr != 0 {
				if visited[ptr] {
					return
				}
				visited[ptr] = true
			}
			trim(v.Elem())
		case reflect.Interface:
			trim(v.Elem())
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				trim(v.Field(i))
			}
		case reflect.String:
			v.SetString(strings.TrimSpace(v.String()))
		}
	}

	trim(reflect.ValueOf(a))

}

func main() {
	type Person struct {
		Name string
		Age  int
		Next *Person
	}

	a := &Person{
		Name: " name ",
		Age:  20,
		Next: &Person{
			Name: " name2 ",
			Age:  21,
			Next: &Person{
				Name: " name3 ",
				Age:  22,
			},
		},
	}

	TrimAllStrings(&a)

	m, _ := json.Marshal(a)

	fmt.Println(string(m))

	a.Next = a

	TrimAllStrings(&a)

	fmt.Println(a.Next.Next.Name == "name")
}
