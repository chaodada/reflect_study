package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"github.com/gookit/goutil/dump"
)

// 一文说反射
// 反射的本质就是在程序运行的时候，获取对象的类型信息和内存结构，反射是把双刃剑，功能强大但可读性差，反射代码无法在编译阶段静态发现错误，反射的代码常常比正常代码效率低1~2个数量级，
// 如果在关键位置使用反射会直接导致代码效率问题，所以，如非必要，不建议使用。

// 静态类型是指在编译的时候就能确定的类型（常见的变量声明类型都是静态类型）；
// 动态类型是指在运行的时候才能确定的类型（比如接口，也只有接口才有反射）。  # 只有接口才能反射

// ##########################################################################################
// 使用反射的三个步骤：
// 先有一个接口类型的变量
// 把它转成reflect对象 一般就是type 或者 value类型
// 然后根据不同的情况调用相应的函数

//func main() {
//	var number int
//	fmt.Println("type : ", reflect.TypeOf(number))   // 获取静态类型
//	fmt.Printf("%T\n", reflect.TypeOf(number))
//	fmt.Println("value : ", reflect.ValueOf(number)) // 获取变量值
//	fmt.Printf("%T\n", reflect.ValueOf(number))
//}

// ##########################################################################################
// 获取接口变量信息
// 事先知道原有类型的时候:
//func main() {
//	var num int = 3000
//	// 转成reflect Value 对象
//	value := reflect.ValueOf(num)
//	fmt.Printf("当前类型：%T , 当前值：%#v\n", value, value)
//	//从一个反射类型对象得到接口类型变量
//	v := value.Interface().(int)
//	fmt.Printf("当前类型：%T , 当前值：%#v\n", v, v)
//
//	Pval := reflect.ValueOf(&num)
//	fmt.Printf("当前类型：%T , 当前值：%#v\n", Pval, Pval)
//	pv := Pval.Interface().(*int)
//	fmt.Printf("当前类型：%T , 当前值：%#v\n", pv, pv) // 断言成了 int指针  值为内存地址
//	fmt.Printf("当前类型：%T , 当前值：%#v\n", *pv, *pv)
//}

// ##########################################################################################
// 事先不知道原有类型的时候:
// 这时候我们一般需要遍历探测一下Field

//type Person struct {
//	Name   string `json:"name" form:"name"`
//	Age    int    `json:"age" form:"name"`
//	Gender string `json:"gender" form:"name"`
//}
//
//func (p Person) Say(msg string) {
//	fmt.Println("hello, ", p.Name+msg)
//}
//func (p Person) Sum(a ,b int)  int{
//	return  a+b
//}
//func (p Person) PrintInfo() {
//	fmt.Printf("Name: %s, Age: %d, Gender: %s", p.Name, p.Age, p.Gender)
//}
//
//func Test(input interface{}) {
//
//	s := reflect.TypeOf(input)
//	fmt.Println()
//	fmt.Println("传入数据的类型", s.Kind())
//	fmt.Println()
//
//	// 判断是否指针类型
//	var getType reflect.Type
//	var getValue reflect.Value
//	if s.Kind() == reflect.Ptr {
//		// 转成reflect Type 对象
//		getType = reflect.TypeOf(input).Elem()
//
//		// 转成reflect Value 对象
//		getValue = reflect.ValueOf(input).Elem()
//
//	} else {
//		// 转成reflect Type 对象
//		getType = reflect.TypeOf(input)
//		// 转成reflect Value 对象
//		getValue = reflect.ValueOf(input)
//	}
//	fmt.Printf("当前类型：%T , 当前值：%#v\n", getType, getType)
//	fmt.Println("数据类型： ", getType.Name())
//	fmt.Println("类型种类： ", getType.Kind())
//	fmt.Printf("当前类型：%T , 当前值：%#v\n", getValue, getValue)
//
//	// 遍历结构体字段
//	fmt.Println("获取结构体字段：")
//	for i := 0; i < getType.NumField(); i++ {
//		field := getType.Field(i)
//		tag := getType.Field(i).Tag
//		value := getValue.Field(i).Interface()
//		fmt.Printf("字段名称： %s, 字段类型： %s, 字段值： %v, TAG： %v\n", field.Name, field.Type, value, tag)
//	}
//	// 遍历结构体方法
//	fmt.Println("获取结构体方法：")
//	for i := 0; i < getType.NumMethod(); i++ {
//		method := getType.Method(i)
//		fmt.Printf("方法名称： %s, 方法类型： %v\n", method.Name, method.Type)
//
//		agsNum := method.Type.NumIn() // 获取函数参数
//		fmt.Println("函数入参个数", agsNum)
//		var args []reflect.Value
//		switch agsNum {
//		case 1:
//			args = []reflect.Value{getValue}
//		case 2:
//			args = []reflect.Value{getValue, reflect.ValueOf("你好")}
//		case 3:
//			args = []reflect.Value{getValue, reflect.ValueOf(1),reflect.ValueOf(2)}
//		}
//
//		ret:=method.Func.Call(args)
//		for _, value := range ret {
//			fmt.Printf("函数执行返回值： 当前类型：%T , 当前值：%#v\n", value, value)
//			vv:=value.Interface().(int)
//			fmt.Printf("函数执行返回值： 当前类型：%T , 当前值：%#v\n", vv, vv)
//		}
//
//	}
//
//}
//
//func main() {
//
//	p1 := Person{"CDD", 16, "Male"}
//	Test(p1)
//	Test(&p1)
//
//
//}

// ##########################################################################################
// 创建结构体

type App struct {
	AppId int `json:"app_id" form:"app_id"`
}

func main() {
	app := new(App)
	if err := TestBindParam(app); err == nil {
		fmt.Printf("%#v\n",app)
		dump.P(app)
	}
}

func TestBindParam(i interface{}) error {
	s := reflect.TypeOf(i)
	if s.Kind() != reflect.Ptr {
		return errors.New("传入的不是一个指针类型")
	}
	ptrType := reflect.TypeOf(i) //获取call的指针的reflect.Type
	fmt.Printf("当前类型：%T , 当前值：%#v\n", ptrType, ptrType)
	trueType := ptrType.Elem() //获取type的真实类型
	fmt.Printf("当前类型：%T , 当前值：%#v\n", trueType, trueType)
	ptrValue := reflect.ValueOf(i)
	fmt.Printf("当前类型：%T , 当前值：%#v\n", ptrValue, ptrValue)
	trueValue := ptrValue.Elem() //获取真实的结构体类型
	fmt.Printf("当前类型：%T , 当前值：%#v\n", trueValue, trueValue)
	tp:=trueValue.FieldByName("AppId").Type()
	fmt.Printf("当前类型：%T , 当前值：%#v\n", tp, tp)
	switch tp.Kind()  {
	case reflect.String:
		trueValue.FieldByName("AppId").SetString("appid_xxxxxxx") // 设置值
	case reflect.Int:
		trueValue.FieldByName("AppId").SetInt(100)
	}

	return nil
}

func SetStructFieldByJsonName(ptr interface{}, fields map[string]interface{}) {
	v := reflect.ValueOf(ptr).Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("json")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}

		//去掉逗号后面内容 如 `json:"voucher_usage,omitempty"`
		name = strings.Split(name, ",")[0]
		if value, ok := fields[name]; ok {
			if reflect.ValueOf(value).Type() == v.FieldByName(fieldInfo.Name).Type() {
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value))
			}
		}
	}

	return

}
