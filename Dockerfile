FROM golang:1.16-alpine

# 在容器中创建工作目录
WORKDIR /app
# 设置代理
ENV GOPROXY=https://goproxy.cn,direct
# 将当前目录下的所有文件复制到工作目录
COPY . .

# 安装依赖
RUN go mod download

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 编译Go应用
RUN go build ./cmd/cli/main.go -o app

# 运行应用程序
CMD ["./app"]
