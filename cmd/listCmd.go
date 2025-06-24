package cmd

import (
	"cbk/pkg/globals"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// listCmdMain 查询并打印任务列表
func listCmdMain(db *sqlx.DB) error {
	// 创建表格
	t := table.NewWriter()

	// 设置表格样式
	if style, ok := globals.TableStyle[listTableStyle.Get()]; ok {
		t.SetStyle(style)
	} else {
		return fmt.Errorf("表格样式不存在: %s, 可选样式: %v", listTableStyle.Get(), globals.TableStyleList)
	}

	// 查询所有任务
	querySql := "SELECT task_id, task_name, target_directory, backup_directory, retention_count, retention_days, no_compression, exclude_rules FROM backup_tasks;"

	// 定义存储查询结果的结构体
	var tasks globals.BackupTasks

	// 查询任务列表
	if err := db.Select(&tasks, querySql); err != nil {
		return fmt.Errorf("查询任务失败: %w", err)
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
				// 根据NoCompression字段返回true或false
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
