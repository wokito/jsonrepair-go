# jsonrepair-go

[![Go Reference](https://pkg.go.dev/badge/github.com/wujinduo/jsonrepair-go.svg)](https://pkg.go.dev/github.com/wujinduo/jsonrepair-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/wujinduo/jsonrepair-go)](https://goreportcard.com/report/github.com/wujinduo/jsonrepair-go)
[![License: ISC](https://img.shields.io/badge/License-ISC-blue.svg)](https://opensource.org/licenses/ISC)

Go 语言实现的 JSON 修复库，用于修复无效的 JSON 文档。这是 [jsonrepair](https://github.com/josdejong/jsonrepair) (TypeScript/JavaScript) 的 Go 移植版本。

## 功能特性

jsonrepair 可以修复以下 JSON 问题：

- ✅ 为键名添加缺失的引号
- ✅ 添加缺失的转义字符
- ✅ 添加缺失的逗号
- ✅ 添加缺失的闭合括号
- ✅ 修复截断的 JSON
- ✅ 将单引号替换为双引号
- ✅ 替换特殊引号字符（如中文引号 `""`、`''`）
- ✅ 替换特殊空白字符
- ✅ 替换 Python 常量（`None`、`True`、`False`）
- ✅ 移除尾部逗号
- ✅ 移除注释（`/* */` 和 `//`）
- ✅ 移除 Markdown 代码块标记
- ✅ 移除数组和对象中的省略号 `...`
- ✅ 移除 JSONP 包装
- ✅ 处理 MongoDB 数据类型（`ObjectId`、`NumberLong` 等）
- ✅ 字符串连接（`"hello" + " world"`）
- ✅ 将换行分隔的 JSON 转换为有效的 JSON 数组

## 安装

```bash
go get github.com/wujinduo/jsonrepair-go
```

## 使用方法

### 作为库使用

```go
package main

import (
    "fmt"
    "log"

    "github.com/wujinduo/jsonrepair-go"
)

func main() {
    // 修复缺少引号的键
    input := `{name: "John", age: 30}`
    result, err := jsonrepair.JSONRepair(input)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result) // {"name": "John", "age": 30}

    // 使用 MustJSONRepair（失败时会 panic）
    result = jsonrepair.MustJSONRepair(`[1, 2, 3,]`)
    fmt.Println(result) // [1, 2, 3]
}
```

### 作为命令行工具使用

```bash
# 从标准输入读取
echo "{name: 'John'}" | jsonrepair

# 从文件读取
jsonrepair input.json

# 输出到文件
jsonrepair input.json -o output.json
```

## 示例

### 修复缺少引号的键

```go
input := `{name: "John", age: 30}`
result, _ := jsonrepair.JSONRepair(input)
// 输出: {"name": "John", "age": 30}
```

### 修复单引号

```go
input := `{'name': 'Alice'}`
result, _ := jsonrepair.JSONRepair(input)
// 输出: {"name": "Alice"}
```

### 移除尾部逗号

```go
input := `{"items": [1, 2, 3,]}`
result, _ := jsonrepair.JSONRepair(input)
// 输出: {"items": [1, 2, 3]}
```

### 移除注释

```go
input := `{
    "name": "Bob", // 用户名
    "active": true /* 是否激活 */
}`
result, _ := jsonrepair.JSONRepair(input)
// 输出: {"name": "Bob", "active": true}
```

### 修复截断的 JSON

```go
input := `{"message": "hello`
result, _ := jsonrepair.JSONRepair(input)
// 输出: {"message": "hello"}
```

### 修复 Python 常量

```go
input := `{"enabled": True, "data": None}`
result, _ := jsonrepair.JSONRepair(input)
// 输出: {"enabled": true, "data": null}
```

### 字符串连接

```go
input := `"hello" + " world"`
result, _ := jsonrepair.JSONRepair(input)
// 输出: "hello world"
```

### 处理 MongoDB 数据类型

```go
input := `{"_id": ObjectId("123"), "count": NumberLong("456")}`
result, _ := jsonrepair.JSONRepair(input)
// 输出: {"_id": "123", "count": "456"}
```

## API

### JSONRepair

```go
func JSONRepair(text string) (string, error)
```

修复无效的 JSON 字符串。如果无法修复，返回错误。

### MustJSONRepair

```go
func MustJSONRepair(text string) string
```

修复无效的 JSON 字符串。如果无法修复，会 panic。

## 与原始 TypeScript 版本的差异

本项目是 [josdejong/jsonrepair](https://github.com/josdejong/jsonrepair) 的 Go 移植版本，功能基本一致，但有以下已知差异：

| 功能 | TypeScript | Go |
|------|------------|-----|
| MongoDB Timestamp | ✅ 支持 | ❌ 不支持 |
| Triple quotes `'''` | ✅ 修复 | ❌ 报错 |
| 正则表达式转义 | 双反斜杠 | 单反斜杠 |

## 测试

```bash
go test -v
```

当前测试覆盖 78 个测试用例，全部通过。

## 许可证

ISC License

本项目是 [jsonrepair](https://github.com/josdejong/jsonrepair) 的 Go 移植版本，原项目由 Jos de Jong 创建，采用 ISC 许可证。

## 致谢

- [Jos de Jong](https://github.com/josdejong) - 原始 jsonrepair TypeScript 项目的作者
