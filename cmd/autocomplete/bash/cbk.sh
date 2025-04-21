#!/bin/bash

# 定义一个名为 _cbk 的函数, 用于为 cbk 命令提供自动补全功能
_cbk() {
    # 定义局部变量
    local cur prev sub_opts opts

    # 初始化 COMPREPLY 数组, 用于存储补全选项
    COMPREPLY=()

    # 获取当前输入的单词
    cur="${COMP_WORDS[COMP_CWORD]}"

    # 获取前一个输入的单词
    prev="${COMP_WORDS[COMP_CWORD - 1]}"

    # 定义所有可用的子命令和选项
    opts="list run add delete edit log show unpack zip unzip uz clear init export version help --help -h -v -vv"

    # 根据前一个单词(prev)来决定补全的内容
    case "${prev}" in
    cbk)
        # 如果前一个单词是 cbk, 补全所有子命令和选项
        COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
        return 0
        ;;
    add)
        # 如果前一个单词是 add, 补全 add 命令的选项
        sub_opts="-n -t -b -c -d -bn -h -nc -f -ex"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    a)
        # 如果前一个单词是 a, 补全 a 命令的选项
        sub_opts="-n -t -b -c -d -bn -h -nc -f -ex"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    delete)
        # 如果前一个单词是 delete, 补全 delete 命令的选项
        sub_opts="-id -n -d -v -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    d)
        # 如果前一个单词是 d, 补全 d 命令的选项
        sub_opts="-id -n -d -v -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    edit)
        # 如果前一个单词是 edit, 补全 edit 命令的选项
        sub_opts="-id -ids -n -c -d -bn -h -nc -ex"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    e)
        # 如果前一个单词是 e, 补全 e 命令的选项
        sub_opts="-id -ids -n -c -d -bn -h -nc -ex"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    list)
        # 如果前一个单词是 list, 补全 list 命令的选项
        sub_opts="-ts -no-table -nt -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    l)
        # 如果前一个单词是 l, 补全 l 命令的选项
        sub_opts="-ts -no-table -nt -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    log)
        # 如果前一个单词是 log, 补全 log 命令的选项
        sub_opts="-l -v -ts -no-table -nt -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    run)
        # 如果前一个单词是 run, 补全 run 命令的选项
        sub_opts="-id -h -ids"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    r)
        # 如果前一个单词是 r, 补全 r 命令的选项
        sub_opts="-id -h -ids"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    show)
        # 如果前一个单词是 show, 补全 show 命令的选项
        sub_opts="-id -v -ts -no-table -nt -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    s)
        # 如果前一个单词是 s, 补全 s 命令的选项
        sub_opts="-id -v -ts -no-table -nt -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    unpack)
        # 如果前一个单词是 unpack, 补全 unpack 命令的选项
        sub_opts="-id -v -o -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    u)
        # 如果前一个单词是 u, 补全 u 命令的选项
        sub_opts="-id -v -o -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    zip)
        # 如果前一个单词是 zip, 补全 zip 命令的选项
        sub_opts="-o -t -h -nc -ex"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    z)
        # 如果前一个单词是 z, 补全 z 命令的选项
        sub_opts="-o -t -h -nc -ex"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    unzip)
        # 如果前一个单词是 unzip, 补全 unzip 命令的选项
        sub_opts="-f -d -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    uz)
        # 如果前一个单词是 uz, 补全 uz 命令的选项
        sub_opts="-f -d -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    clear)
        # 如果前一个单词是 clear, 补全 clear 命令的选项
        sub_opts="-confirm -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    init)
        # 如果前一个单词是 init, 补全 init 命令的选项
        sub_opts="-type -h"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    export)
        # 如果前一个单词是 export, 补全 export 命令的选项
        sub_opts="-id -h -all"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
        ;;
    *)
        # 如果前一个单词不匹配任何已知命令, 不做任何操作
        ;;
    esac

    # 如果前一个单词是 -id, 补全任务 ID
    if [[ ${prev} == "-id" ]]; then
        # 获取任务 ID 列表
        local task_ids=$(cbk list -nt | awk 'NR>1 {print $2}')
        COMPREPLY=($(compgen -W "${task_ids}" -- ${cur}))
        return 0
    fi

    # 如果前一个单词是 -ids, 补全任务 ID
    if [[ ${prev} == "-ids" ]]; then
        # 获取任务 ID 列表
        local task_ids=$(cbk list -nt | awk 'NR>1 {print $2}')
        COMPREPLY=($(compgen -W "${task_ids}" -- ${cur}))
        return 0
    fi

    # 如果前一个单词是 -ts, 补全表格样式
    if [[ ${prev} == "-ts" ]]; then
        # 定义所有可用的表格样式
        local table_styles="default bold colorbright colordark double light rounded bd cb cd de lt ro"
        COMPREPLY=($(compgen -W "${table_styles}" -- ${cur}))
        return 0
    fi

    # 如果前一个单词是 -type, 补全类型
    if [[ ${prev} == "-type" ]] || [[ ${prev} == "--type" ]]; then
        # 定义所有可用的类型
        local completion_types="bash addtask"
        COMPREPLY=($(compgen -W "${completion_types}" -- ${cur}))
        return 0
    fi

    # 如果前一个单词是-n, 则补全文件名和目录名
    if [[ ${prev} == "-n" ]]; then
        sub_opts="$(ls)"
        COMPREPLY=($(compgen -W "${sub_opts}" -f -d -- ${cur}))
        return 0
    fi

    # 如果前一个单词是-t, 则补全文件名和目录名
    if [[ ${prev} == "-t" ]]; then
        sub_opts="$(ls)"
        COMPREPLY=($(compgen -W "${sub_opts}" -f -d -- ${cur}))
        return 0
    fi

    # 如果前一个单词是-b, 则补全文件名和目录名
    if [[ ${prev} == "-b" ]]; then
        sub_opts="$(ls)"
        COMPREPLY=($(compgen -W "${sub_opts}" -f -d -- ${cur}))
        return 0
    fi

    # 如果前一个单词是-c, 则补全推荐的保留数
    if [[ ${prev} == "-c" ]]; then
        sub_opts="1 3 5 7 9 12"
        COMPREPLY=($(compgen -W "${sub_opts}" -- ${cur}))
        return 0
    fi

    # 如果前一个单词是-bn, 则补全文件名和目录名
    if [[ ${prev} == "-bn" ]]; then
        sub_opts="$(ls)"
        COMPREPLY=($(compgen -W "${sub_opts}" -f -d -- ${cur}))
        return 0
    fi

    # 如果前一个单词是-o, 则补全文件名和目录名
    if [[ ${prev} == "-o" ]]; then
        sub_opts="$(ls)"
        COMPREPLY=($(compgen -W "${sub_opts}" -f -d -- ${cur}))
        return 0
    fi

    # 如果前一个单词是-f, 则补全文件名和目录名
    if [[ ${prev} == "-f" ]]; then
        sub_opts="$(ls)"
        COMPREPLY=($(compgen -W "${sub_opts}" -f -d -- ${cur}))
        return 0
    fi

    # 如果前一个单词是-d, 则补全文件名和目录名
    if [[ ${prev} == "-d" ]]; then
        sub_opts="$(ls)"
        COMPREPLY=($(compgen -W "${sub_opts}" -f -d -- ${cur}))
        return 0
    fi

    # 如果前一个单词是-ex, 则提示常见的排除规则
    if [[ ${prev} == "-ex" ]]; then
        sub_opts="*.log *.txt logs log"
        COMPREPLY=($(compgen -W "${sub_opts}" -f -d -- ${cur}))
        return 0
    fi
}

# 将 _cbk 函数与 cbk 命令关联, 实现自动补全功能
complete -F _cbk cbk
