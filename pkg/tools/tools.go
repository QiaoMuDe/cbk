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
	"github.com/jmoiron/sqlx"
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
func GetZipFiles(dirPath, FileExtension string) ([]string, error) {
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("在获取文件列表时，目录不存在: %s", dirPath)
	}

	// 检查文件扩展名是否为空
	if FileExtension == "" {
		return nil, fmt.Errorf("在获取文件列表时，文件扩展名不能为空")
	}

	// 用于存储匹配的文件路径
	var zipFiles []string

	// 遍历目录及其子目录
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("遍历文件时出错: %w", err)
		}

		// 检查文件是否以 指定扩展名 结尾
		if !info.IsDir() && strings.HasSuffix(info.Name(), FileExtension) {
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
//	db - 数据库连接
//	files - 文件路径列表
//	retainCount - 保留的文件数量
//
// 返回值：
//
//	error - 如果发生错误，返回错误信息；否则返回 nil
func RetainLatestFiles(db *sqlx.DB, files []string, retainCount int) error {
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

		// 按小数点分割文件名，获取文件名称部分
		baseName := filepath.Base(filePath) // 获取文件名（带扩展名）
		parts := strings.Split(baseName, ".")

		// 按下划线分割文件名称部分，获取文件名称和扩展名部分
		nameParts := strings.Split(parts[0], "_")

		// 构建更新sql
		updateSql := `update backup_records set data_status =? where task_name =? and timestamp =?;`

		// 执行更新sql
		if _, err := db.Exec(updateSql, "0", nameParts[0], nameParts[1]); err != nil {
			return fmt.Errorf("更新备份记录时出错: %w", err)
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
func CompressFilesByOS(db *sqlx.DB, targetDir, targetName, backupFileNamePath string) (string, error) {
	// 检查操作系统类型
	if runtime.GOOS != "linux" && runtime.GOOS != "windows" {
		return "", fmt.Errorf("不支持的操作系统类型: %s", runtime.GOOS)
	}

	// 查询sql
	qurySql := `SELECT os_type, compress_tool, compress_args, file_extension FROM compress_config where os_type = ?;`
	var compressConfig struct {
		OsType        string `db:"os_type"`        // 操作系统类型
		CompressTool  string `db:"compress_tool"`  // 压缩工具名称
		CompressArgs  string `db:"compress_args"`  // 压缩工具的参数
		FileExtension string `db:"file_extension"` // 文件扩展名
	}

	// 查询压缩配置
	err := db.Get(&compressConfig, qurySql, runtime.GOOS)
	if err != nil {
		return "", fmt.Errorf("查询压缩配置时出错: %w", err)
	} else if compressConfig.CompressTool == "" {
		return "", fmt.Errorf("未找到对应的压缩配置")
	}

	// // 打印压缩配置
	// CL.PrintInfof("压缩配置: %+v", compressConfig)

	// 构建完整的压缩文件路径
	backupFilePath := fmt.Sprintf("%s%s", backupFileNamePath, compressConfig.FileExtension)

	// 检查tar命令是否可用
	if _, err := exec.LookPath(compressConfig.CompressTool); err != nil {
		return "", fmt.Errorf("%s 命令不可用: %w", compressConfig.CompressTool, err)
	}

	// 拆分 CompressArgs 为独立的参数
	args := strings.Fields(compressConfig.CompressArgs)
	args = append(args, backupFilePath, targetName)

	// 执行命令进行压缩
	cmd := exec.Command(compressConfig.CompressTool, args...)
	cmd.Dir = targetDir // 设置工作目录为目标目录
	// 设置标准输出和标准错误输出
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out

	// // 打印执行的命令
	// CL.PrintInfof("执行命令: %s", cmd.String())

	// // 打印运行的工作目录
	// CL.PrintInfof("运行的工作目录: %s", cmd.Dir)

	if err := cmd.Run(); err != nil {
		CL.PrintErrorf("执行打包命令时出错: %s", out.String())
		return "", fmt.Errorf("压缩文件时出错: %w", err)
	}

	return backupFilePath, nil
}

// UncompressFilesByOS 函数根据操作系统类型执行不同的解压命令
// 参数：
//
//	 db - 数据库连接
//		zipDir - 压缩文件所在目录
//		zipFileName - 压缩文件名
//		outputPath - 解压输出路径
//
// 返回值：
//
//	string - 解压文件的路径
//	error - 如果发生错误，返回错误信息；否则返回 nil
func UncompressFilesByOS(db *sqlx.DB, zipDir, zipFileName, outputPath string) (string, error) {
	// 检查操作系统类型
	if runtime.GOOS != "linux" && runtime.GOOS != "windows" {
		return "", fmt.Errorf("不支持的操作系统类型: %s", runtime.GOOS)
	}

	// 检查解压输出路径是否存在
	if _, err := CheckPath(outputPath); err != nil {
		return "", fmt.Errorf("解压输出路径不存在: %w", err)
	}

	// 查询sql
	querySql := `select decompress_tool, decompress_args, file_extension from decompress_config where os_type = ?;`
	var decompressConfig struct {
		DecompressTool string `db:"decompress_tool"` // 解压工具名称
		DecompressArgs string `db:"decompress_args"` // 解压工具的参数
		FileExtension  string `db:"file_extension"`  // 文件扩展名
	}

	// 查询解压配置
	err := db.Get(&decompressConfig, querySql, runtime.GOOS)
	if err != nil {
		return "", fmt.Errorf("查询解压配置时出错: %w", err)
	} else if decompressConfig.DecompressTool == "" {
		return "", fmt.Errorf("未找到对应的解压配置")
	}

	// 检查解压工具是否可用
	if _, err := exec.LookPath(decompressConfig.DecompressTool); err != nil {
		return "", fmt.Errorf("%s 命令不可用: %w", decompressConfig.DecompressTool, err)
	}

	// 获取解压缩文件的完整路径
	zipFilePath := filepath.Join(zipDir, zipFileName)

	// 检查解压缩文件是否存在
	if _, err := CheckPath(zipFilePath); err != nil {
		return "", fmt.Errorf("解压文件不存在: %w", err)
	}

	// 检查输出路径下是否存在同名
	baseName := strings.Split(zipFileName, "_")[0] // 按下划线分割文件名，获取文件名称部分
	outputPath = filepath.Join(outputPath, baseName)
	if info, err := CheckPath(outputPath); os.IsNotExist(err) {
		// 不存在则忽略
	} else if info.Exists {
		return "", fmt.Errorf("解压文件路径已存在: %s", outputPath)
	}

	// 拆分 DecompressArgs 为独立的参数
	args := strings.Split(decompressConfig.DecompressArgs, "|")
	// 构建完整的解压命令, 例如1: tar -xvf test.tar.gz -C /home/test，例如2: 7z x -tzip test.zip -o /home/test
	args = append(args, args[0], zipFilePath, args[1], outputPath)

	// 执行命令进行解压
	cmd := exec.Command(decompressConfig.DecompressTool, args...)
	//cmd.Dir = zipDir // 设置工作目录为目标目录
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out

	// 执行命令
	if err := cmd.Run(); err != nil {
		CL.PrintErrorf("执行解压命令时出错: %s", out.String())
		return "", fmt.Errorf("解压文件时出错: %w", err)
	}

	return outputPath, nil
}
