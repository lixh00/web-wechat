FROM golang:alpine as builder
MAINTAINER LiXunHuan(lxh@cxh.cn)
# 创建工作目录，修改alpine源为中科大的源，安装必要工具
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
  apk update && \
  apk upgrade && \
  apk add ca-certificates gcc g++ && update-ca-certificates && \
  apk add --update tzdata && \
  rm -rf /var/cache/apk/*
ENV TZ=Asia/Shanghai
WORKDIR /builder
COPY . .
RUN go mod download && go build -o wechat
RUN ls -lh && chmod +x ./wechat
FROM golang:alpine as runner
MAINTAINER LiXunHuan(lxh@cxh.cn)
WORKDIR /app
COPY --from=builder /builder/wechat ./wechat
CMD ./wechat