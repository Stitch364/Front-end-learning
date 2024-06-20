//package main
//
//import (
//	"fmt"
//	"reflect"
//)
//
//func main() {
//	//反射操作： 通过反射。可以获取一个接口类型变量的 类型和数值
//
//	var a int
//	a = 10
//	fmt.Println("type；", reflect.TypeOf(a))    //type； int
//	fmt.Println("value: ", reflect.ValueOf(a)) //value:  10
//
//	fmt.Println("------------------------")
//	//根据反射的值，来获取对应的类型和数值
//	v := reflect.ValueOf(a)
//	fmt.Println("kind is int", v.Kind() == reflect.Int) //kind is int true
//	fmt.Println("type:", v.Type())                      //type: int
//	fmt.Println("value:", v.Int())                      //value: 10
//
//}

//type Person struct {
//	Name string
//	Age  int
//	sex  string
//}
//
//func (p Person) Say(msg string) {
//	fmt.Println("hello ", msg)
//}
//
//func (p Person) PrintInfo(msg string) {
//	fmt.Printf("名字: %s, 年龄: %d，性别: %s\n", p.Name, p.Age, p.sex)
//}
//
//func main() {
//	p1 := Person{"Stitch", 26, "男"}
//	GetMessage(p1)
//}
//
//// 获取input的信息
//func GetMessage(input interface{}) {
//	getType := reflect.TypeOf(input) //先获取input 的类型
//	fmt.Println(getType.Name())      //Person
//	fmt.Println(getType.Kind())      //struct
//
//	getValue := reflect.ValueOf(input)
//	fmt.Println(getValue) //{Stitch 26 男}
//
//	//获取字段
//	//step1:先获取Type对象
//	//		NumField()		//返回字段总数
//	//		Field(index)
//	//step2:通过Filed()获取每一个Filed字段
//	//step3:Interface(),得到对应的Value
//
//	for i := 0; i < getType.NumField(); i++ {
//		field := getType.Field(i)
//		value := getValue.Field(i).Interface() //获取第i个值
//		fmt.Printf("字段名称:%s, 字段类型:%s, 字段数值:%v \n", field.Name, field.Type, value)
//	}
//
//	//获取方法
//	for i := 0; i < getType.NumMethod(); i++ {
//		method := getType.Method(i)
//		fmt.Println(method.Name, method.Type)
//	}
//}

package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
	Sex  string
}

func (p Person) Say(msg string) {
	fmt.Println("hello，", msg)
}
func (p Person) PrintInfo() {
	fmt.Printf("姓名：%s,年龄：%d，性别：%s\n", p.Name, p.Age, p.Sex)
}

func main() {
	p1 := Person{"王二狗", 30, "男"}

	DoFiledAndMethod(p1)

}

// DoFiledAndMethod 通过接口来获取任意参数
func DoFiledAndMethod(input interface{}) {

	getType := reflect.TypeOf(input)              //先获取input的类型
	fmt.Println("get Type is :", getType.Name())  // Person
	fmt.Println("get Kind is : ", getType.Kind()) // struct

	getValue := reflect.ValueOf(input)
	fmt.Println("get all Fields is:", getValue) //{王二狗 30 男}

	// 获取方法字段
	// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
	// 2. 再通过reflect.Type的Field获取其Field（字段）
	// 3. 最后通过Field的Interface()得到对应的value
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		value := getValue.Field(i).Interface() //获取第i个值
		fmt.Printf("字段名称:%s, 字段类型:%s, 字段数值:%v \n", field.Name, field.Type, value)
	}

	// 通过反射，操作方法
	// 1. 先获取interface的reflect.Type，然后通过.NumMethod进行遍历
	// 2. 再通过reflect.Type的Method获取其Method（方法）
	for i := 0; i < getType.NumMethod(); i++ {
		method := getType.Method(i)
		fmt.Printf("方法名称:%s, 方法类型:%v \n", method.Name, method.Type)
	}
}
