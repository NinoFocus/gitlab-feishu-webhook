##### Build Stage #####

# 基于golang:1.20镜像构建
FROM golang:1.20 AS build

# 设置工作目录
WORKDIR /app

# 将项目拷贝到容器中
COPY . .

# 在 GOPATH 中构建（使用了 go mod）
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o gitlab-feishu-webhook src/main.go


##### Release Stage #####

# 基于 alpine 镜像构建 Release Stage
FROM gcr.io/distroless/base-debian11 AS release

WORKDIR /

# 从 Build Stage 拷贝二进制文件
COPY --from=build /app/gitlab-feishu-webhook /gitlab-feishu-webhook

# 暴露 8083 端口
EXPOSE 8083

USER nonroot:nonroot

# 启动命令
ENTRYPOINT ["/gitlab-feishu-webhook"]