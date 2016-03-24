package main

import (
	"fmt"
	"reflect"
)

type Get interface {
	Url(string) string
}

type Test struct {
	I int
	D float64
}

func (t Test) Url(s string) string {
	return s + " Test"
}

func (t *Test) SetI(i int) {
	t.I = i

	fmt.Println("设置Test.I = ", t.I)
}

func main() {
	var inter = Test{I: 10, D: 9.9}

	m := reflect.ValueOf(&inter).Elem()

	for i := 0; i < m.NumField(); i++ {
		fmt.Println(m.Field(i))
	}
	param := make([]reflect.Value, 1)
	param[0] = reflect.ValueOf("hehe")

	for i := 0; i < m.NumMethod(); i++ {
		fmt.Println(m.Method(i).Type())
	}
	// fmt.Println(m.MethodByName("Url").Call(param)[0])

	inter.SetI(12)

	fmt.Println(inter.I)

	var g Get
	g = inter

	fmt.Println(g.Url("hehe"))
	fmt.Println(reflect.TypeOf(&g))

	// s := reflect.TypeOf(inter).Elem()
	// for i := 0; i < s.NumField(); i++ {
	// 	fmt.Println(s.Field(i).Name, " : ", s.Field(i).Type)
	// 	fmt.Println(s.Field(i).Index)
	// }

}
