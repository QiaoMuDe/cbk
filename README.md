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

### 全局命令

- `-v`：显示简单版本信息。
- `-vv`：显示详细版本信息。
- `-h` 或 `--help`：显示帮助信息。

### 子命令

#### 1. 添加任务 (`add` 或 `a`)
```bash
cbk add -n <任务名> -t <目标目录> [-b <备份目录>] [-k <保留数量>]
```
- `-n`：任务名（必填）。
- `-t`：目标目录（必填）。
- `-b`：备份目录（可选，默认为用户主目录下的 `.cbk/data/<任务名>`）。
- `-k`：保留数量（可选，默认为 3）。

#### 2. 删除任务 (`delete` 或 `d`)
```bash
cbk delete -id <任务ID> [-d]
```
或
```bash
cbk delete -n <任务名> [-d]
```
- `-id`：任务 ID（与 `-n` 二选一）。
- `-n`：任务名（与 `-id` 二选一）。
- `-d`：是否同时删除备份文件（可选）。

#### 3. 编辑任务 (`edit` 或 `e`)
```bash
cbk edit -id <任务ID> [-n <新任务名>] [-k <新保留数量>]
```
- `-id`：任务 ID（必填）。
- `-n`：新任务名（可选）。
- `-k`：新保留数量（可选）。

#### 4. 运行任务 (`run` 或 `r`)
```bash
cbk run -id <任务ID>
```
- `-id`：任务 ID（必填）。

#### 5. 查看任务列表 (`list` 或 `l`)
```bash
cbk list
```

#### 6. 查看日志 (`log`)
```bash
cbk log [-l <行数>]
```
- `-l`：显示的行数（可选，默认为 10）。

#### 7. 查看任务详情 (`show` 或 `s`)
```bash
cbk show -id <任务ID>
```
- `-id`：任务 ID（必填）。

#### 8. 解压备份 (`unpack` 或 `u`)
```bash
cbk unpack -id <任务ID> -v <版本ID> [-o <输出目录>]
```
- `-id`：任务 ID（必填）。
- `-v`：版本 ID（必填）。
- `-o`：输出目录（可选，默认为当前目录）。

#### 9. 查看版本信息 (`version`)
```bash
cbk version
```

#### 10. 查看帮助 (`help`)
```bash
cbk help
```

## 配置

- 数据库文件默认存储在用户主目录下的 `.cbk/cbk.db`。
- 备份文件默认存储在用户主目录下的 `.cbk/data/<任务名>`。

## 示例

1. 添加一个备份任务：
   ```bash
   cbk add -n mytask -t /path/to/target -b /path/to/backup -k 5
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
