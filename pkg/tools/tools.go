package tools

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
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

// FileWithModTime 表示文件及其最后修改时间
type FileWithModTime struct {
	Path    string
	ModTime time.Time
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

// MergeStringFlags 函数用于合并长选项和短选项，返回一个字符串和一个错误。
// 如果同时指定了长选项和短选项，则返回一个错误。
// 如果只指定了一个选项，则返回该选项的值。
// 如果两个选项都未指定，则返回空字符串。
func MergeStringFlags(longPtr, shortPtr string) (string, error) {
	// 检查长选项和短选项是否同时不为空
	if longPtr != "" && shortPtr != "" {
		// 如果同时指定了长选项和短选项，则返回一个错误
		return "", errors.New("不能同时指定长选项和短选项")
	}

	// 如果长选项不为空, 则返回长选项
	if longPtr != "" {
		return longPtr, nil
	}

	// 如果短选项不为空, 则返回短选项
	if shortPtr != "" {
		return shortPtr, nil
	}

	// 如果长选项和短选项都为空, 则返回空字符串
	return "", nil
}

// MoveDir 将一个目录移动到另一个位置
// 参数：
//
//	src - 源目录路径
//	dst - 目标目录路径
//
// 返回：
//
//	error - 如果发生错误，返回错误信息；否则返回nil
func MoveDir(src, dst string) error {
	// 检查源目录是否存在
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("源目录 %s 不存在", src)
	}

	// 遍历源目录
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算目标路径
		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relativePath)

		// 如果是目录，创建对应的目录
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 如果是文件，移动文件
		return moveFile(path, dstPath)
	})
	if err != nil {
		return err
	}

	// 删除源目录
	return os.RemoveAll(src)
}

// moveFile 移动单个文件
func moveFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 拷贝内容
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// 获取源文件的 FileInfo
	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// 设置目标文件的权限
	if err := os.Chmod(dst, srcFileInfo.Mode()); err != nil {
		return err
	}

	// 删除源文件
	return os.Remove(src)
}

// GetFileMD5Last8 获取文件的 MD5 哈希值的后 8 位
// 参数：
//
//	filePath - 文件路径
//
// 返回值：
//
//	string - 文件 MD5 哈希值的后 8 位
//	error - 如果发生错误，返回错误信息；否则返回 nil
func GetFileMD5Last8(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件时出错: %w", err)
	}
	defer file.Close()

	// 创建 MD5 哈希对象
	hash := md5.New()

	// 分块读取文件内容
	buffer := make([]byte, 32*1024) // 32KB 缓冲区
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("读取文件内容时出错: %w", err)
		}
		if n == 0 {
			break // 文件读取完成
		}
		hash.Write(buffer[:n]) // 将读取的内容写入哈希对象
	}

	// 获取完整的 MD5 哈希值
	sum := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashStr := fmt.Sprintf("%x", sum)

	// 返回哈希值的后 8 位
	return hashStr[len(hashStr)-8:], nil
}

// HumanReadableSize 获取文件大小并转换为人性化单位显示
// 参数：
//
//	filePath - 文件路径
//
// 返回值：
//
//	string - 文件大小的人性化表示
//	error - 如果发生错误，返回错误信息；否则返回 nil
func HumanReadableSize(filePath string) (string, error) {
	// 打开文件
	file, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("获取文件信息时出错: %w", err)
	}

	// 获取文件大小（以字节为单位）
	size := file.Size()

	// 定义单位和换算关系
	units := []string{"B", "KB", "MB", "GB"}
	base := float64(1024)

	// 转换为合适的单位
	var unit string
	sizeFloat := float64(size) // 将 size 转换为 float64
	if sizeFloat < base {
		unit = units[0]
	} else if sizeFloat < base*base {
		unit = units[1]
		sizeFloat /= base
	} else if sizeFloat < base*base*base {
		unit = units[2]
		sizeFloat /= base * base
	} else {
		unit = units[3]
		sizeFloat /= base * base * base
	}

	// 格式化输出
	return fmt.Sprintf("%.2f%s", sizeFloat, unit), nil
}

// GetZipFiles 获取指定目录下所有以 .zip 结尾的文件列表
// 参数：
//
//	dirPath - 目录路径
//
// 返回值：
//
//	[]string - 匹配的文件路径列表
//	error - 如果发生错误，返回错误信息；否则返回 nil
func GetZipFiles(dirPath string) ([]string, error) {
	// 用于存储匹配的文件路径
	var zipFiles []string

	// 遍历目录及其子目录
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("遍历文件时出错: %w", err)
		}

		// 检查文件是否以 .zip 结尾
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zip") {
			zipFiles = append(zipFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录时出错: %w", err)
	}

	return zipFiles, nil
}

// SortFilesByModTime 按照文件的最后修改时间对文件列表进行排序
func SortFilesByModTime(files []FileWithModTime) []FileWithModTime {
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.Before(files[j].ModTime)
	})
	return files
}

// RetainLatestFiles 处理文件列表，仅保留最新的指定数量的文件
// 参数：
//
//	files - 文件路径列表
//	retainCount - 保留的文件数量
//
// 返回值：
//
//	error - 如果发生错误，返回错误信息；否则返回 nil
func RetainLatestFiles(files []string, retainCount int) error {
	// 用于存储文件及其最后修改时间
	var fileInfos []FileWithModTime

	// 获取每个文件的最后修改时间
	for _, filePath := range files {
		info, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			CL.PrintErrorf("文件不存在，跳过: %s", filePath)
			continue
		} else if err != nil {
			return fmt.Errorf("获取文件信息时出错: %w", err)
		}
		fileInfos = append(fileInfos, FileWithModTime{Path: filePath, ModTime: info.ModTime()})
	}

	// 按照文件的最后修改时间排序
	fileInfos = SortFilesByModTime(fileInfos)

	// 如果文件数量小于或等于保留数量，直接返回
	if len(fileInfos) <= retainCount {
		CL.PrintSuccessf("文件数量小于或等于保留数量，无需清理")
		return nil
	}

	// 删除最早的文件，仅保留最新的指定数量的文件
	for i := 0; i < len(fileInfos)-retainCount; i++ {
		filePath := fileInfos[i].Path
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("清理文件时出错: %w", err)
		}
		CL.PrintSuccessf("清理历史文件: %s", filePath)
	}

	return nil
}

// CompressFilesByOS 函数根据操作系统类型执行不同的压缩命令
// 参数：
//
//	targetDir - 目标目录路径
//	targetName - 要压缩的目标名称
//	backupFileNamePath - 备份文件的名称
//
// 返回值：
//
//	string - 压缩文件的路径
//	error - 如果发生错误，返回错误信息；否则返回 nil
func CompressFilesByOS(targetDir, targetName, backupFileNamePath string) (string, error) {
	// 检查操作系统类型
	if runtime.GOOS == "linux" {
		// 构建完整的压缩文件路径
		backupFilePath := filepath.Join(targetDir, fmt.Sprintf("%s.tgz", backupFileNamePath))

		// 检查tar命令是否可用
		if _, err := exec.LookPath("tar"); err != nil {
			return "", fmt.Errorf("tar 命令不可用: %w", err)
		}

		// 执行 tar 命令进行压缩
		cmd := exec.Command("tar", "-czf", backupFilePath, targetName)
		cmd.Dir = targetDir // 设置工作目录为目标目录
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("压缩文件时出错: %w", err)
		}

		return backupFilePath, nil
	}

	if runtime.GOOS == "windows" {
		// 构建完整的压缩文件路径
		backupFilePath := filepath.Join(targetDir, fmt.Sprintf("%s.zip", backupFileNamePath))

		// 检查7z命令是否可用
		if _, err := exec.LookPath("7z"); err != nil {
			return "", fmt.Errorf("7z 命令不可用: %w", err)
		}

		// 执行 7z 命令进行压缩
		cmd := exec.Command("7z", "a", "-tzip", backupFilePath, targetName)
		cmd.Dir = targetDir // 设置工作目录为目标目录
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("压缩文件时出错: %w", err)
		}

		return backupFilePath, nil
	}

	// 如果操作系统类型不支持，返回错误
	return "", fmt.Errorf("不支持的操作系统类型")

}
