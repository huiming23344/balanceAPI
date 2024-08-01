FROM golang:latest as builder

# 设置工作目录
WORKDIR /src

# 复制源代码到容器中
COPY . .

# 编译二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o balanceapi .

FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从builder镜像中复制编译好的二进制文件
COPY --from=builder /src/balanceapi /app/balanceapi
COPY conf/app-docker.yaml /app/app.yaml

EXPOSE 20001

# 设置执行权限
RUN chmod +x /app/balanceapi

# 运行二进制文件
CMD ["./balanceapi"]