package tools

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
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

// PathInfo 是一个结构体，用于封装路径的信息
type PathInfo struct {
	Path    string      // 路径
	Exists  bool        // 是否存在
	IsFile  bool        // 是否为文件
	IsDir   bool        // 是否为目录
	Size    int64       // 文件大小（字节）
	Mode    os.FileMode // 文件权限
	ModTime time.Time   // 文件修改时间
}

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

// CheckPath 检查给定路径的信息
// path: 要检查的路径
func CheckPath(path string) (PathInfo, error) {
	// 创建一个 PathInfo 结构体
	var info PathInfo

	// 清理路径，确保没有多余的斜杠
	path = filepath.Clean(path)

	// 设置路径
	info.Path = path

	// 使用 os.Stat 获取文件状态
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果路径不存在, 则直接返回
			info.Exists = false
			return info, fmt.Errorf("路径 '%s' 不存在，请检查路径是否正确: %s", path, err)
		} else {
			return info, fmt.Errorf("无法访问路径 '%s': %s", path, err)
		}
	}

	// 路径存在，填充信息
	info.Exists = true                // 标记路径存在
	info.IsFile = !fileInfo.IsDir()   // 通过取反判断是否为文件，因为 IsDir 返回 false 表示是文件
	info.IsDir = fileInfo.IsDir()     // 直接使用 IsDir 方法判断是否为目录
	info.Size = fileInfo.Size()       // 获取文件大小
	info.Mode = fileInfo.Mode()       // 获取文件权限
	info.ModTime = fileInfo.ModTime() // 获取文件的最后修改时间

	// 返回路径信息结构体
	return info, nil
}
