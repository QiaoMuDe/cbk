package cmd

import (
	"cbk/pkg/globals"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jmoiron/sqlx"
)

// listCmdMain 查询并打印任务列表
func listCmdMain(db *sqlx.DB) error {
	// 创建表格
	t := table.NewWriter()

	// 设置表格样式
	if style, ok := TableStyle[*listTableStyle]; ok {
		t.SetStyle(style)
	} else {
		// 定义样式列表
		var styleList []string
		for k := range TableStyle {
			styleList = append(styleList, k)
		}
		return fmt.Errorf("表格样式不存在: %s, 可选样式: %v", *listTableStyle, styleList)
	}

	// 查询所有任务
	querySql := "SELECT task_id, task_name, target_directory, backup_directory, retention_count, retention_days, no_compression, exclude_rules FROM backup_tasks;"

	// 定义存储查询结果的结构体
	var tasks globals.BackupTasks

	// 查询任务列表
	if err := db.Select(&tasks, querySql); err != nil {
		return fmt.Errorf("查询任务失败: %w", err)
	}

	// 禁用表格的输出
	if *listNoTable || *listNoTableShort {
		// 打印任务列表
		fmt.Printf("%-30s %-10s %-15s %-15s %-30s %-30s %-20s %-30s\n",
			"任务名", "任务ID", "保留数量", "保留天数", "目标目录", "备份目录", "是否禁用压缩", "排除规则")
		for _, task := range tasks {
			fmt.Printf("%-30s %-10d %-15d %-15d %-30s %-30s %-10s %-30s\n", task.TaskName, task.TaskID, task.RetentionCount, task.RetentionDays, task.TargetDirectory, task.BackupDirectory, func() string {
				if task.NoCompression == 0 {
					return "false"
				} else {
					return "true"
				}
			}(), task.ExcludeRules)
		}

		return nil
	}

	// 使用标准输出作为输出目标
	t.SetOutputMirror(os.Stdout)

	// 设置表头
	t.AppendHeader(table.Row{"ID", "任务名", "保留数量", "保留天数", "目标目录", "备份目录", "是否禁用压缩", "排除规则"})

	// 设置列配置
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "ID", Align: text.AlignCenter, WidthMaxEnforcer: text.WrapHard},
		{Name: "任务名", Align: text.AlignLeft, WidthMaxEnforcer: text.WrapHard},
		{Name: "保留数量", Align: text.AlignCenter, WidthMaxEnforcer: text.WrapHard},
		{Name: "保留天数", Align: text.AlignCenter, WidthMaxEnforcer: text.WrapHard},
		{Name: "目标目录", Align: text.AlignLeft, WidthMaxEnforcer: text.WrapHard},
		{Name: "备份目录", Align: text.AlignLeft, WidthMaxEnforcer: text.WrapHard},
		{Name: "是否禁用压缩", Align: text.AlignCenter, WidthMaxEnforcer: text.WrapHard},
		{Name: "排除规则", Align: text.AlignLeft, WidthMaxEnforcer: text.WrapHard},
	})

	// 添加数据行
	for _, task := range tasks {
		t.AppendRow(table.Row{
			task.TaskID,
			task.TaskName,
			task.RetentionCount,
			task.RetentionDays,
			task.TargetDirectory,
			task.BackupDirectory,
			func() string {
				if task.NoCompression == 0 {
					return "false"
				} else {
					return "true"
				}
			}(),
			task.ExcludeRules,
		})
	}

	// 渲染表格
	t.Render()

	return nil
}
