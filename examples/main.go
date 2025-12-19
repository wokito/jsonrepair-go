package main

import (
	"fmt"

	jsonrepair "github.com/wokito/jsonrepair-go"
)

func main() {
	// 示例 1: 修复缺少引号的键
	input1 := `{name: "John", age: 30}`
	result1, _ := jsonrepair.JSONRepair(input1)
	fmt.Println("1. 修复缺少引号的键:")
	fmt.Printf("   输入: %s\n", input1)
	fmt.Printf("   输出: %s\n\n", result1)

	// 示例 2: 修复单引号为双引号
	input2 := `{'name': 'Alice'}`
	result2, _ := jsonrepair.JSONRepair(input2)
	fmt.Println("2. 修复单引号:")
	fmt.Printf("   输入: %s\n", input2)
	fmt.Printf("   输出: %s\n\n", result2)

	// 示例 3: 修复尾部逗号
	input3 := `{"items": [1, 2, 3,]}`
	result3, _ := jsonrepair.JSONRepair(input3)
	fmt.Println("3. 修复尾部逗号:")
	fmt.Printf("   输入: %s\n", input3)
	fmt.Printf("   输出: %s\n\n", result3)

	// 示例 4: 移除注释
	input4 := `{
		"name": "Bob", // 用户名
		"active": true /* 是否激活 */
	}`
	result4, _ := jsonrepair.JSONRepair(input4)
	fmt.Println("4. 移除注释:")
	fmt.Printf("   输入: %s\n", input4)
	fmt.Printf("   输出: %s\n\n", result4)

	// 示例 5: 修复截断的 JSON
	input5 := `{"message": "hello`
	result5, _ := jsonrepair.JSONRepair(input5)
	fmt.Println("5. 修复截断的 JSON:")
	fmt.Printf("   输入: %s\n", input5)
	fmt.Printf("   输出: %s\n\n", result5)

	// 示例 6: 修复 Python 常量
	input6 := `{"enabled": True, "data": None}`
	result6, _ := jsonrepair.JSONRepair(input6)
	fmt.Println("6. 修复 Python 常量:")
	fmt.Printf("   输入: %s\n", input6)
	fmt.Printf("   输出: %s\n\n", result6)

	// 示例 7: 修复缺少的逗号
	input7 := `{"a": 1 "b": 2}`
	result7, _ := jsonrepair.JSONRepair(input7)
	fmt.Println("7. 修复缺少的逗号:")
	fmt.Printf("   输入: %s\n", input7)
	fmt.Printf("   输出: %s\n\n", result7)

	// 示例 8: 字符串连接
	input8 := `"hello" + " world"`
	result8, _ := jsonrepair.JSONRepair(input8)
	fmt.Println("8. 字符串连接:")
	fmt.Printf("   输入: %s\n", input8)
	fmt.Printf("   输出: %s\n\n", result8)

	// 示例 9: 处理 MongoDB 数据类型
	input9 := `{"_id": ObjectId("123"), "count": NumberLong("456")}`
	result9, _ := jsonrepair.JSONRepair(input9)
	fmt.Println("9. 处理 MongoDB 数据类型:")
	fmt.Printf("   输入: %s\n", input9)
	fmt.Printf("   输出: %s\n\n", result9)

	// 示例 10: 使用 MustJSONRepair (会 panic 如果无法修复)
	input10 := `[1, 2, 3`
	result10 := jsonrepair.MustJSONRepair(input10)
	fmt.Println("10. 使用 MustJSONRepair:")
	fmt.Printf("    输入: %s\n", input10)
	fmt.Printf("    输出: %s\n", result10)
}
