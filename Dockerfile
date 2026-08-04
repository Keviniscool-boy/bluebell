FROM golang:1.25-alpine AS builder
    WORKDIR /app
    # 国内环境换 Go 模块代理，否则 proxy.golang.org 连不上
    ENV GOPROXY=https://goproxy.cn,direct
    COPY go.mod go.sum ./
    RUN go mod download
    COPY . .
    RUN CGO_ENABLED=0 go build -o bluebell_app .

#阶段二是把编译好的二进制文件拷贝到一个更小的镜像中
FROM alpine:3.20
WORKDIR /app
# 装 netcat：wait-for.sh 里的 nc -z 依赖它，alpine 默认不带
RUN apk add --no-cache netcat-openbsd
COPY --from=builder /app/bluebell_app .
COPY settings/config.docker.yaml ./config.yaml
COPY wait-for.sh .
RUN chmod +x wait-for.sh
EXPOSE 8081
