package cmd

import (
	"fmt"
	"runtime"
)

// completeCmdMain 自动补全主逻辑
func completeCmdMain(t string) error {
	// 检查自动补全类型是否为空
	if t == "" {
		return fmt.Errorf("请指定自动补全类型, 例如: 'cbk complete -type bash'")
	}

	switch t {
	case "bash":
		// 检查是否为Linux或Mac系统
		if runtime.GOOS != "linux" {
			return fmt.Errorf("自动补全类型 'bash' 仅在Linux上受支持")
		}

		// 打印自动补全脚本
		fmt.Println(BashCompletion)
		return nil
	default:
		return fmt.Errorf("未知的自动补全类型: %s", t)
	}
}
