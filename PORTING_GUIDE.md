# JSONRepair Go 移植指南

## 项目概述

### 原项目信息
- **项目名称**: jsonrepair
- **原作者**: Jos de Jong
- **GitHub仓库**: https://github.com/josdejong/jsonrepair
- **编程语言**: TypeScript/JavaScript
- **许可证**: ISC License
- **Star数**: 2.1k+

### 核心功能
JSONRepair 是一个用于修复无效 JSON 文档的库。它能够智能地修复各种 JSON 格式错误，使其变为有效的 JSON。

## 功能特性列表

### 1. 引号相关修复
- ✅ 为对象键添加缺失的引号
- ✅ 替换单引号为双引号
- ✅ 替换特殊引号字符（如 `"..."`, `'...'`, `` `...` ``）为标准双引号
- ✅ 为字符串添加缺失的转义字符
- ✅ 修复未转义的引号
- ✅ 修复转义字符串（如 `{\"stringified\": \"content\"}`）

### 2. 标点符号修复
- ✅ 添加缺失的逗号
- ✅ 移除尾随逗号
- ✅ 添加缺失的冒号

### 3. 括号修复
- ✅ 添加缺失的闭合括号（`]` 和 `}`）
- ✅ 移除多余的闭合括号

### 4. 注释与特殊标记移除
- ✅ 移除块注释 `/* ... */`
- ✅ 移除行注释 `// ...`
- ✅ 移除代码围栏块（Markdown fenced code blocks）` ```json ... ``` `
- ✅ 移除省略号 `[1, 2, 3, ...]`
- ✅ 移除 JSONP 包装 `callback({ ... })`

### 5. 数据类型修复
- ✅ 替换 Python 常量（`None` → `null`, `True` → `true`, `False` → `false`）
- ✅ 修复未引号字符串（如 `{name: value}` → `{"name": "value"}`）
- ✅ 移除 MongoDB 数据类型（如 `NumberLong(2)`, `ISODate(...)`）
- ✅ 修复截断的数字（如 `2.` → `2.0`, `2e` → `2e0`）
- ✅ 修复前导零数字（如 `00789` → `"00789"`）
- ✅ 将 `undefined` 转换为 `null`

### 6. 空白字符处理
- ✅ 替换特殊空白字符为标准空格
- ✅ 保留必要的空白
- ✅ 在适当位置插入空白

### 7. 字符串处理
- ✅ 连接字符串（如 `"long text" + "more text"` → `"long textmore text"`）
- ✅ 修复 Unicode 转义序列
- ✅ 处理 URL 字符串（避免将 `://` 误判为注释）
- ✅ 支持正则表达式字面量

### 8. 特殊格式支持
- ✅ 修复换行分隔的 JSON (NDJSON) 为标准 JSON 数组
- ✅ 修复截断的 JSON

### 9. 控制字符处理
- ✅ 转义控制字符（`\b`, `\f`, `\n`, `\r`, `\t`）

## 技术架构

### 核心实现

原项目采用**状态机 + 递归下降解析器**的设计模式：

#### 1. 解析器结构
```
jsonrepair()
  ├── parseValue()
  │   ├── parseObject()
  │   ├── parseArray()
  │   ├── parseString()
  │   ├── parseNumber()
  │   ├── parseKeywords()
  │   ├── parseUnquotedString()
  │   └── parseRegex()
  ├── parseWhitespaceAndSkipComments()
  ├── parseMarkdownCodeBlock()
  └── parseNewlineDelimitedJSON()
```

#### 2. 核心变量
- `i`: 当前解析位置索引（输入指针）
- `output`: 修复后的 JSON 字符串（输出缓冲区）
- `text`: 输入的待修复 JSON 字符串

#### 3. 主要算法流程

1. **初始化**
   - 设置输入指针 `i = 0`
   - 初始化输出 `output = ""`

2. **预处理**
   - 跳过可能存在的 Markdown 代码块开始标记

3. **解析主值**
   - 递归调用 `parseValue()` 解析根级值
   - 如果解析失败，抛出错误

4. **后处理**
   - 跳过可能存在的 Markdown 代码块结束标记
   - 处理尾随逗号
   - 检测并处理换行分隔的 JSON（NDJSON）
   - 移除多余的闭合括号

5. **验证完成**
   - 确保所有输入都已处理
   - 返回修复后的 JSON

### 字符串解析策略

字符串解析是最复杂的部分，采用**双阶段解析策略**：

#### 第一阶段：乐观解析
- 假设字符串有有效的结束引号
- 正常解析直到遇到匹配的引号

#### 第二阶段：保守解析
- 如果第一阶段失败（没有找到有效的结束引号）
- 在第一个遇到的分隔符处停止
- 插入缺失的引号

### 关键修复技术

#### 1. 插入策略
```typescript
// 在最后一个空白字符前插入内容
insertBeforeLastWhitespace(output, ',')
```

#### 2. 删除策略
```typescript
// 删除最后一次出现的内容
stripLastOccurrence(output, ',')
```

#### 3. 前瞻/后顾
解析器会查看前后字符来决定如何修复：
- 检查分隔符判断是否需要插入逗号
- 检查引号后的字符判断是否为真正的结束引号

## Go 移植计划

### 包结构设计

```
jsonrepair-go/
├── go.mod
├── go.sum
├── README.md
├── PORTING_GUIDE.md
├── LICENSE
├── jsonrepair.go          # 主入口和核心修复函数
├── parser.go              # 解析器实现
├── stringutils.go         # 字符串工具函数
├── types.go               # 类型定义和错误处理
├── constants.go           # 常量定义
├── examples/              # 示例代码
│   └── basic/
│       └── main.go
└── jsonrepair_test.go     # 测试文件
```

### 类型映射

| TypeScript | Go | 说明 |
|-----------|-----|------|
| `string` | `string` | 字符串 |
| `number` | `int` | 索引位置 |
| `boolean` | `bool` | 布尔值 |
| `Error` | `error` | 错误类型 |
| `{ [key: string]: string }` | `map[string]string` | 字符映射 |

### 核心类型定义

```go
// JSONRepairError 自定义错误类型
type JSONRepairError struct {
    Message  string
    Position int
}

// Parser 解析器结构体
type Parser struct {
    text   string       // 输入文本
    output strings.Builder // 输出缓冲区
    i      int          // 当前位置索引
}
```

### 函数映射

#### 主函数
- `jsonrepair(text: string): string` → `func JSONRepair(text string) (string, error)`

#### 解析函数
- `parseValue()` → `func (p *Parser) parseValue() bool`
- `parseObject()` → `func (p *Parser) parseObject() bool`
- `parseArray()` → `func (p *Parser) parseArray() bool`
- `parseString()` → `func (p *Parser) parseString(stopAtDelimiter bool, stopAtIndex int) bool`
- `parseNumber()` → `func (p *Parser) parseNumber() bool`
- `parseKeywords()` → `func (p *Parser) parseKeywords() bool`

#### 工具函数
- `isWhitespace()` → `func isWhitespace(text string, index int) bool`
- `isDelimiter()` → `func isDelimiter(char byte) bool`
- `isQuote()` → `func isQuote(char rune) bool`
- `insertBeforeLastWhitespace()` → `func insertBeforeLastWhitespace(text, textToInsert string) string`
- `stripLastOccurrence()` → `func stripLastOccurrence(text, textToStrip string, stripRemainingText bool) string`

### 实现要点

#### 1. Unicode 处理
Go 中需要注意：
- 使用 `rune` 类型处理 Unicode 字符
- 使用 `[]rune(text)` 转换字符串以正确处理多字节字符
- 或使用 `utf8.DecodeRuneInString()` 进行安全的字符访问

#### 2. 字符串构建
- 使用 `strings.Builder` 代替字符串拼接，提高性能
- 避免频繁的字符串分配

#### 3. 错误处理
Go 的错误处理模式：
```go
func JSONRepair(text string) (string, error) {
    parser := NewParser(text)
    result, err := parser.parse()
    if err != nil {
        return "", err
    }
    return result, nil
}
```

#### 4. 正则表达式
将 TypeScript 正则转换为 Go：
```go
// TypeScript: /^[0-9A-Fa-f]$/
// Go:
var hexRegex = regexp.MustCompile(`^[0-9A-Fa-f]$`)

// TypeScript: /^[[{\w-]$/
// Go:
var startOfValueRegex = regexp.MustCompile(`^[[\{\w-]$`)
```

#### 5. 字符比较
Go 中字符索引访问：
```go
// 获取字节
char := text[i]  // type: byte

// 获取 rune（Unicode 字符）
r, size := utf8.DecodeRuneInString(text[i:])
```

### 测试策略

#### 1. 单元测试
为每个解析函数编写独立测试：
- 测试对象解析
- 测试数组解析
- 测试字符串解析（各种引号情况）
- 测试数字解析
- 测试关键字解析

#### 2. 集成测试
测试完整的修复场景：
```go
func TestJSONRepair(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"missing quotes", "{name: 'John'}", `{"name": "John"}`},
        {"trailing comma", "[1,2,3,]", "[1,2,3]"},
        // ... 更多测试用例
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := JSONRepair(tt.input)
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if result != tt.expected {
                t.Errorf("expected %q, got %q", tt.expected, result)
            }
        })
    }
}
```

#### 3. 基准测试
```go
func BenchmarkJSONRepair(b *testing.B) {
    input := "{name: 'John', age: 30, city: 'New York'}"
    for i := 0; i < b.N; i++ {
        JSONRepair(input)
    }
}
```

### 性能优化建议

1. **避免字符串频繁分配**
   - 使用 `strings.Builder` 而非 `+` 操作符
   - 预分配合适的容量

2. **缓存编译后的正则表达式**
   - 使用包级变量存储 `*regexp.Regexp`

3. **减少函数调用开销**
   - 内联简单的工具函数（Go 编译器会自动优化）

4. **使用字节级操作**
   - 对于 ASCII 字符，直接使用 `byte` 而非 `rune`

## API 设计

### 基础 API

```go
package jsonrepair

// JSONRepair 修复无效的 JSON 字符串
// 如果无法修复，返回错误
func JSONRepair(text string) (string, error)

// MustJSONRepair 修复无效的 JSON 字符串
// 如果无法修复，panic
func MustJSONRepair(text string) string
```

### 选项配置（可选扩展）

```go
// Options 修复选项
type Options struct {
    // 是否保留注释（默认移除）
    PreserveComments bool
    // 是否处理 NDJSON
    HandleNDJSON bool
}

// JSONRepairWithOptions 使用自定义选项修复 JSON
func JSONRepairWithOptions(text string, opts Options) (string, error)
```

## 兼容性要求

- **Go 版本**: 1.18 - 1.23
- **依赖**: 仅使用 Go 标准库
- **平台**: 跨平台支持（Linux, macOS, Windows）

## 测试覆盖率目标

- 单元测试覆盖率: ≥ 90%
- 集成测试: 覆盖所有主要功能场景
- 边界测试: 包括空字符串、超大字符串等

## 文档要求

1. **README.md**
   - 项目介绍
   - 安装说明
   - 快速开始
   - API 文档
   - 示例代码

2. **代码注释**
   - 所有导出的函数和类型必须有文档注释
   - 复杂逻辑添加行内注释

3. **示例代码**
   - 基本使用示例
   - 高级功能示例
   - 错误处理示例

## 实现阶段

### Phase 1: 核心框架（优先）
1. ✅ 创建项目结构
2. ✅ 定义核心类型和接口
3. ✅ 实现字符串工具函数
4. ✅ 实现 Parser 结构体基础

### Phase 2: 基础解析（核心）
1. ✅ 实现 `parseValue()`
2. ✅ 实现 `parseObject()`
3. ✅ 实现 `parseArray()`
4. ✅ 实现 `parseString()`（基础版本）
5. ✅ 实现 `parseNumber()`
6. ✅ 实现 `parseKeywords()`

### Phase 3: 高级功能（重要）
1. ✅ 完善 `parseString()`（处理各种引号和转义）
2. ✅ 实现注释处理
3. ✅ 实现 Markdown 代码块处理
4. ✅ 实现 NDJSON 处理
5. ✅ 实现 MongoDB 数据类型处理
6. ✅ 实现 JSONP 处理

### Phase 4: 错误处理和修复策略（关键）
1. ✅ 实现缺失逗号修复
2. ✅ 实现缺失引号修复
3. ✅ 实现缺失括号修复
4. ✅ 实现尾随逗号修复
5. ✅ 实现截断 JSON 修复

### Phase 5: 测试和优化
1. ✅ 编写单元测试
2. ✅ 编写集成测试
3. ✅ 性能测试和优化
4. ✅ 边界测试

### Phase 6: 文档和发布
1. ✅ 编写完整文档
2. ✅ 创建示例代码
3. ✅ 准备发布到 GitHub
4. ✅ （可选）发布到 pkg.go.dev

## 注意事项

### 关键难点

1. **字符串解析的双阶段策略**
   - 需要仔细实现回溯机制
   - 正确处理各种引号类型

2. **空白字符处理**
   - 在修复时保持必要的空白
   - 正确插入和删除空白

3. **Unicode 处理**
   - Go 的字符串索引是基于字节的
   - 需要使用 `rune` 正确处理多字节字符

4. **递归解析**
   - 防止栈溢出（对于超深嵌套的 JSON）
   - 正确维护解析状态

### 与原项目的差异

1. **错误处理**
   - TypeScript 使用异常（throw）
   - Go 使用返回值 error

2. **字符串不可变性**
   - Go 字符串不可变，需要使用 `strings.Builder`

3. **字符类型**
   - TypeScript: `string[i]` 返回字符串
   - Go: `text[i]` 返回字节 (byte)，需要使用 `rune` 处理 Unicode

## 许可证

保持与原项目相同的 ISC 许可证，并在文件头部注明：
```
// Package jsonrepair provides functionality to repair invalid JSON documents.
// This is a Go port of the original jsonrepair project by Jos de Jong.
// Original project: https://github.com/josdejong/jsonrepair
// Licensed under the ISC License
```

## 参考资源

- 原项目仓库: https://github.com/josdejong/jsonrepair
- JSON 规范: https://www.json.org/
- RFC 8259: https://tools.ietf.org/html/rfc8259
- Go 编码规范: https://golang.org/doc/effective_go
- Go 字符串处理: https://blog.golang.org/strings

## 总结

本移植项目的目标是创建一个功能完整、性能优良、符合 Go 语言习惯的 JSON 修复库。通过严格遵循原项目的实现逻辑，同时充分利用 Go 语言的特性，我们将创建一个高质量的 Go 版本 jsonrepair 库。
