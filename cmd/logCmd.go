package cmd

import (
	"cbk/pkg/globals"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jmoiron/sqlx"
)

// logCmdMain 函数，支持分页查询
// 参数:
//
//	db - 数据库连接
//	page - 页码
//	pageSize - 每页记录数
//
// 返回值:
//
//	error - 如果发生错误，返回错误信息；否则返回 nil
func logCmdMain(db *sqlx.DB, page, pageSize int) error {
	// 验证分页参数
	if page < 1 {
		return fmt.Errorf("页码必须从1开始")
	}
	if pageSize <= 0 {
		return fmt.Errorf("每页记录数必须大于0")
	}

	// 计算 OFFSET
	offset := (page - 1) * pageSize

	// 定义查询语句
	querySql := `
		SELECT version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash
		FROM backup_records
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?;
	`

	// 定义结构体来接收查询结果
	var records globals.BackupRecords

	// 执行查询
	if err := db.Select(&records, querySql, pageSize, offset); err == sql.ErrNoRows {
		return fmt.Errorf("没有查询到任何记录")
	} else if err != nil {
		return fmt.Errorf("查询备份记录失败: %w", err)
	}

	// 选择完整格式
	if *logView {
		// 禁用表格的输出
		if *logNoTable || *logNoTableShort {
			// 打印备份记录
			fmt.Printf("%-25s%-18s%-15s%-20s%-10s%-40s%-30s%-25s%-10s\n", "备份时间", "版本ID", "任务ID", "任务名", "备份状态", "备份文件名", "备份文件大小", "备份存放目录", "版本哈希")
			for _, record := range records {
				// 将时间戳转换为时间对象并格式化为易读格式
				timestamp, err := time.Parse("20060102150405", record.Timestamp)
				if err != nil {
					return fmt.Errorf("解析时间戳失败: %w", err)
				}
				formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")
				fmt.Printf("%-25s%-25s%-15d%-20s%-10s%-40s%-30s%-30s%-10s\n", formattedTimestamp, record.VersionID, record.TaskID, record.TaskName, record.BackupStatus, record.BackupFileName, record.BackupSize, record.BackupPath, record.VersionHash)
			}

			return nil
		}

		// 创建表格
		t := table.NewWriter()

		// 设置表格输出到标准输出
		t.SetOutputMirror(os.Stdout)

		// 设置表格样式
		if style, ok := TableStyle[*logTableStyle]; ok {
			t.SetStyle(style)
		} else {
			// 定义样式列表
			var styleList []string
			for k := range TableStyle {
				styleList = append(styleList, k)
			}
			return fmt.Errorf("表格样式不存在: %s, 可选样式: %v", *logTableStyle, styleList)
		}

		// 添加表头
		t.AppendHeader(table.Row{"备份时间", "版本ID", "任务ID", "任务名", "备份状态", "备份文件名", "备份文件大小", "备份存放目录", "版本哈希"})

		// 遍历查询结果，将数据添加到表格中
		for _, record := range records {
			// 将时间戳转换为时间对象并格式化为易读格式
			timestamp, err := time.Parse("20060102150405", record.Timestamp)
			if err != nil {
				return fmt.Errorf("解析时间戳失败: %w", err)
			}
			formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")

			// 将数据添加到表格中
			t.AppendRow(table.Row{
				formattedTimestamp, // 备份时间
				record.VersionID,
				record.TaskID,
				record.TaskName,
				record.BackupStatus,
				record.BackupFileName,
				record.BackupSize,
				record.BackupPath,
				record.VersionHash,
			})
		}

		// 设置表格样式
		t.SetColumnConfigs([]table.ColumnConfig{
			{Name: "版本ID", WidthMax: 10, WidthMaxEnforcer: text.WrapHard},
			{Name: "任务ID", WidthMax: 10, WidthMaxEnforcer: text.WrapHard},
			{Name: "备份时间", WidthMax: 20, WidthMaxEnforcer: text.WrapHard},
			{Name: "任务名", WidthMax: 20, WidthMaxEnforcer: text.WrapHard},
			{Name: "备份状态", WidthMax: 10, WidthMaxEnforcer: text.WrapHard},
			{Name: "备份文件名", WidthMax: 20, WidthMaxEnforcer: text.WrapHard},
			{Name: "备份文件大小", WidthMax: 10, WidthMaxEnforcer: text.WrapHard},
			{Name: "备份存放目录", WidthMax: 30, WidthMaxEnforcer: text.WrapHard},
			{Name: "版本哈希", WidthMax: 20, WidthMaxEnforcer: text.WrapHard},
		})
		t.SetColumnConfigs([]table.ColumnConfig{
			{Name: "版本ID", Align: text.AlignCenter},
			{Name: "任务ID", Align: text.AlignCenter},
			{Name: "备份时间", Align: text.AlignLeft},
			{Name: "任务名", Align: text.AlignLeft},
			{Name: "备份状态", Align: text.AlignCenter},
			{Name: "备份文件名", Align: text.AlignLeft},
			{Name: "备份文件大小", Align: text.AlignCenter},
			{Name: "备份存放目录", Align: text.AlignLeft},
			{Name: "版本哈希", Align: text.AlignCenter},
		})

		// 打印表格
		t.Render()

		return nil
	}

	// 禁用表格的输出
	if *logNoTable || *logNoTableShort {
		// 打印备份记录
		fmt.Printf("%-25s%-20s%-10s%-40s%-30s%-25s\n", "备份时间", "任务名", "备份状态", "备份文件名", "备份文件大小", "备份存放目录")
		for _, record := range records {
			// 将时间戳转换为时间对象并格式化为易读格式
			timestamp, err := time.Parse("20060102150405", record.Timestamp)
			if err != nil {
				return fmt.Errorf("解析时间戳失败: %w", err)
			}
			formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")
			fmt.Printf("%-25s%-20s%-10s%-40s%-30s%-30s\n", formattedTimestamp, record.TaskName, record.BackupStatus, record.BackupFileName, record.BackupSize, record.BackupPath)
		}

		return nil
	}

	// 默认为简略格式
	// 创建表格
	t := table.NewWriter()

	// 设置表格输出到标准输出
	t.SetOutputMirror(os.Stdout)

	// 设置表格样式
	if style, ok := TableStyle[*logTableStyle]; ok {
		t.SetStyle(style)
	} else {
		// 定义样式列表
		var styleList []string
		for k := range TableStyle {
			styleList = append(styleList, k)
		}
		return fmt.Errorf("表格样式不存在: %s, 可选样式: %v", *logTableStyle, styleList)
	}

	// 添加表头
	t.AppendHeader(table.Row{"备份时间", "任务名", "备份状态", "备份文件名", "备份文件大小", "备份存放目录"})

	// 遍历查询结果，将数据添加到表格中
	for _, record := range records {
		// 将时间戳转换为时间对象并格式化为易读格式
		timestamp, err := time.Parse("20060102150405", record.Timestamp)
		if err != nil {
			return fmt.Errorf("解析时间戳失败: %w", err)
		}
		formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")

		// 将数据添加到表格中
		t.AppendRow(table.Row{
			formattedTimestamp, // 备份时间
			record.TaskName,
			record.BackupStatus,
			record.BackupFileName,
			record.BackupSize,
			record.BackupPath,
		})
	}

	// 设置表格样式
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "备份时间", WidthMax: 20, WidthMaxEnforcer: text.WrapHard},
		{Name: "任务名", WidthMax: 20, WidthMaxEnforcer: text.WrapHard},
		{Name: "备份状态", WidthMax: 10, WidthMaxEnforcer: text.WrapHard},
		{Name: "备份文件名", WidthMax: 20, WidthMaxEnforcer: text.WrapHard},
		{Name: "备份文件大小", WidthMax: 10, WidthMaxEnforcer: text.WrapHard},
		{Name: "备份存放目录", WidthMax: 30, WidthMaxEnforcer: text.WrapHard},
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "备份时间", Align: text.AlignLeft},
		{Name: "任务名", Align: text.AlignLeft},
		{Name: "备份状态", Align: text.AlignCenter},
		{Name: "备份文件名", Align: text.AlignLeft},
		{Name: "备份文件大小", Align: text.AlignCenter},
		{Name: "备份存放目录", Align: text.AlignLeft},
	})

	// 打印表格
	t.Render()

	return nil
}
