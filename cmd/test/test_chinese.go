package main

import (
	"encoding/json"
	"fmt"

	jsonrepair "github.com/wokito/jsonrepair-go"
)

func main() {
	// 原始输入 - 包含未转义的双引号
	input := `{ "镜号": "01", "调用图片": "[挂历]，[科比·布莱恩特海报]", "生图提示词": "特写，平拍，框架内构图，图中的[挂历]中心位置显示"2004年9月15日"，日期用红色标注，挂历纸张略微泛黄边缘有卷曲，挂历下方贴着图中的[科比·布莱恩特海报]中的科比穿着湖人队8号球衣年轻面孔充满斗志，墙面是斑驳的白色涂料有些地方露出水泥，均匀的自然光营造怀旧氛围，游戏截图风格。"}`

	fmt.Println("=== 原始输入 ===")
	fmt.Println(input)
	fmt.Println()

	// 尝试直接解析原始输入
	var data1 map[string]interface{}
	err1 := json.Unmarshal([]byte(input), &data1)
	fmt.Println("=== 原始 JSON 解析 ===")
	if err1 != nil {
		fmt.Println("失败:", err1)
	} else {
		fmt.Println("成功")
	}
	fmt.Println()

	// 修复
	result, repairErr := jsonrepair.JSONRepair(input)
	fmt.Println("=== 修复结果 ===")
	fmt.Println(result)
	fmt.Println()
	fmt.Println("=== 修复错误 ===")
	fmt.Printf("%v\n", repairErr)
	fmt.Println()

	// 尝试解析修复后的结果
	var data2 map[string]interface{}
	err2 := json.Unmarshal([]byte(result), &data2)
	fmt.Println("=== 修复后 JSON 解析 ===")
	if err2 != nil {
		fmt.Println("失败:", err2)
	} else {
		fmt.Println("成功")
		fmt.Println("生图提示词:", data2["生图提示词"])
	}
}
