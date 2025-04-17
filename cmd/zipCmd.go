package cmd

import (
	"cbk/pkg/tools"
	"fmt"
	"strings"
)

// zipCmdMain 压缩指定目录下的文件
func zipCmdMain() error {
	// 检查是否指定了ZIP文件名
	if *zipOutput == "" {
		return fmt.Errorf("打包ZIP文件时, 必须指定ZIP文件名")
	}

	// 检查是否指定了目录路径
	if *zipTarget == "" {
		return fmt.Errorf("打包ZIP文件时, 必须指定目录路径")
	}

	// 基本格式检查
	if !strings.HasSuffix(*zipOutput, ".zip") {
		return fmt.Errorf("ZIP文件路径必须以.zip结尾: %s", *zipOutput)
	}

	// 清理路径并获取绝对路径
	if err := tools.SanitizePath(zipOutput); err != nil {
		return fmt.Errorf("清理路径并获取绝对路径失败: %w", err)
	}

	// 清理路径并获取绝对路径
	if err := tools.SanitizePath(zipTarget); err != nil {
		return fmt.Errorf("清理路径并获取绝对路径失败: %w", err)
	}

	// 检查指定的ZIP文件路径是否存在
	if info, err := tools.CheckPath(*zipOutput); err == nil {
		// 如果路径存在
		if info.Exists {
			// 如果路径存在且是一个文件
			if info.IsFile {
				return fmt.Errorf("指定的ZIP文件已存在: %s", *zipOutput)
			}
			// 如果路径存在但不是一个文件（例如是一个目录）
			return fmt.Errorf("指定的路径存在，但不是一个文件: %s", *zipOutput)
		}
	}

	// 检查指定的目录路径是否存在
	if _, err := tools.CheckPath(*zipTarget); err != nil {
		return fmt.Errorf("指定的目录路径不存在: %s", *zipTarget)
	}

	// 获取是否禁用压缩
	var noCompression int
	if *zipNoCompression {
		noCompression = 1 // 禁用压缩
	} else {
		noCompression = 0 // 默认启用压缩
	}

	// 创建ZIP文件
	if err := tools.CreateZip(*zipOutput, *zipTarget, noCompression); err != nil {
		return fmt.Errorf("创建ZIP文件失败: %w", err)
	}

	return nil
}
