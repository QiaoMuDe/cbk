package cmd

import (
	"cbk/pkg/tools"
	"fmt"
	"path/filepath"
	"strings"
)

// unzipCmdMain 解压指定的ZIP文件
func unzipCmdMain() error {
	// 检查是否指定了ZIP文件路径
	if *unzipFile == "" {
		return fmt.Errorf("解压ZIP文件时, 必须指定ZIP文件路径")
	}

	// 对指定的ZIP文件路径进行清理和获取绝对路径
	if err := tools.SanitizePath(unzipFile); err != nil {
		return fmt.Errorf("获取ZIP文件绝对路径失败: %w", err)
	}

	// 对指定的输出目录进行清理和获取绝对路径
	if err := tools.SanitizePath(unzipOutputDir); err != nil {
		return fmt.Errorf("获取输出目录绝对路径失败: %w", err)
	}

	// 检查ZIP文件路径是否以.zip结尾
	if !strings.HasSuffix(*unzipFile, ".zip") {
		return fmt.Errorf("ZIP文件路径必须以.zip结尾: %s", *unzipFile)
	}

	// 检查指定的ZIP文件路径是否存在
	if _, err := tools.CheckPath(*unzipFile); err != nil {
		return fmt.Errorf("指定的ZIP文件路径不存在: %s", *unzipFile)
	}

	// 检查输出目录是否存在
	if _, err := tools.CheckPath(*unzipOutputDir); err != nil {
		return fmt.Errorf("指定的输出目录不存在: %s", *unzipOutputDir)
	}

	// 以点号分隔文件名，获取文件名
	unzipDirName := strings.Split(filepath.Base(*unzipFile), ".")[0]

	// 构建解压后的目录路径
	unzipTargetDir := filepath.Join(*unzipOutputDir, unzipDirName)

	// 检查解压后的目录路径是否存在
	if _, err := tools.CheckPath(unzipTargetDir); err == nil {
		return fmt.Errorf("该路径疑似和解压后的ZIP文件冲突: %s, 请先重命名或移动该路径", unzipTargetDir)
	}

	// 解压ZIP文件
	if err := tools.Unzip(*unzipFile, *unzipOutputDir); err != nil {
		return fmt.Errorf("解压ZIP文件失败: %s", err)
	}

	return nil
}
