FROM golang:alpine as builder

WORKDIR /builder
COPY . .
RUN apk add upx && go mod download && go build -o app && upx -9 app
RUN ls -lh && chmod +x ./app

FROM repo.lxh.io/alpine:3.16.0 as runner
MAINTAINER LiXunHuan(lxh@cxh.cn)
LABEL org.opencontainers.image.source = "https://github.com/lixh00/web-wechat"

WORKDIR /app
ENV TZ=Asia/Shanghai
COPY --from=builder /builder/app ./app
COPY --from=builder /builder/resource/carbon-language-zh-CN.json ./resource/carbon-language-zh-CN.json
CMD ./app