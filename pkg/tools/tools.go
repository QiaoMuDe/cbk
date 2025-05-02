package tools

import (
	"archive/zip"
	"cbk/pkg/globals"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gitee.com/MM-Q/colorlib"
	"github.com/jmoiron/sqlx"
	"github.com/schollz/progressbar/v3"
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
	builder.Grow(8 + randomLength) // 预分配内存

	// 获取时间戳的后8位并转换为字符串
	timestampStr := strconv.FormatInt(timestamp%1e8, 10)

	// 将时间戳转换为字符串并拼接到builder中
	builder.WriteString(timestampStr)

	// 生成随机字符部分
	for i := 0; i < randomLength; i++ {
		builder.WriteByte(charset[globalRand.Intn(len(charset))])
	}

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
			return info, fmt.Errorf("路径 '%s' 不存在，请检查路径是否正确: %s", path, os.ErrNotExist)
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

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("获取文件大小时出错: %w", err)
	}
	fileSize := fileInfo.Size()

	// 创建进度条
	bar := progressbar.DefaultBytes(
		fileSize,
		"正在计算MD5",
	)

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
		// 更新进度条// progressbar 库要求传入 int64 类型
		if err := bar.Add64(int64(n)); err != nil {
			return "", fmt.Errorf("更新进度条失败: %w", err)
		}
	}

	// 获取完整的 MD5 哈希值
	sum := hash.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashStr := fmt.Sprintf("%x", sum)

	// 确保进度条完成
	if err := bar.Finish(); err != nil {
		return "", fmt.Errorf("进度条完成失败: %w", err)
	}

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
// ascending: true 表示旧文件在前(升序)，false 表示新文件在前(降序)
func SortFilesByModTime(files []FileWithModTime, ascending bool) []FileWithModTime {
	// 检查文件列表是否为空
	if len(files) == 0 {
		return files
	}

	// 对文件列表按照最后修改时间进行排序
	sort.Slice(files, func(i, j int) bool {
		if ascending {
			return files[i].ModTime.Before(files[j].ModTime) // 升序: 旧文件在前
		}
		return files[i].ModTime.After(files[j].ModTime) // 降序: 新文件在前
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
func RetainLatestFiles(db *sqlx.DB, files []string, retainCount, retainDays int) error {
	// 用于存储文件及其最后修改时间
	var fileInfos []FileWithModTime

	// 检查获取到的文件列表是否为空
	if len(files) == 0 {
		return nil
	}

	// 获取每个文件的最后修改时间
	for _, filePath := range files {
		// 获取文件信息
		info, err := os.Stat(filePath)
		// 检查文件是否存在
		if os.IsNotExist(err) {
			CL.PrintErrf("备份文件不存在, 跳过: %s", filePath)
			continue
		} else if err != nil {
			CL.PrintErrf("获取文件信息时出错: %v", err)
			continue
		}

		// 将文件路径和最后修改时间存储到 fileInfos 切片中
		fileInfos = append(fileInfos, FileWithModTime{Path: filePath, ModTime: info.ModTime()})
	}

	// 按照文件的最后修改时间排序, 旧的在前面
	fileInfos = SortFilesByModTime(fileInfos, true)

	// 根据保留天数和保留数量筛选文件
	filesToDelete := filterFiles(fileInfos, retainDays, retainCount)

	// 按照筛选后的文件列表删除文件
	if err := deleteFiles(db, filesToDelete); err != nil {
		return fmt.Errorf("清理文件时出错: %w", err)
	}

	return nil
}

// filterFiles 根据保留天数和保留数量筛选文件
// 参数：
//
//	files - 文件列表
//	retainDays - 保留的天数
//	retainCount - 保留的文件数量
//
// 返回值：
//
//	[]FileWithModTime - 筛选后的文件列表
func filterFiles(files []FileWithModTime, retainDays, retainCount int) []FileWithModTime {
	var filesToDelete []FileWithModTime

	// 如果为空，直接返回
	if len(files) == 0 {
		return filesToDelete
	}

	// 仅设置保留数量，没有设置保留天数
	if retainDays == 0 && retainCount > 0 {
		// 保留最新的retainCount个文件, 如果文件数量大于retainCount, 则将从0到len(files)-retainCount的文件添加到filesToDelete中
		if len(files) > retainCount {
			filesToDelete = files[0 : len(files)-retainCount] // 从0到len(files)-retainCount的文件添加到filesToDelete中
		}

		return filesToDelete
	}

	// 同时设置保留天数和保留数量
	if retainDays > 0 && retainCount > 0 {
		// 先按天数筛选分组，然后在每天分组中保留最新的retainCount个文件
		retainTime := time.Now().AddDate(0, 0, -retainDays)
		dayMap := make(map[time.Time][]FileWithModTime)

		// 按天分组
		for _, file := range files {
			// 检查文件的最后修改时间是否在保留时间范围内
			if file.ModTime.After(retainTime) {
				// 将文件按天分组
				day := time.Date(file.ModTime.Year(), file.ModTime.Month(), file.ModTime.Day(), 0, 0, 0, 0, file.ModTime.Location())
				dayMap[day] = append(dayMap[day], file)
			}
		}

		// 在每天分组中保留最新的retainCount个文件
		for _, dayFiles := range dayMap {
			// 检查每天分组的文件数量是否大于保留数量
			if len(dayFiles) > retainCount {
				// 如果大于保留数量，则将从0到len(dayFiles)-retainCount的文件添加到filesToDelete中
				filesToDelete = append(filesToDelete, dayFiles[0:len(dayFiles)-retainCount]...)
			}
		}

		return filesToDelete
	}

	// 返回空切片
	return []FileWithModTime{}
}

// deleteFiles 根据文件列表删除文件
// 参数：
//
//	db - 数据库连接
//	files - 文件列表
//
// 返回值：
//
// error - 如果发生错误，返回错误信息；否则返回 nil
func deleteFiles(db *sqlx.DB, files []FileWithModTime) error {
	// 构建删除sql
	deleteSql := `delete from backup_records where task_name = ? and timestamp = ?;`

	// 遍历文件列表, 删除文件
	for _, file := range files {
		// 按小数点分割文件名，获取文件名称部分
		baseName := filepath.Base(file.Path) // 获取文件名（带扩展名）
		parts := strings.Split(baseName, ".")

		// 按下划线分割文件名称部分，获取文件名称和扩展名部分
		nameParts := strings.Split(parts[0], "_")

		// 检查 nameParts 的长度是否足够
		if len(nameParts) < 2 {
			CL.PrintErrf("文件名格式不正确，无法解析出足够的部分: %s, 请在稍后手动删除", file.Path)
			continue
		} else {
			// 检查文件是否存在, 如果存在, 则删除
			if _, err := CheckPath(file.Path); err == nil {
				if err := os.Remove(file.Path); err != nil {
					CL.PrintErrf("清理 %s 文件时出错: %v, 请在稍后手动删除", file.Path, err)
				}
			}
		}

		taskName := nameParts[0]  // 任务名称
		timestamp := nameParts[1] // 时间戳

		// 执行删除操作sql, 参数: 任务名称, 时间戳
		if _, err := db.Exec(deleteSql, taskName, timestamp); err != nil {
			return fmt.Errorf("删除备份记录时出错: %w", err)
		}
	}

	return nil
}

// CreateZipFromOSPaths 根据目标目录和文件名创建ZIP压缩文件
// 参数:
//
//	db - 数据库连接(当前未使用，保留参数)
//	targetDir - 需要压缩的目标目录路径
//	targetName - 需要压缩的目标名称(文件或目录名)
//	backupFileNamePath - 备份文件的基础路径(不含扩展名)
//	noCompression - 是否禁用压缩，默认为false
//	excludeFunc - 排除函数，用于决定是否跳过文件或目录
//
// 返回值:
//
//	string - 生成的ZIP文件完整路径
//	error - 操作过程中遇到的错误
func CreateZipFromOSPaths(db *sqlx.DB, targetDir, targetName, backupFileNamePath string, noCompression int, filter globals.ExcludeFunc) (string, error) {
	// 构建完整的压缩文件路径(添加扩展名)
	zipFilePath := fmt.Sprintf("%s%s", backupFileNamePath, ".zip")

	// 切换到目标目录以便后续操作
	if err := os.Chdir(targetDir); err != nil {
		return "", fmt.Errorf("切换到目标目录时出错: %w", err)
	}

	// 调用CreateZip函数执行实际压缩操作
	if err := CreateZip(zipFilePath, targetName, noCompression, filter); err != nil {
		return "", fmt.Errorf("压缩文件时出错: %w", err)
	}

	// 返回生成的ZIP文件完整路径
	return zipFilePath, nil
}

// UncompressFilesByOS 根据目标目录和文件名解压ZIP压缩文件
// 参数:
//
// zipDir - 需要解压的ZIP文件所在目录路径
//
//	zipFileName - 需要解压的ZIP文件名
//	outputPath - 解压后的文件存放路径
//
// 返回值:
//
//	string - 解压后的文件存放路径
//	error - 操作过程中遇到的错误
func UncompressFilesByOS(zipDir, zipFileName, outputPath string) (string, error) {
	// 检查解压输出路径是否存在
	if _, err := CheckPath(outputPath); err != nil {
		return "", fmt.Errorf("解压输出路径不存在: %w", err)
	}

	// 获取解压缩文件的完整路径, 例如: /home/backup/zip/20240506_123456.zip
	zipFilePath := filepath.Join(zipDir, zipFileName)

	// 检查解压缩文件是否存在
	if _, err := CheckPath(zipFilePath); err != nil {
		return "", fmt.Errorf("解压文件不存在: %w", err)
	}

	// 检查输出路径下是否存在同名
	baseName := strings.Split(zipFileName, "_")[0]  // 按下划线分割文件名，获取文件名称部分
	tempPath := filepath.Join(outputPath, baseName) // 构建临时路径
	if _, err := CheckPath(tempPath); err == nil {
		return "", fmt.Errorf("解压输出路径下存在同名: %s", tempPath)
	}

	// 调用解压函数
	if err := Unzip(zipFilePath, outputPath); err != nil {
		return "", fmt.Errorf("解压文件时出错: %w", err)
	}

	return outputPath, nil
}

// CreateZip 函数用于创建ZIP压缩文件
// 参数:
//
//	zipFilePath - 生成的ZIP文件路径
//	sourceDir - 需要压缩的源目录路径
//	noCompression - 是否禁用压缩，默认为false
//	excludeFunc - 排除函数，用于决定是否跳过文件或目录
//
// 返回值:
//
//	error - 操作过程中遇到的错误
func CreateZip(zipFilePath string, sourceDir string, noCompression int, excludeFunc globals.ExcludeFunc) error {
	// 检查zipFilePath是否为绝对路径，如果不是，将其转换为绝对路径
	if !filepath.IsAbs(zipFilePath) {
		absPath, err := filepath.Abs(zipFilePath)
		if err != nil {
			return fmt.Errorf("转换zipFilePath为绝对路径失败: %w", err)
		}
		zipFilePath = absPath
	}

	// 检查sourceDir是否为绝对路径，如果不是，将其转换为绝对路径
	if !filepath.IsAbs(sourceDir) {
		absPath, err := filepath.Abs(sourceDir)
		if err != nil {
			return fmt.Errorf("转换sourceDir为绝对路径失败: %w", err)
		}
		sourceDir = absPath
	}

	// 如果没有提供排除函数，使用默认的排除函数
	if excludeFunc == nil {
		excludeFunc = func(path string, info os.FileInfo) bool {
			return false // 默认不跳过任何文件
		}
	}

	// 获取是否压缩的标志
	var deflateOrStore uint16
	// 根据noCompression参数设置压缩方法, 默认使用Deflate压缩
	if noCompression == 1 {
		deflateOrStore = zip.Store // 等于1时，表示禁用压缩
	} else {
		deflateOrStore = zip.Deflate // 等于0时，表示启用压缩
	}

	// 创建 ZIP 文件
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("创建 ZIP 文件失败: %w", err)
	}
	defer zipFile.Close()

	// 创建 ZIP 写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close() // 确保在函数结束时关闭 ZIP 写入器

	// 获取源目录的总大小，用于进度条
	totalSize := int64(0)

	// 创建一个不确定进度的进度条
	iBar := progressbar.DefaultBytes(
		-1, // 设置总大小为 -1，表示不确定进度
		"正在计算大小...",
	)

	// 遍历目录并计算总大小
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("遍历目录时出错: %w", err)
		}

		// 检查是否需要跳过当前文件或目录
		if excludeFunc(path, info) {
			if info.IsDir() {
				// 如果是目录，跳过其所有子文件和子目录
				return filepath.SkipDir
			}
			// 如果是文件，直接跳过
			return nil
		}

		// 跳过目录本身
		if !info.IsDir() {
			// 获取文件的详细状态
			fileStat, err := os.Lstat(path)
			if err != nil {
				return fmt.Errorf("获取文件状态失败: %w", err)
			}

			// 根据文件类型处理
			switch mode := fileStat.Mode(); {
			case mode.IsRegular():
				// 普通文件
				totalSize += info.Size()

				// 更新进度条
				if err := iBar.Add64(info.Size()); err != nil {
					return fmt.Errorf("更新进度条失败: %w", err)
				}
			case mode.IsDir():
				// 目录，跳过
				return nil
			case mode&os.ModeSymlink != 0:
				// 软链接，跳过
				return nil
			case mode&os.ModeDevice != 0:
				// 设备文件，跳过
				return nil
			default:
				// 其他特殊文件类型，跳过
				return nil
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("获取源目录大小失败: %w", err)
	}

	// 关闭不确定进度的进度条
	if err := iBar.Finish(); err != nil {
		return fmt.Errorf("关闭进度条失败: %w", err)
	}

	// 初始化进度条
	bar := progressbar.DefaultBytes(
		totalSize,
		"正在打包",
	)

	// 遍历目录并添加文件到 ZIP 包
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("遍历目录时出错: %w", err)
		}

		// 检查是否需要跳过当前文件或目录
		if excludeFunc(path, info) {
			if info.IsDir() {
				// 如果是目录，跳过其所有子文件和子目录
				return filepath.SkipDir
			}
			// 如果是文件，直接跳过
			return nil
		}

		// 获取相对路径，保留顶层目录
		headerName, err := filepath.Rel(filepath.Dir(sourceDir), path)
		if err != nil {
			return fmt.Errorf("获取相对路径失败: %w", err)
		}

		// 替换路径分隔符为正斜杠（ZIP 文件格式要求）
		headerName = filepath.ToSlash(headerName)

		// 获取文件的详细状态
		fileStat, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("获取文件状态失败: %w", err)
		}

		// 根据文件类型处理
		switch mode := fileStat.Mode(); {
		case mode.IsRegular():
			// 普通文件
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return fmt.Errorf("创建 ZIP 文件头失败: %w", err)
			}
			// 设置文件头的名称
			header.Name = headerName

			// 设置压缩方法为 Deflate
			header.Method = deflateOrStore

			// 创建 ZIP 写入器
			fileWriter, err := zipWriter.CreateHeader(header)
			if err != nil {
				return fmt.Errorf("创建 ZIP 写入器失败: %w", err)
			}

			// 打开文件并写入 ZIP 包
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("打开文件失败: %w", err)
			}
			defer file.Close()

			// 创建一个自定义多路写入器，用于同时写入文件和进度条
			multiWriter := io.MultiWriter(fileWriter, bar)

			// 使用缓冲区进行文件复制，提高性能
			buffer := make([]byte, 512*1024) // 512KB 缓冲区大小
			if _, err := io.CopyBuffer(multiWriter, file, buffer); err != nil {
				return fmt.Errorf("写入 ZIP 文件失败: %w", err)
			}

		case mode.IsDir():
			// 目录
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return fmt.Errorf("创建 ZIP 文件头失败: %w", err)
			}
			// 设置目录的名称，末尾添加斜杠
			header.Name = headerName + "/"

			// 设置压缩方法为 Store（不压缩）
			header.Method = zip.Store

			// 创建目录
			if _, err := zipWriter.CreateHeader(header); err != nil {
				return fmt.Errorf("创建 ZIP 目录失败: %w", err)
			}

		case mode&os.ModeSymlink != 0:
			// 软链接
			target, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("读取软链接目标失败: %w", err)
			}

			// 创建软链接文件头
			header := &zip.FileHeader{
				Name:   headerName,
				Method: zip.Store,
			}
			// 设置软链接的元数据
			header.SetMode(mode)

			// 创建软链接
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return fmt.Errorf("创建 ZIP 软链接失败: %w", err)
			}
			if _, err := writer.Write([]byte(target)); err != nil {
				return fmt.Errorf("写入软链接目标失败: %w", err)
			}

		case mode&os.ModeDevice != 0:
			// 设备文件
			header := &zip.FileHeader{
				Name:   headerName,
				Method: zip.Store,
			}
			// 设置设备文件的元数据
			header.SetMode(mode)

			// 创建设备文件
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return fmt.Errorf("创建 ZIP 设备文件失败: %w", err)
			}
			// 设备文件通常不包含数据，只记录其元数据
			if _, err := writer.Write([]byte{}); err != nil {
				return fmt.Errorf("写入设备文件失败: %w", err)
			}

		default:
			// 其他特殊文件类型
			header := &zip.FileHeader{
				Name:   headerName,
				Method: zip.Store,
			}
			// 设置特殊文件的元数据
			header.SetMode(mode)

			// 创建特殊文件
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return fmt.Errorf("创建 ZIP 特殊文件失败: %w", err)
			}
			// 特殊文件通常不包含数据，只记录其元数据
			if _, err := writer.Write([]byte{}); err != nil {
				return fmt.Errorf("写入特殊文件失败: %w", err)
			}

		}

		return nil
	})

	// 检查是否有错误发生
	if err != nil {
		return fmt.Errorf("打包目录到 ZIP 失败: %w", err)
	}

	// 关闭进度条
	if err := bar.Finish(); err != nil {
		return fmt.Errorf("关闭进度条失败: %w", err)
	}

	return nil
}

// Unzip 解压缩 ZIP 文件到指定目录
// 参数:
//   - zipFilePath: 要解压缩的 ZIP 文件路径
//   - targetDir: 解压缩后的目标目录路径
//
// 返回值:
//   - error: 解压缩过程中发生的错误
func Unzip(zipFilePath string, targetDir string) error {
	// 打开 ZIP 文件
	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("打开 ZIP 文件失败: %w", err)
	}
	defer zipReader.Close()

	// 检查目标目录是否存在, 如果不存在, 则创建
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0644); err != nil {
			return fmt.Errorf("创建目标目录失败: %w", err)
		}
	}

	// 获取 ZIP 文件的总大小
	var totalSize uint64

	// 遍历 ZIP 文件中的每个文件或目录, 计算总大小
	for _, file := range zipReader.File {
		totalSize += file.UncompressedSize64 // 通过 UncompressedSize64 获取未压缩的文件大小
	}

	// 创建进度条
	bar := progressbar.DefaultBytes(
		int64(totalSize), // progressbar 库要求传入 int64 类型
		"正在解压",
	)

	// 遍历 ZIP 文件中的每个文件或目录
	for _, file := range zipReader.File {
		// 获取目标路径
		targetPath := filepath.Join(targetDir, file.Name)

		// 获取文件的模式
		mode := file.Mode()

		// 使用 switch 语句处理不同类型的文件
		switch {
		case mode.IsDir():
			// 目录
			if err := os.MkdirAll(targetPath, mode.Perm()); err != nil {
				return fmt.Errorf("创建目录失败: %w", err)
			}
		case mode&os.ModeSymlink != 0:
			// 软链接
			zipFileReader, err := file.Open()
			if err != nil {
				return fmt.Errorf("打开 ZIP 文件中的软链接失败: %w", err)
			}
			defer zipFileReader.Close()

			var target string
			// 读取软链接的目标路径
			if _, err := fmt.Fscanln(zipFileReader, &target); err != nil {
				return fmt.Errorf("读取软链接目标失败: %w", err)
			}

			// 检查软链接的父目录是否存在，如果不存在，则创建
			parentDir := filepath.Dir(targetPath)
			if _, err := os.Stat(parentDir); os.IsNotExist(err) {
				if err := os.MkdirAll(parentDir, 0644); err != nil {
					return fmt.Errorf("创建软链接的父目录失败: %w", err)
				}
			}

			// 创建软链接
			if err := os.Symlink(target, targetPath); err != nil {
				return fmt.Errorf("创建软链接失败: %w", err)
			}
		default:
			// 普通文件

			// 检查file的父目录是否存在, 如果不存在, 则创建
			parentDir := filepath.Dir(targetPath)
			if _, err := os.Stat(parentDir); os.IsNotExist(err) {
				if err := os.MkdirAll(parentDir, 0644); err != nil {
					return fmt.Errorf("创建父目录失败: %w", err)
				}
			}

			// 创建目标文件
			fileWriter, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
			if err != nil {
				return fmt.Errorf("创建文件失败: %w", err)
			}
			defer fileWriter.Close()

			zipFileReader, err := file.Open()
			if err != nil {
				return fmt.Errorf("打开 ZIP 文件中的文件失败: %w", err)
			}
			defer zipFileReader.Close()

			// 自定义写入器，用于更新进度条
			progressWriter := io.MultiWriter(fileWriter, bar) // bar 是一个全局的进度条对象

			// 使用 io.CopyBuffer 并指定缓冲区大小
			buffer := make([]byte, 512*1024) // 512KB 缓冲区大小
			if _, err := io.CopyBuffer(progressWriter, zipFileReader, buffer); err != nil {
				return fmt.Errorf("写入文件失败: %w", err)
			}
		}
	}

	// 关闭进度条
	if err := bar.Finish(); err != nil {
		return fmt.Errorf("关闭进度条失败: %w", err)
	}

	return nil
}

// ContainsSpecialChars 检测字符串是否包含特殊字符或危险字符
// 参数：s - 需要检测的字符串
// 返回：true 如果包含特殊字符或危险字符，false 否则
func ContainsSpecialChars(s string) bool {
	// 定义危险字符集合，使用哈希表存储以提高查找效率
	dangerousChars := map[rune]bool{
		'<': true, '>': true, '"': true, '\'': true, '\\': true, '`': true, ';': true,
		'%': true, '$': true, '#': true, '@': true, '&': true, '*': true, '(': true,
		')': true, '{': true, '}': true, '[': true, ']': true, '|': true, '!': true,
		'^': true, '~': true, '=': true, '+': true, '-': true, '.': true, ',': true,
		'/': true, '?': true, ':': true, '，': true, '。': true, '？': true, '！': true,
		'；': true, '：': true, '‘': true, '’': true, '“': true, '”': true, '（': true,
		'）': true, '【': true, '】': true, '《': true, '》': true, '…': true, '—': true,
		'～': true, '￥': true, '·': true, '、': true,
	}

	// 遍历字符串中的每个字符
	for _, char := range s {
		// 检查是否为危险字符
		if dangerousChars[char] {
			return true
		}
		// 检查是否为控制字符或非打印字符
		if unicode.IsControl(char) || !unicode.IsPrint(char) {
			return true
		}
	}

	// 如果没有找到任何特殊或危险字符，返回false
	return false
}

// SanitizePath 清理并获取路径的绝对路径
// 参数：path - 需要清理的路径指针
// 返回：error - 如果发生错误，返回错误信息；否则返回 nil
func SanitizePath(path *string) error {
	*path = filepath.Clean(*path)
	absPath, err := filepath.Abs(*path)
	if err != nil {
		return fmt.Errorf("获取路径绝对路径失败: %w", err)
	}
	*path = absPath
	return nil
}

// RenameBackupDirectory 重命名备份目录
// 参数：
//
//	rootPath - 备份存放目录的根路径
//	oldDirName - 旧的备份目录名称
//	newDirName - 新的备份目录名称
//
// 返回值：
//
//	error - 如果发生错误，返回错误信息；否则返回 nil
func RenameBackupDirectory(rootPath, oldDirName, newDirName string) error {
	// 切换到备份存放目录的rootPath
	if err := os.Chdir(rootPath); err != nil {
		return fmt.Errorf("切换到备份存放目录的rootPath失败: %w, Path: %s", err, rootPath)
	}

	// 检查新的备份目录是否存在
	if _, err := CheckPath(newDirName); err == nil {
		return fmt.Errorf("备份目录已存在: %s, 请重试", filepath.Join(rootPath, newDirName))
	}

	// 检查旧的备份目录是否存在
	if _, err := CheckPath(oldDirName); err != nil {
		return fmt.Errorf("旧的备份目录不存在: %s", filepath.Join(rootPath, oldDirName))
	}

	// 重命名备份目录
	if err := os.Rename(oldDirName, newDirName); err != nil {
		return fmt.Errorf("重命名备份目录失败: %w, Old: %s, New: %s", err, oldDirName, newDirName)
	}

	CL.PrintOkf("备份目录重命名成功: %s", filepath.Join(rootPath, newDirName))
	return nil
}

// EnsureDirExists 确保指定目录存在，如果不存在则创建该目录
// 参数:
//
//	dir - 需要检查/创建的目录路径
//
// 返回值:
//
//	error - 如果目录不存在且创建失败，返回错误信息；否则返回nil
func EnsureDirExists(dir string) error {
	// 检查目录是否存在
	if _, err := CheckPath(dir); err != nil {
		// 如果目录不存在，则递归创建目录
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("目录创建失败: %w", err)
		}
	}
	return nil
}

// ParseExclude 解析 --exclude 标志的值并生成排除函数
// 参数:
//
//	excludeValue - 传入的 --exclude 标志的值，多个排除模式使用 | 分隔
//
// 返回值:
//
//	globals.ExcludeFunc - 生成的排除函数，根据不同模式判断文件或目录是否需要排除
//	error - 解析过程中一般不会产生错误，当前固定返回 nil
//
// 特殊值:
//
//	"none" - 表示不排除任何内容，返回一个不排除任何内容的函数
func ParseExclude(excludeValue string) (globals.ExcludeFunc, error) {
	// 特殊值 "none" 表示不排除任何内容
	if excludeValue == "none" {
		// 返回一个不排除任何内容的函数
		return func(path string, info os.FileInfo) bool {
			return false
		}, nil
	}

	// 将传入的排除值按 | 分割成多个排除模式
	patterns := strings.Split(excludeValue, "|")

	// 返回一个排除函数，用于判断文件或目录是否需要排除
	return func(path string, info os.FileInfo) bool {
		// 获取路径的基础文件名
		base := filepath.Base(path)
		// 获取文件的扩展名
		ext := filepath.Ext(path)

		// 遍历所有排除模式
		for _, pattern := range patterns {
			// 跳过空的排除模式
			if pattern == "" {
				continue
			}

			// 根据不同的模式类型进行判断
			switch {
			case strings.HasSuffix(pattern, "/"):
				// 带正斜杠结尾的视为目录排除规则
				// 去掉末尾的 / 得到目录名
				dirPattern := pattern[:len(pattern)-1]
				// 如果当前路径是目录且目录名匹配，则返回 true 表示需要排除
				if info.IsDir() && filepath.Base(path) == dirPattern {
					return true
				}
			case strings.HasPrefix(pattern, "."):
				// 以点开头的视为扩展名排除规则
				// 如果文件扩展名匹配，则返回 true 表示需要排除
				if ext == pattern {
					return true
				}
			default:
				// 其他情况视为普通文件名或通配符模式排除规则
				// 使用 filepath.Match 函数进行匹配
				matched, err := filepath.Match(pattern, base)
				// 如果匹配成功且没有错误，则返回 true 表示需要排除
				if err == nil && matched {
					return true
				}
			}
		}
		// 所有模式都不匹配，返回 false 表示不需要排除
		return false
	}, nil
}
