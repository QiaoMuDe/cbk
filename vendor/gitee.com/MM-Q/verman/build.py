import os
import subprocess
import sys
import argparse
import time
import zipfile
from datetime import datetime, timezone
import platform


# 可根据实际情况修改区域 #
# 项目名称|包名，可根据实际情况修改(在注入git信息时会用到)
PROJECT_NAME = "your_project_name"
# 基础输出文件名(不带扩展名)
BASE_OUTPUT_NAME = "myapp"
# 默认入口文件的位置
DEFAULT_ENTRY_FILE = "./main.go"
# 默认构建时的链接器标志
DEFAULT_LDFLAGS = "-s -w"
# ldflags模板字符串
LD_FLAGS_TEMPLATE = "-X 'gitee.com/MM-Q/verman.appName={app_name}' -X 'gitee.com/MM-Q/verman.gitVersion={git_version}' -X 'gitee.com/MM-Q/verman.gitCommit={git_commit}' -X 'gitee.com/MM-Q/verman.gitCommitTime={commit_time}' -X 'gitee.com/MM-Q/verman.buildTime={build_time}' -X 'gitee.com/MM-Q/verman.gitTreeState={tree_state}' -s -w"
# 默认的 Go 编译器，使用全局 PATH 中的 go
DEFAULT_GO_COMPILER = "go"
# 默认不使用 vendor 克隆依赖
DEFAULT_USE_VENDOR = False
# 是否注入git信息，默认为True
DEFAULT_INJECT_GIT_INFO = True


# 默认无需修改区域 #
# 定义颜色转义字符
RED_BOLD = '\033[1;31m'  # 红色加粗
GREEN_BOLD = '\033[1;32m'  # 绿色加粗
RESET = '\033[0m'  # 重置颜色


# 函数定义 #
def print_success(message):
    """打印成功信息"""
    print(f"{GREEN_BOLD}ok: {message}{RESET}")

def print_error(message):
    """打印错误信息"""
    print(f"{RED_BOLD}error: {message}{RESET}")
    
def check_go_installed(go_compiler):
    """检查指定的 Go 编译器是否可用"""
    try:
        subprocess.run([go_compiler, "version"], capture_output=True, text=True, check=True)
        print_success(f"{go_compiler} 编译器已安装。")
        return True
    except subprocess.CalledProcessError:
        print_error(f"未检测到 {go_compiler} 编译器，请确保已安装并添加到 PATH 中，或者指定正确的路径。")
        return False

def check_go_mod_file():
    """检查当前目录是否存在 go.mod 文件"""
    if os.path.exists("go.mod"):
        print_success("go.mod 文件存在。")
        return True
    else:
        print_error("当前目录下未找到 go.mod 文件。")
        return False

def check_entry_file(entry_file):
    """检查指定的入口文件是否存在"""
    if os.path.exists(entry_file):
        print_success(f"入口文件 {entry_file} 存在。")
        return True
    else:
        print_error(f"未找到入口文件 {entry_file}。")
        return False

def run_go_mod_vendor(go_compiler):
    """执行 go mod vendor 克隆依赖"""
    try:
        subprocess.run([go_compiler, "mod", "vendor"], capture_output=True, text=True, check=True)
        print_success("go mod vendor 执行成功，依赖已克隆到 vendor 目录。")
    except subprocess.CalledProcessError as e:
        print_error("go mod vendor 执行失败：")
        print_error(e.stderr.strip())
        sys.exit(1)

def run_go_mod_tidy(go_compiler, use_vendor):
    """执行 go mod tidy 并处理输出"""
    command = [go_compiler, "mod", "tidy"]
    if use_vendor:
        command.extend(["-v"])
    try:
        subprocess.run(command, capture_output=True, text=True, check=True)
        print_success("go mod tidy 执行成功。")
    except subprocess.CalledProcessError as e:
        print_error("go mod tidy 执行失败：")
        print_error(e.stderr.strip())
        sys.exit(1)

def run_go_vet(go_compiler, use_vendor):
    """执行 go vet 检查代码并处理输出"""
    command = [go_compiler, "vet"]
    command.append("./...")
    try:
        subprocess.run(command, capture_output=True, text=True, check=True)
        print_success("go vet 检查成功。")
    except subprocess.CalledProcessError as e:
        print_error("go vet 检查失败：")
        print_error(e.stderr.strip())
        sys.exit(1)

def run_gofmt(go_compiler):
    """执行 gofmt -w . 格式化代码"""
    try:
        subprocess.run([go_compiler, "fmt", "./..."], capture_output=True, text=True, check=True)
        print_success("代码格式化成功。")
    except subprocess.CalledProcessError as e:
        print_error("代码格式化失败：")
        print_error(e.stderr.strip())
        sys.exit(1)

def build_go_app(go_compiler, output_file, entry_file, ldflags, use_vendor):
    """组装并执行构建命令"""
    command = [go_compiler, "build", "-o", output_file, "-ldflags", ldflags]
    if use_vendor:
        command.extend(["-mod=vendor"])
    command.append(entry_file)
    try:
        # 使用指定的链接器标志进行构建
        subprocess.run(
            command,
            capture_output=True,
            text=True,
            check=True
        )
        print_success(f"构建成功，输出文件：{output_file}。")
        return True
    except subprocess.CalledProcessError as e:
        print_error("构建失败：")
        print_error(e.stderr.strip())
        return False

def zip_executable(output_file, zip_file):
    """将构建成功的可执行文件打包到 ZIP 文件中"""
    try:
        with zipfile.ZipFile(zip_file, 'w', zipfile.ZIP_DEFLATED) as zipf:
            zipf.write(output_file)
        print_success(f"成功将 {output_file} 打包到 {zip_file} 中。")
    except Exception as e:
        print_error(f"打包到 {zip_file} 失败：{str(e)}")

def get_git_info():
    """获取 Git 版本信息"""
    try:
        git_version = subprocess.run(["git", "describe", "--tags", "--always", "--dirty"], capture_output=True, text=True, check=True, timeout=10).stdout.strip()
        git_commit = subprocess.run(["git", "rev-parse", "--short", "HEAD"], capture_output=True, text=True, check=True, timeout=10).stdout.strip()
        git_commit_time = subprocess.run(["git", "log", "-1", "--format=%cd", "--date=iso"], capture_output=True, text=True, check=True, timeout=10).stdout.strip()
        git_status = "dirty" if subprocess.run(["git", "status", "--porcelain"], capture_output=True, text=True, timeout=10).stdout.strip() else "clean"
        format_time = datetime.strptime(git_commit_time, '%Y-%m-%d %H:%M:%S %z').strftime('%Y-%m-%dT%H:%M:%SZ')
        return git_version, git_commit, format_time, git_status
    except subprocess.CalledProcessError:
        print_error("警告: 无法获取 Git 版本信息，可能是当前目录不是 Git 仓库。")
        return None
    except subprocess.TimeoutExpired:
        print_error("获取 Git 版本信息超时。")
        return None

def generate_output_file_name(base_name, system):
    """根据操作系统生成默认输出文件名"""
    if system == "windows":
        return f"{base_name}.exe"
    return base_name

def generate_zip_file_name(project_name, system, architecture):
    """根据项目名称、操作系统和架构生成默认的 zip 文件名"""
    return f"{project_name}_{system}_{architecture}.zip"

def pre_build_checks(go_compiler, entry_file, use_vendor):
    """构建前的检查工作"""
    if not check_go_installed(go_compiler):
        return False
    if not check_go_mod_file():
        return False
    if not check_entry_file(entry_file):
        return False
    if use_vendor:
        try:
            run_go_mod_vendor(go_compiler)
        except Exception as e:
            print_error(str(e))
            return False
    try:
        run_go_mod_tidy(go_compiler, use_vendor)
        run_go_vet(go_compiler, use_vendor)
        run_gofmt(go_compiler)
    except Exception as e:
        print_error(str(e))
        return False
    return True

def parse_arguments():
    """解析命令行参数"""
    parser = argparse.ArgumentParser(description="构建 Go 应用程序")
    parser.add_argument("-o", "--output", help="指定输出文件名", default=None)
    parser.add_argument("-e", "--entry", help="指定入口文件路径", default=DEFAULT_ENTRY_FILE)
    parser.add_argument("-l", "--ldflags", help="指定构建时的链接器标志", default=DEFAULT_LDFLAGS)
    parser.add_argument("-g", "--go-compiler", help="指定 Go 编译器的绝对路径", default=DEFAULT_GO_COMPILER)
    parser.add_argument("-v", "--use-vendor", action="store_true", help="是否通过 vendor 克隆依赖到项目目录下", default=DEFAULT_USE_VENDOR)
    parser.add_argument("-z", "--zip", action="store_true", help="是否将构建成功的可执行文件打包到 ZIP 文件中", default=False)
    parser.add_argument("-p", "--project-name", help="指定项目名称", default=PROJECT_NAME)
    parser.add_argument("--zip-file", help="指定打包输出的 ZIP 文件名", default=None)
    parser.add_argument("-git", action="store_true", help="在构建时注入 Git 仓库的版本信息", default=DEFAULT_INJECT_GIT_INFO)
    args = parser.parse_args() # 解析命令行参数
    return args

# 主程序入口 # 
def main():
    # 记录开始时间
    start_time = time.time()

    # 解析命令行参数
    args = parse_arguments()

    entry_file = args.entry # 指定入口文件路径
    ldflags = args.ldflags # 指定构建时的链接器标志
    go_compiler = args.go_compiler # 指定 Go 编译器的绝对路径
    use_vendor = args.use_vendor # 是否通过 vendor 克隆依赖到项目目录下
    zip_flag = args.zip # 是否将构建成功的可执行文件打包到 ZIP 文件中
    inject_git = args.git # 是否注入 Git 信息
    project_name = args.project_name # 指定项目名称
    
    
    # 获取操作系统和架构信息
    system = platform.system().lower()
    architecture = platform.machine().lower()
    
    # 生成默认的 zip 文件名
    if args.zip_file is None:
        zip_file = generate_zip_file_name(project_name, system, architecture)
    else:
        zip_file = args.zip_file

    # 生成默认的输出文件名
    if args.output is None:
        base_name = f"{BASE_OUTPUT_NAME}_{system}_{architecture}"
        output_file = generate_output_file_name(base_name, system)
    else:
        output_file = args.output        

    print("开始构建 Go 应用程序...")
    
    # 执行构建前的检查工作
    if not pre_build_checks(go_compiler, entry_file, use_vendor):
        sys.exit(1)

    # 获取构建时间
    build_time = datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')

    # 根据参数注入 Git 信息
    if inject_git:
        git_info = get_git_info()
        if git_info is None:
            sys.exit(1)
        # 获取 Git 版本信息
        git_version, git_commit, format_time, git_status = git_info
        # 注入 Git 信息到链接器标志
        ldflags = LD_FLAGS_TEMPLATE.format(
            app_name=project_name,
            git_version=git_version,
            git_commit=git_commit,
            commit_time=format_time,
            build_time=build_time,
            tree_state=git_status
        )
        print_success(f"注入 Git 信息: {git_version}, {git_commit}, {format_time}, {git_status}")

    # 执行构建命令
    build_result = build_go_app(go_compiler, output_file, entry_file, ldflags, use_vendor)

    if build_result:
        print_success("构建完成。")
        if zip_flag:
            zip_executable(output_file, zip_file)
    else:
        print_error("构建失败，请检查错误信息。")

    # 记录结束时间
    end_time = time.time()
    # 计算构建耗时
    elapsed_time = end_time - start_time
    print(f"本次构建耗时: {elapsed_time:.2f} 秒")

if __name__ == "__main__":
    main()
    