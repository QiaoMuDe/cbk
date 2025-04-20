package cmd

import "fmt"

// 打印帮助信息
func helpCmdMain(cmd string) error {
	// 检查是否指定了命令
	if cmd == "" {
		return fmt.Errorf("请指定要查看帮助的命令, 例如: 'cbk help 指定命令'")
	}

	// 根据命令打印帮助信息
	switch cmd {
	case "list":
		fmt.Println(HelpListText)
		return nil
	case "add":
		fmt.Println(HelpAddText)
		return nil
	case "unpack":
		fmt.Println(HelpUnpackText)
		return nil
	case "show":
		fmt.Println(HelpShowText)
		return nil
	case "log":
		fmt.Println(HelpLogText)
		return nil
	case "run":
		fmt.Println(HelpRunText)
		return nil
	case "delete":
		fmt.Println(HelpDeleteText)
		return nil
	case "zip":
		fmt.Println(HelpZipText)
		return nil
	case "unzip":
		fmt.Println(HelpUnzipText)
		return nil
	case "edit":
		fmt.Println(HelpEditText)
		return nil
	case "init":
		fmt.Println(HelpInitText)
		return nil
	case "export":
		fmt.Println(HelpExportText)
		return nil
	default:
		return fmt.Errorf("未知命令: %s", cmd)
	}
}
