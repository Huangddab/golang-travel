package main

import (
	"fmt"
	"reflect"
)

/*
// type any = interface{}

// 获取反射对象类型reflect.Type
// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i any) Type

// 获取反射对象值reflect.Value
// ValueOf returns a new Value initialized to the concrete value stored in the interface i.
// ValueOf(nil) returns the zero Value.
func ValueOf(i any) Value

// 反射对象转换回interface
func (v Value) Interface() (i any)
*/

// func main() {
// 	age := 19
// 	fmt.Println("type:", reflect.TypeOf(age))
// 	value := reflect.ValueOf(age)
// 	fmt.Println("value:", value)
// 	fmt.Println(value.Interface().(int))
// }

/*
通过reflect.Value的SetXX相关方法，可以设置真实变量的值。reflect.Value是通过reflect.ValueOf(x)获得的，只有当x是指针的时候，才可以通过reflec.Value修改实际变量x的值。
// Set assigns x to the value v. It panics if Value.CanSet returns false.
// As in Go, x's value must be assignable to v's type and must not be derived from an unexported field.
func (v Value) Set(x Value)
func (v Value) SetInt(x int64)
...

// Elem returns the value that the interface v contains or that the pointer v points to. It panics if v's Kind is not Interface or Pointer. It returns the zero Value if v is nil.
func (v Value) Elem() Value

*/

func main() {
	age := 19
	// 参数必须是指针才能修改值
	pointerValue := reflect.ValueOf(&age)
	// Elem和Set方法结合，相当于给指针指向的变量赋值*p=值
	newValue := pointerValue.Elem()
	newValue.SetInt(28)
	fmt.Println(age) // 值被修改
	// reflect.ValueOf参数不是指针
	pointerValue = reflect.ValueOf(age)
	fmt.Println(pointerValue)
	// 如果非指针，直接panic: reflect: call of reflect.Value.Elem on int Value
	// newValue = pointerValue.Elem()
}
