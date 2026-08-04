#!/bin/sh

# 等待所有 host:port 就绪后，再执行 -- 后面的命令
# 用法: ./wait-for.sh mysql8019:3306 redis507:6379 -- ./bluebell_app

while [ $# -gt 0 ]; do
    # 遇到 -- 就停：-- 之前是等待列表，-- 之后是真正要执行的命令
    [ "$1" = "--" ] && shift && break

    # 把 "mysql8019:3306" 拆成 host=mysql8019, port=3306
    host="${1%%:*}"
    port="${1##*:}"
    shift

    echo "Waiting for $host:$port ..."
    # 端口没通就每秒重试，直到 nc 返回成功
    until nc -z "$host" "$port" 2>/dev/null; do
        sleep 1
    done
    echo "$host:$port is ready"
done

# 执行真正的命令（如 ./bluebell_app）
exec "$@"
