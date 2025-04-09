#!/bin/bash

# 定义一个名为 _cbk 的函数，用于为 cbk 命令提供自动补全功能
_cbk()
{
    # 定义局部变量
    local cur prev opts

    # 初始化 COMPREPLY 数组，用于存储补全选项
    COMPREPLY=()

    # 获取当前输入的单词
    cur="${COMP_WORDS[COMP_CWORD]}"

    # 获取前一个输入的单词
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # 定义所有可用的子命令和选项
    opts="list run add delete edit log show unpack zip unzip clear version help --help -h -v -vv"

    # 根据前一个单词（prev）来决定补全的内容
    case "${prev}" in
        cbk)
            # 如果前一个单词是 cbk，补全所有子命令和选项
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        add)
            # 如果前一个单词是 add，补全 add 命令的选项
            opts="-n -t -b -k -bn"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        delete)
            # 如果前一个单词是 delete，补全 delete 命令的选项
            opts="-id -n -d -v"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        edit)
            # 如果前一个单词是 edit，补全 edit 命令的选项
            opts="-id -n -k -bn"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        list)
            # 如果前一个单词是 list，补全 list 命令的选项
            opts="-ts -no-table -nt"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        log)
            # 如果前一个单词是 log，补全 log 命令的选项
            opts="-l -v -ts -no-table -nt"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        run)
            # 如果前一个单词是 run，补全 run 命令的选项
            opts="-id"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        show)
            # 如果前一个单词是 show，补全 show 命令的选项
            opts="-id -v -ts -no-table -nt"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        unpack)
            # 如果前一个单词是 unpack，补全 unpack 命令的选项
            opts="-id -v -o"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        zip)
            # 如果前一个单词是 zip，补全 zip 命令的选项
            opts="-o -t"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        unzip)
            # 如果前一个单词是 unzip，补全 unzip 命令的选项
            opts="-f -d"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        clear)
            # 如果前一个单词是 clear，补全 clear 命令的选项
            opts="-confirm"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        complete)
            # 如果前一个单词是 complete，补全 complete 命令的选项
            opts="-type --type"
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        *)
            # 如果前一个单词不匹配任何已知命令，不做任何操作
            ;;
    esac

    # 如果前一个单词是 -id，补全任务 ID
    if [[ ${prev} == "-id" ]]; then
        # 获取任务 ID 列表
        local task_ids=$(cbk list -nt | awk '{print $1}' | grep -v "ID")
        COMPREPLY=( $(compgen -W "${task_ids}" -- ${cur}) )
        return 0
    fi

    # 如果前一个单词是 -ts，补全表格样式
    if [[ ${prev} == "-ts" ]]; then
        # 定义所有可用的表格样式
        local table_styles="default bold colorbright colordark double light rounded bd cb cd de lt ro"
        COMPREPLY=( $(compgen -W "${table_styles}" -- ${cur}) )
        return 0
    fi

    # 如果前一个单词是 -type，补全补全类型
    if [[ ${prev} == "-type" ]] || [[ ${prev} == "--type" ]]; then
        # 定义所有可用的补全类型
        local completion_types="bash"
        COMPREPLY=( $(compgen -W "${completion_types}" -- ${cur}) )
        return 0
    fi
}

# 将 _cbk 函数与 cbk 命令关联，实现自动补全功能
complete -F _cbk cbk