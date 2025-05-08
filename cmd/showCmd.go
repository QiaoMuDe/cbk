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

// showCmdMain 查询指定任务ID的备份记录并以表格形式输出
func showCmdMain(db *sqlx.DB) error {
	// 检查任务ID是否指定
	if *showID == 0 {
		return fmt.Errorf("查询指定备份任务时, 必须指定任务ID")
	}

	// 构建查询sql语句
	querySql := "SELECT version_id, task_id, timestamp, task_name, backup_status, backup_file_name, backup_size, backup_path, version_hash FROM backup_records WHERE task_id = ? ORDER BY timestamp DESC"

	// 定义存储查询结果的结构体
	var records globals.BackupRecords

	// 执行查询
	if err := db.Select(&records, querySql, *showID); err == sql.ErrNoRows {
		return fmt.Errorf("未找到指定任务ID %d 的备份记录", *showID)
	} else if err != nil {
		return fmt.Errorf("查询备份记录失败: %w", err)
	}

	// 检查是否需要选择完整格式
	if *showView {
		// 禁用表格的输出
		if *showNoTable || *showNoTableShort {
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
		if style, ok := TableStyle[*showTableStyle]; ok {
			t.SetStyle(style)
		} else {
			// 定义样式列表
			var styleList []string
			for k := range TableStyle {
				styleList = append(styleList, k)
			}
			return fmt.Errorf("表格样式不存在: %s, 可选样式: %v", *showTableStyle, styleList)
		}

		// 添加表头
		t.AppendHeader(table.Row{"备份时间", "版本ID", "任务ID", "任务名", "备份状态", "备份文件名", "备份文件大小", "备份文件路径", "版本哈希"})

		// 将查询结果添加到表格
		for _, record := range records {
			// 将时间戳转换为时间对象并格式化为易读格式
			timestamp, err := time.Parse("20060102150405", record.Timestamp)
			if err != nil {
				return fmt.Errorf("解析时间戳失败: %w", err)
			}
			formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")

			// 将数据添加到表格中
			t.AppendRow(table.Row{
				formattedTimestamp, // 格式化后的备份时间
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

		// 输出表格
		t.Render()

		return nil
	}

	// 禁用表格的输出
	if *showNoTable || *showNoTableShort {
		// 打印备份记录
		fmt.Printf("%-25s%-18s%-15s%-20s\n", "备份时间", "版本ID", "任务ID", "任务名")
		for _, record := range records {
			// 将时间戳转换为时间对象并格式化为易读格式
			timestamp, err := time.Parse("20060102150405", record.Timestamp)
			if err != nil {
				return fmt.Errorf("解析时间戳失败: %w", err)
			}
			formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")
			fmt.Printf("%-25s%-25s%-15d%-20s\n", formattedTimestamp, record.VersionID, record.TaskID, record.TaskName)
		}

		return nil

	}

	// 默认为简略格式
	// 创建表格
	t := table.NewWriter()
	// 设置表格输出到标准输出
	t.SetOutputMirror(os.Stdout)
	// 设置表格样式
	if style, ok := TableStyle[*showTableStyle]; ok {
		t.SetStyle(style)
	} else {
		// 定义样式列表
		var styleList []string
		for k := range TableStyle {
			styleList = append(styleList, k)
		}
		return fmt.Errorf("表格样式不存在: %s, 可选样式: %v", *showTableStyle, styleList)
	}
	// 添加表头
	t.AppendHeader(table.Row{"备份时间", "版本ID", "任务ID", "任务名"})

	// 将查询结果添加到表格
	for _, record := range records {
		// 将时间戳转换为时间对象并格式化为易读格式
		timestamp, err := time.Parse("20060102150405", record.Timestamp)
		if err != nil {
			return fmt.Errorf("解析时间戳失败: %w", err)
		}
		formattedTimestamp := timestamp.Format("2006-01-02 15:04:05")

		// 将数据添加到表格中
		t.AppendRow(table.Row{
			formattedTimestamp, // 格式化后的时间戳
			record.VersionID,
			record.TaskID,
			record.TaskName,
		})
	}

	// 设置表格样式
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "版本ID", WidthMax: 10, WidthMaxEnforcer: text.WrapHard},
		{Name: "任务ID", WidthMax: 10, WidthMaxEnforcer: text.WrapHard},
		{Name: "备份时间", WidthMax: 20, WidthMaxEnforcer: text.WrapHard},
		{Name: "任务名", WidthMax: 20, WidthMaxEnforcer: text.WrapHard},
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "版本ID", Align: text.AlignCenter},
		{Name: "任务ID", Align: text.AlignCenter},
		{Name: "备份时间", Align: text.AlignLeft},
		{Name: "任务名", Align: text.AlignLeft},
	})

	// 输出表格
	t.Render()

	return nil
}
