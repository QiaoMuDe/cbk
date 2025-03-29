package tools

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gitee.com/MM-Q/colorlib"
)

// 定义随机字符集
const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// 全局随机数生成器
var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// 全局颜色渲染器
var CL = colorlib.NewColorLib()

// GenerateID 生成包含时间戳和指定长度随机字符的ID
// randomLength: 随机部分的长度
func GenerateID(randomLength int) string {
	if randomLength < 0 {
		return "" // 如果随机部分长度为负数，返回空字符串
	}

	// 获取当前时间的纳秒级时间戳
	timestamp := time.Now().UnixNano()

	// 使用strings.Builder进行字符串拼接
	var builder strings.Builder

	// 将时间戳转换为字符串并拼接到builder中
	builder.WriteString(fmt.Sprintf("%d", timestamp))

	// 预先分配足够的内存
	randomPart := make([]byte, randomLength)
	for i := 0; i < randomLength; i++ {
		randomPart[i] = charset[globalRand.Intn(len(charset))]
	}

	// 将随机字符部分拼接到builder中
	builder.Write(randomPart)

	// 返回最终拼接的字符串
	return builder.String()
}
