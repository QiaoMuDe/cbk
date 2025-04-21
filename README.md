# CBK - 备份管理工具

CBK 是一个基于命令行的备份管理工具，支持任务的添加、删除、编辑、运行、日志查看等功能。通过简单的命令行操作，用户可以轻松管理备份任务，并查看备份记录。

## 安装

1. 确保已安装 Go 语言环境（建议 Go 1.23.4 及以上版本）。
2. 克隆项目到本地：
   ```bash
   git clone https://gitee.com/MM-Q/cbk.git
   ```
3. 进入项目目录并编译：
   ```bash
   cd cbk
   
   # linux
   sh build.sh
   
   # windows
   build.bat
   ```
4. 将生成的可执行文件 `cbk` 添加到系统 PATH 中，或直接使用 `./cbk` 运行。

## 使用说明

### help
```bash
cbk - 命令行备份任务管理工具

功能：用于管理备份任务，包括添加、运行、编辑、删除备份任务，查看任务日志，显示任务详情等。

用法：
  cbk [主参数] <子命令> [子参数]

子命令：
  list                列出所有备份任务
  run                 运行指定备份任务
  add                 添加新的备份任务
  delete              删除指定备份任务
  edit                编辑指定备份任务
  log                 查看备份任务的日志
  show                显示指定备份任务的详细信息
  unpack              解压指定版本的备份
  zip                 将指定目标路径打包为ZIP压缩包
  unzip               将指定ZIP压缩包解压到指定路径
  clear               清除数据库记录及其数据目录
  init                生成预设的自动补全脚本或配置文件
  export              导出备份任务到控制台
  version             显示当前版本信息
  help                显示帮助信息

参数：
  -h, --help                  显示帮助信息
  -v, --v                     显示简约的版本信息 
  -vv, --vv                   显示详细版本信息

支持的表格样式：
  "default":      默认样式
	"bold":         加粗样式
	"colorbright":  亮色样式
	"colordark":    暗色样式
	"double":       双边框样式
	"light":        浅色样式
	"rounded":      圆角样式
	"bd":           加粗样式(缩写版)
	"cb":           亮色样式(缩写版)
	"cd":           暗色样式(缩写版)
	"de":           双边框样式(缩写版)
	"lt":           浅色样式(缩写版)
	"ro":           圆角样式(缩写版)

示例：
  cbk list
    列出所有备份任务的基本信息。

  cbk run -id 1
    运行ID为1的备份任务。

  cbk add -n "每日备份" -t "/data" -b "/backup" -c 7
    添加一个名为“每日备份”的备份任务，目标目录为“/data”，备份文件存放在“/backup”，保留最近7个备份。

  cbk delete -id 2
    删除ID为2的备份任务。

  cbk edit -id 3 -n "周备份" -c 4
    将ID为3的备份任务名称改为“周备份”，并将保留数量设置为4。

  cbk log -l 20
    查看最近20行备份任务的日志。

  cbk show -id 4
    显示ID为4的备份任务的详细信息。

  cbk log -l 20 -v -ts "light"
    查看最近20行备份任务的详细日志，并使用“light”样式显示。

  cbk unpack -id 5 -o "/data/unpacked" -v 151nmgd1
    解压ID为5的备份文件，解压后的目录为“/data/unpacked”，解压的版本ID为151nmgd1。

  source <(cbk init -type bash)
    启用bash自动补全。
  
  cbk version
    显示当前版本信息。

注意：
- 在使用子命令时，请确保参数的正确性，避免因参数错误导致任务执行失败。
- 对于某些操作（如删除和编辑），建议在执行前仔细确认任务ID或名称，以避免误操作。
```

## 配置

- 数据库文件默认存储在用户主目录下的 `.cbk/cbk.db`。
- 备份文件默认存储在用户主目录下的 `.cbk/data/<任务名>`。

## 示例

1. 添加一个备份任务：
   ```bash
   cbk add -n mytask -t /path/to/target -b /path/to/backup -c 5
   ```

2. 运行备份任务：
   ```bash
   cbk run -id 1
   ```

3. 查看备份日志：
   ```bash
   cbk log -l 20
   ```

4. 解压指定版本的备份：
   ```bash
   cbk unpack -id 1 -v 123456 -o /path/to/output
   ```

## 依赖

- [sqlx](https://github.com/jmoiron/sqlx)：用于数据库操作。
- [go-pretty](https://github.com/jedib0t/go-pretty)：用于表格输出。
- [colorlib](https://gitee.com/MM-Q/colorlib)：用于命令行颜色渲染。 
