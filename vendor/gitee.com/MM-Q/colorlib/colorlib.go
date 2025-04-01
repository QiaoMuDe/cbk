package colorlib

import (
	"fmt"
	"strings"
)

type ColorLib struct {
	// 空结构体，用于实现接口
	LevelMap map[string]string // LevelMap 是一个映射，用于将日志级别映射到对应的前缀,// 日志级别映射到对应的前缀, 后面留空, 方便后面拼接提示内容
	ColorMap map[string]string // colorMap 是一个映射，用于将颜色名称映射到对应的 ANSI 颜色代码。
	LvlMap   map[string]string // LvlMap 是一个映射，用于将日志级别映射到对应的日志级别名称。
}

// ColorLibInterface 是一个接口，定义了一组方法，用于打印和返回带有颜色的文本。
type ColorLibInterface interface {
	// 需要占位符的方法(自带换行符)
	Bluef(format string, a ...any)           // 打印蓝色信息到控制台（带占位符）
	Greenf(format string, a ...any)          // 打印绿色信息到控制台（带占位符）
	Redf(format string, a ...any)            // 打印红色信息到控制台（带占位符）
	Yellowf(format string, a ...any)         // 打印黄色信息到控制台（带占位符）
	Purplef(format string, a ...any)         // 打印紫色信息到控制台（带占位符）
	Sbluef(format string, a ...any) string   // 返回构造后的蓝色字符串（带占位符）
	Sgreenf(format string, a ...any) string  // 返回构造后的绿色字符串（带占位符）
	Sredf(format string, a ...any) string    // 返回构造后的红色字符串（带占位符）
	Syellowf(format string, a ...any) string // 返回构造后的黄色字符串（带占位符）
	Spurplef(format string, a ...any) string // 返回构造后的紫色字符串（带占位符）
	PrintSuccessf(format string, a ...any)   // 打印成功信息到控制台（带占位符）
	PrintErrorf(format string, a ...any)     // 打印错误信息到控制台（带占位符）
	PrintWarningf(format string, a ...any)   // 打印警告信息到控制台（带占位符）
	PrintInfof(format string, a ...any)      // 打印信息到控制台（带占位符）
	PrintDebugf(format string, a ...any)     // 打印调试信息到控制台（带占位符）

	// 直接打印信息, 无需占位符
	Blue(msg ...any)           // 打印蓝色信息到控制台, 无需占位符
	Green(msg ...any)          // 打印绿色信息到控制台, 无需占位符
	Red(msg ...any)            // 打印红色信息到控制台, 无需占位符
	Yellow(msg ...any)         // 打印黄色信息到控制台, 无需占位符
	Purple(msg ...any)         // 打印紫色信息到控制台, 无需占位符
	Sblue(msg ...any) string   // 返回构造后的蓝色字符串, 无需占位符
	Sgreen(msg ...any) string  // 返回构造后的绿色字符串, 无需占位符
	Sred(msg ...any) string    // 返回构造后的红色字符串, 无需占位符
	Syellow(msg ...any) string // 返回构造后的黄色字符串, 无需占位符
	Spurple(msg ...any) string // 返回构造后的紫色字符串, 无需占位符
	PrintSuccess(msg ...any)   // 打印成功信息到控制台, 无需占位符
	PrintError(msg ...any)     // 打印错误信息到控制台, 无需占位符
	PrintWarning(msg ...any)   // 打印警告信息到控制台, 无需占位符
	PrintInfo(msg ...any)      // 打印信息到控制台, 无需占位符
	PrintDebug(msg ...any)     // 打印调试信息到控制台, 无需占位符

	// 新增扩展颜色的方法
	Black(msg ...any)                         // 打印黑色信息到控制台, 无需占位符
	Blackf(format string, a ...any)           // 打印黑色信息到控制台（带占位符）
	Sblack(msg ...any) string                 // 返回构造后的黑色字符串, 无需占位符
	Sblackf(format string, a ...any) string   // 返回构造后的黑色字符串（带占位符）
	Cyan(msg ...any)                          // 打印青色信息到控制台, 无需占位符
	Cyanf(format string, a ...any)            // 打印青色信息到控制台（带占位符）
	Scyan(msg ...any) string                  // 返回构造后的青色字符串, 无需占位符
	Scyanf(format string, a ...any) string    // 返回构造后的青色字符串（带占位符）
	White(msg ...any)                         // 打印白色信息到控制台, 无需占位符
	Whitef(format string, a ...any)           // 打印白色信息到控制台（带占位符）
	Swhite(msg ...any) string                 // 返回构造后的白色字符串, 无需占位符
	Swhitef(format string, a ...any) string   // 返回构造后的白色字符串（带占位符）
	Gray(msg ...any)                          // 打印灰色信息到控制台, 无需占位符
	Grayf(format string, a ...any)            // 打印灰色信息到控制台（带占位符）
	Sgray(msg ...any) string                  // 返回构造后的灰色字符串, 无需占位符
	Sgrayf(format string, a ...any) string    // 返回构造后的灰色字符串（带占位符）
	Lred(msg ...any)                          // 打印亮红色信息到控制台, 无需占位符
	Lredf(format string, a ...any)            // 打印亮红色信息到控制台（带占位符）
	Slred(msg ...any) string                  // 返回构造后的亮红色字符串, 无需占位符
	Slredf(format string, a ...any) string    // 返回构造后的亮红色字符串（带占位符）
	Lgreen(msg ...any)                        // 打印亮绿色信息到控制台, 无需占位符
	Lgreenf(format string, a ...any)          // 打印亮绿色信息到控制台（带占位符）
	Slgreen(msg ...any) string                // 返回构造后的亮绿色字符串, 无需占位符
	Slgreenf(format string, a ...any) string  // 返回构造后的亮绿色字符串（带占位符）
	Lyellow(msg ...any)                       // 打印亮黄色信息到控制台, 无需占位符
	Lyellowf(format string, a ...any)         // 打印亮黄色信息到控制台（带占位符）
	Slyellow(msg ...any) string               // 返回构造后的亮黄色字符串, 无需占位符
	Slyellowf(format string, a ...any) string // 返回构造后的亮黄色字符串（带占位符）
	Lblue(msg ...any)                         // 打印亮蓝色信息到控制台, 无需占位符
	Lbluef(format string, a ...any)           // 打印亮蓝色信息到控制台（带占位符）
	Slblue(msg ...any) string                 // 返回构造后的亮蓝色字符串, 无需占位符
	Slbluef(format string, a ...any) string   // 返回构造后的亮蓝色字符串（带占位符）
	Lpurple(msg ...any)                       // 打印亮紫色信息到控制台, 无需占位符
	Lpurplef(format string, a ...any)         // 打印亮紫色信息到控制台（带占位符）
	Slpurple(msg ...any) string               // 返回构造后的亮紫色字符串, 无需占位符
	Slpurplef(format string, a ...any) string // 返回构造后的亮紫色字符串（带占位符）
	Lcyan(msg ...any)                         // 打印亮青色信息到控制台, 无需占位符
	Lcyanf(format string, a ...any)           // 打印亮青色信息到控制台（带占位符）
	Slcyan(msg ...any) string                 // 返回构造后的亮青色字符串, 无需占位符
	Slcyanf(format string, a ...any) string   // 返回构造后的亮青色字符串（带占位符）
	Lwhite(msg ...any)                        // 打印亮白色信息到控制台, 无需占位符
	Lwhitef(format string, a ...any)          // 打印亮白色信息到控制台（带占位符）
	Slwhite(msg ...any) string                // 返回构造后的亮白色字符串, 无需占位符
	Slwhitef(format string, a ...any) string  // 返回构造后的亮白色字符串（带占位符）

	// 新增简洁版的方法, 无需占位符
	PrintOk(msg ...any)   // 打印成功信息到控制台, 无需占位符
	PrintErr(msg ...any)  // 打印错误信息到控制台, 无需占位符
	PrintInf(msg ...any)  // 打印信息到控制台, 无需占位符
	PrintDbg(msg ...any)  // 打印调试信息到控制台, 无需占位符
	PrintWarn(msg ...any) // 打印警告信息到控制台, 无需占位符

	// 新增简洁版的方法, 带占位符
	PrintOkf(format string, a ...any)   // 打印成功信息到控制台（带占位符）
	PrintErrf(format string, a ...any)  // 打印错误信息到控制台（带占位符）
	PrintInff(format string, a ...any)  // 打印信息到控制台（带占位符）
	PrintDbgf(format string, a ...any)  // 打印调试信息到控制台（带占位符）
	PrintWarnf(format string, a ...any) // 打印警告信息到控制台（带占位符）
}

// NewColorLib 函数用于创建一个新的 ColorLib 实例。
func NewColorLib() *ColorLib {
	// 创建一个新的 ColorLib 实例
	cl := &ColorLib{
		// LvlMap 是一个映射，用于将日志级别映射到对应的日志级别名称。
		LvlMap: map[string]string{
			"ok":   "ok: ",
			"err":  "err: ",
			"inf":  "info: ",
			"dbg":  "debug: ",
			"warn": "warn:",
		},
		// LevelMap 是一个映射，用于将日志级别映射到对应的前缀。
		LevelMap: map[string]string{
			"success": "[Success] ", // 成功信息级别的前缀
			"error":   "[Error] ",   // 错误信息级别的前缀
			"warning": "[Warning] ", // 警告信息级别的前缀
			"info":    "[Info] ",    // 信息信息级别的前缀
			"debug":   "[Debug] ",   // 调试信息级别的前缀
		},
		// colorMap 是一个映射，用于将颜色名称映射到对应的 ANSI 颜色代码。
		ColorMap: map[string]string{
			"black":   "30", // 黑色文本的 ANSI 颜色代码
			"red":     "31", // 红色文本的 ANSI 颜色代码
			"green":   "32", // 绿色文本的 ANSI 颜色代码
			"yellow":  "33", // 黄色文本的 ANSI 颜色代码
			"blue":    "34", // 蓝色文本的 ANSI 颜色代码
			"purple":  "35", // 紫色文本的 ANSI 颜色代码
			"cyan":    "36", // 青色文本的 ANSI 颜色代码
			"white":   "37", // 白色文本的 ANSI 颜色代码
			"gray":    "90", // 灰色文本的 ANSI 颜色代码
			"lred":    "91", // 亮红色文本的 ANSI 颜色代码
			"lgreen":  "92", // 亮绿色文本的 ANSI 颜色代码
			"lyellow": "93", // 亮黄色文本的 ANSI 颜色代码
			"lblue":   "94", // 亮蓝色文本的 ANSI 颜色代码
			"lpurple": "95", // 亮紫色文本的 ANSI 颜色代码
			"lcyan":   "96", // 亮青色文本的 ANSI 颜色代码
			"lwhite":  "97", // 亮白色文本的 ANSI 颜色代码
		},
	}

	// 返回 ColorLib 实例
	return cl
}

// printWithColor 方法用于将传入的参数以指定颜色文本形式打印到控制台。
func (c *ColorLib) printWithColor(color string, msg ...any) {
	// 获取颜色代码
	code, ok := c.ColorMap[color]
	if !ok {
		fmt.Println("Invalid color:", color)
		return
	}

	// 使用 strings.Builder 拼接字符串
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("\033[1;%sm", code)) // 添加颜色代码，并加粗

	// 检查 msg 是否为空
	if len(msg) > 0 {
		builder.WriteString(fmt.Sprint(msg...)) // 拼接消息内容
	} else {
		builder.WriteString(" ") // 如果没有消息，添加一个空格，避免完全空白的输出
	}

	builder.WriteString("\033[0m") // 添加颜色重置代码
	fmt.Println(builder.String())  // 使用 fmt.Println 自动处理换行
}

// returnWithColor 方法用于将传入的参数以指定颜色文本形式返回。
func (c *ColorLib) returnWithColor(color string, msg ...any) string {
	// 获取颜色代码
	code, ok := c.ColorMap[color]
	if !ok {
		return fmt.Sprintf("Invalid color: %s", color)
	}

	// 检查 msg 是否为空
	if len(msg) == 0 {
		return fmt.Sprintf("\033[1;%sm\033[0m", code) // 返回空字符串，但带有颜色代码
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return fmt.Sprintf("\033[1;%sm%s\033[0m", code, combinedMsg)
}

// Bluef 方法用于将传入的参数以蓝色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Bluef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("blue", formattedMsg)
}

// Greenf 方法用于将传入的参数以绿色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Greenf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("green", formattedMsg)
}

// Redf 方法用于将传入的参数以红色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Redf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("red", formattedMsg)
}

// Yellowf 方法用于将传入的参数以黄色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Yellowf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("yellow", formattedMsg)
}

// Purplef 方法用于将传入的参数以紫色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Purplef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	c.printWithColor("purple", formattedMsg)
}

// Sbluef 方法用于将传入的参数以蓝色文本形式返回（带占位符）。
func (c *ColorLib) Sbluef(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("blue", formattedMsg)
}

// Sgreenf 方法用于将传入的参数以绿色文本形式返回（带占位符）。
func (c *ColorLib) Sgreenf(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("green", formattedMsg)
}

// Sredf 方法用于将传入的参数以红色文本形式返回（带占位符）。
func (c *ColorLib) Sredf(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("red", formattedMsg)
}

// Syellowf 方法用于将传入的参数以黄色文本形式返回（带占位符）。
func (c *ColorLib) Syellowf(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("yellow", formattedMsg)
}

// Spurplef 方法用于将传入的参数以紫色文本形式返回（带占位符）。
func (c *ColorLib) Spurplef(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("purple", formattedMsg)
}

// Blue 方法用于将传入的参数以蓝色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Blue(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("blue", combinedMsg)
}

// Green 方法用于将传入的参数以绿色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Green(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("green", combinedMsg)
}

// Red 方法用于将传入的参数以红色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Red(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("red", combinedMsg)
}

// Yellow 方法用于将传入的参数以黄色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Yellow(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("yellow", combinedMsg)
}

// Purple 方法用于将传入的参数以紫色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Purple(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}
	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("purple", combinedMsg)
}

// Sblue 方法用于将传入的参数以蓝色文本形式返回（不带占位符）。
func (c *ColorLib) Sblue(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("blue", combinedMsg)
}

// Sgreen 方法用于将传入的参数以绿色文本形式返回（不带占位符）。
func (c *ColorLib) Sgreen(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("green", combinedMsg)
}

// Sred 方法用于将传入的参数以红色文本形式返回（不带占位符）。
func (c *ColorLib) Sred(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("red", combinedMsg)
}

// Syellow 方法用于将传入的参数以黄色文本形式返回（不带占位符）。
func (c *ColorLib) Syellow(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("yellow", combinedMsg)
}

// Spurple 方法用于将传入的参数以紫色文本形式返回（不带占位符）。
func (c *ColorLib) Spurple(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("purple", combinedMsg)
}

// PromptMsg 方法用于打印带有颜色和前缀的消息。
func (c *ColorLib) PromptMsg(level, color, format string, a ...any) {
	// 获取指定级别对应的前缀
	prefix, ok := c.LevelMap[level]
	if !ok {
		fmt.Println("Invalid level:", level)
		return
	}

	// 创建一个 strings.Builder 来构建消息
	var message strings.Builder
	message.WriteString(prefix)

	// 如果没有参数，直接打印前缀
	if len(a) == 0 {
		c.printWithColor(color, message.String())
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprintf(format, a...)
	message.WriteString(combinedMsg)

	// 打印最终消息
	c.printWithColor(color, message.String())
}

// PMsg 方法用于打印带有颜色和前缀的消息。
func (c *ColorLib) PMsg(level, color, format string, a ...any) {
	// 获取指定级别对应的前缀
	prefix, ok := c.LvlMap[level]
	if !ok {
		fmt.Println("Invalid level:", level)
		return
	}

	// 创建一个 strings.Builder 来构建消息
	var message strings.Builder
	message.WriteString(prefix)

	// 如果没有参数，直接打印前缀
	if len(a) == 0 {
		c.printWithColor(color, message.String())
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprintf(format, a...)
	message.WriteString(combinedMsg)

	// 打印最终消息
	c.printWithColor(color, message.String())
}

// PrintSuccessf 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（带占位符）。
func (c *ColorLib) PrintSuccessf(format string, a ...any) {
	c.PromptMsg("success", "green", format, a...)
}

// PrintErrorf 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（带占位符）。
func (c *ColorLib) PrintErrorf(format string, a ...any) {
	c.PromptMsg("error", "red", format, a...)
}

// PrintWarningf 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（带占位符）。
func (c *ColorLib) PrintWarningf(format string, a ...any) {
	c.PromptMsg("warning", "yellow", format, a...)
}

// PrintInfof 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（带占位符）。
func (c *ColorLib) PrintInfof(format string, a ...any) {
	c.PromptMsg("info", "blue", format, a...)
}

// PrintDebugf 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（带占位符）。
func (c *ColorLib) PrintDebugf(format string, a ...any) {
	c.PromptMsg("debug", "purple", format, a...)
}

// PrintSuccess 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（不带占位符）。
func (c *ColorLib) PrintSuccess(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.PromptMsg("success", "green", "%s", "")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.PromptMsg("success", "green", "%s", combinedMsg)
}

// PrintError 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（不带占位符）。
func (c *ColorLib) PrintError(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.PromptMsg("error", "red", "%s", "")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.PromptMsg("error", "red", "%s", combinedMsg)
}

// PrintWarning 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（不带占位符）。
func (c *ColorLib) PrintWarning(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.PromptMsg("warning", "yellow", "%s", "")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.PromptMsg("warning", "yellow", "%s", combinedMsg)
}

// PrintInfo 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（不带占位符）。
func (c *ColorLib) PrintInfo(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.PromptMsg("info", "blue", "%s", "")
		return
	}

	// 使用 fmt.Sprint 将 msg 中的所有元素拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.PromptMsg("info", "blue", "%s", combinedMsg)
}

// PrintDebug 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（不带占位符）。
func (c *ColorLib) PrintDebug(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.PromptMsg("debug", "purple", "%s", "")
		return
	}

	// 使用 fmt.Sprint 将 msg 中的所有元素拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.PromptMsg("debug", "purple", "%s", combinedMsg)
}

// 扩展补充颜色 黑，青，白，灰

// Black 方法用于将传入的参数以黑色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Black(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("black", combinedMsg)
}

// Blackf 方法用于将传入的参数以黑色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Blackf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("black", formattedMsg)
}

// Sblack 方法用于将传入的参数以黑色文本形式返回（不带占位符）。
func (c *ColorLib) Sblack(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("black", combinedMsg)
}

// Sblackf 方法用于将传入的参数以黑色文本形式返回（带占位符）。
func (c *ColorLib) Sblackf(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("black", formattedMsg)
}

// Cyan 方法用于将传入的参数以青色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Cyan(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("cyan", combinedMsg)
}

// Cyanf 方法用于将传入的参数以青色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Cyanf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("cyan", formattedMsg)
}

// Scyan 方法用于将传入的参数以青色文本形式返回（不带占位符）。
func (c *ColorLib) Scyan(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("cyan", combinedMsg)
}

// Scyanf 方法用于将传入的参数以青色文本形式返回（带占位符）。
func (c *ColorLib) Scyanf(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("cyan", formattedMsg)
}

// White 方法用于将传入的参数以白色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) White(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("white", combinedMsg)
}

// Whitef 方法用于将传入的参数以白色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Whitef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("white", formattedMsg)
}

// Swhite 方法用于将传入的参数以白色文本形式返回（不带占位符）。
func (c *ColorLib) Swhite(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("white", combinedMsg)
}

// Swhitef 方法用于将传入的参数以白色文本形式返回（带占位符）。
func (c *ColorLib) Swhitef(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("white", formattedMsg)
}

// Gray 方法用于将传入的参数以灰色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Gray(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("gray", combinedMsg)
}

// Grayf 方法用于将传入的参数以灰色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Grayf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("gray", formattedMsg)
}

// Sgray 方法用于将传入的参数以灰色文本形式返回（不带占位符）。
func (c *ColorLib) Sgray(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("gray", combinedMsg)
}

// Sgrayf 方法用于将传入的参数以灰色文本形式返回（带占位符）。
func (c *ColorLib) Sgrayf(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("gray", formattedMsg)
}

// 扩展颜色 亮红色 亮绿色 亮黄色 亮蓝色 亮紫色 亮青色 亮白色
// Lred 方法用于将传入的参数以亮红色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lred(msg ...any) {
	// 检查传入的参数数量，如果为0，则直接返回
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	// 调用 ColorLib 类型的 printWithColor 方法，传入颜色 "lred" 和拼接后的字符串
	c.printWithColor("lred", combinedMsg)
}

// Lredf 方法用于将传入的参数以亮红色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lredf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lred", formattedMsg)
}

// Slred 方法用于将传入的参数以亮红色文本形式返回（不带占位符）。
func (c *ColorLib) Slred(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lred", combinedMsg)
}

// Slredf 方法用于将传入的参数以亮红色文本形式返回（带占位符）。
func (c *ColorLib) Slredf(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lred", formattedMsg)
}

// Lgreen 方法用于将传入的参数以亮绿色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lgreen(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("lgreen", combinedMsg)
}

// Lgreenf 方法用于将传入的参数以亮绿色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lgreenf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lgreen", formattedMsg)
}

// Slgreen 方法用于将传入的参数以亮绿色文本形式返回（不带占位符）。
func (c *ColorLib) Slgreen(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}
	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lgreen", combinedMsg)
}

// Slgreenf 方法用于将传入的参数以亮绿色文本形式返回（带占位符）。
func (c *ColorLib) Slgreenf(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lgreen", formattedMsg)
}

// Lyellow 方法用于将传入的参数以亮黄色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lyellow(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("lyellow", combinedMsg)
}

// Lyellowf 方法用于将传入的参数以亮黄色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lyellowf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lyellow", formattedMsg)
}

// Slblue 方法用于将传入的参数以亮蓝色文本形式返回（不带占位符）。
func (c *ColorLib) Slyellow(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lyellow", combinedMsg)
}

// Slbluef 方法用于将传入的参数以亮蓝色文本形式返回（带占位符）。
func (c *ColorLib) Slyellowf(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lyellow", formattedMsg)
}

// Lblue 方法用于将传入的参数以亮蓝色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lblue(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("lblue", combinedMsg)
}

// Lbluef 方法用于将传入的参数以亮蓝色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lbluef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lblue", formattedMsg)
}

// Slblue 方法用于将传入的参数以亮蓝色文本形式返回（不带占位符）。
func (c *ColorLib) Slblue(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lblue", combinedMsg)
}

// Slbluef 方法用于将传入的参数以亮蓝色文本形式返回（带占位符）。
func (c *ColorLib) Slbluef(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lblue", formattedMsg)
}

// Lgreen 方法用于将传入的参数以亮绿色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lpurple(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("lpurple", combinedMsg)
}

// Lgreenf 方法用于将传入的参数以亮绿色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lpurplef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lpurple", formattedMsg)
}

// Slgreen 方法用于将传入的参数以亮绿色文本形式返回（不带占位符）。
func (c *ColorLib) Slpurple(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lpurple", combinedMsg)
}

// Slgreenf 方法用于将传入的参数以亮绿色文本形式返回（带占位符）。
func (c *ColorLib) Slpurplef(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lpurple", formattedMsg)
}

// Lcyan 方法用于将传入的参数以亮青色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lcyan(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("lcyan", combinedMsg)
}

// Lcyanf 方法用于将传入的参数以亮青色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lcyanf(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lcyan", formattedMsg)
}

// Slcyan 方法用于将传入的参数以亮青色文本形式返回（不带占位符）。
func (c *ColorLib) Slcyan(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lcyan", combinedMsg)
}

// Slcyanf 方法用于将传入的参数以亮青色文本形式返回（带占位符）。
func (c *ColorLib) Slcyanf(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lcyan", formattedMsg)
}

// Lwhite 方法用于将传入的参数以亮白色文本形式打印到控制台（不带占位符）。
func (c *ColorLib) Lwhite(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.printWithColor("lwhite", combinedMsg)
}

// Lwhitef 方法用于将传入的参数以亮白色文本形式打印到控制台（带占位符）。
func (c *ColorLib) Lwhitef(format string, a ...any) {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 printWithColor 方法，传入格式化后的字符串
	c.printWithColor("lwhite", formattedMsg)
}

// Slwhite 方法用于将传入的参数以亮白色文本形式返回（不带占位符）。
func (c *ColorLib) Slwhite(msg ...any) string {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串
		return ""
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	return c.returnWithColor("lwhite", combinedMsg)
}

// Slwhitef 方法用于将传入的参数以亮白色文本形式返回（带占位符）。
func (c *ColorLib) Slwhitef(format string, a ...any) string {
	// 使用 fmt.Sprintf 格式化参数
	formattedMsg := fmt.Sprintf(format, a...)

	// 调用 returnWithColor 方法，传入格式化后的字符串
	return c.returnWithColor("lwhite", formattedMsg)
}

// 扩展简洁方法 错误 成功 警告 信息 调试
// PrintOk 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（不带占位符）。
func (c *ColorLib) PrintOk(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.PMsg("ok", "green", "%s", "")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.PMsg("ok", "green", "%s", combinedMsg)
}

// PrintErr 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（不带占位符）。
func (c *ColorLib) PrintErr(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.PMsg("err", "red", "%s", "")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.PMsg("err", "red", "%s", combinedMsg)
}

// PrintWarn 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（不带占位符）。
func (c *ColorLib) PrintWarn(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.PMsg("warn", "yellow", "%s", "")
		return
	}

	// 使用 fmt.Sprint 将所有参数拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.PMsg("warn", "yellow", "%s", combinedMsg)
}

// PrintInf 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（不带占位符）。
func (c *ColorLib) PrintInf(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.PMsg("inf", "blue", "%s", "")
		return
	}

	// 使用 fmt.Sprint 将 msg 中的所有元素拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.PMsg("inf", "blue", "%s", combinedMsg)
}

// PrintDbg 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（不带占位符）。
func (c *ColorLib) PrintDbg(msg ...any) {
	if len(msg) == 0 {
		// 如果没有传入任何参数，直接返回空字符串或默认消息
		c.PMsg("dbg", "purple", "%s", "")
		return
	}

	// 使用 fmt.Sprint 将 msg 中的所有元素拼接成一个字符串
	combinedMsg := fmt.Sprint(msg...)
	c.PMsg("dbg", "purple", "%s", combinedMsg)
}

// PrintOkf 方法用于将传入的参数以绿色文本形式打印到控制台，并在文本前添加一个表示成功的标志（带占位符）。
func (c *ColorLib) PrintOkf(format string, a ...any) {
	// 调用 PMsg 方法，传入格式化后的字符串
	c.PMsg("ok", "green", format, a...)
}

// PrintErrf 方法用于将传入的参数以红色文本形式打印到控制台，并在文本前添加一个表示错误的标志（带占位符）。
func (c *ColorLib) PrintErrf(format string, a ...any) {
	// 调用 PMsg 方法，传入格式化后的字符串
	c.PMsg("err", "red", format, a...)
}

// PrintWarnf 方法用于将传入的参数以黄色文本形式打印到控制台，并在文本前添加一个表示警告的标志（带占位符）。
func (c *ColorLib) PrintWarnf(format string, a ...any) {
	// 调用 PMsg 方法，传入格式化后的字符串
	c.PMsg("warn", "yellow", format, a...)
}

// PrintInff 方法用于将传入的参数以蓝色文本形式打印到控制台，并在文本前添加一个表示信息的标志（带占位符）。
func (c *ColorLib) PrintInff(format string, a ...any) {
	// 调用 PMsg 方法，传入格式化后的字符串
	c.PMsg("inf", "blue", format, a...)
}

// PrintDbgf 方法用于将传入的参数以紫色文本形式打印到控制台，并在文本前添加一个表示调试的标志（带占位符）。
func (c *ColorLib) PrintDbgf(format string, a ...any) {
	// 调用 PMsg 方法，传入格式化后的字符串
	c.PMsg("dbg", "purple", format, a...)
}
