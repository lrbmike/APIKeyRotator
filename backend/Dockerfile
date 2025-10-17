FROM golang:1.24-alpine AS builder

# 设置Go模块代理，加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./api-key-rotator .

FROM alpine:3.22.0

RUN apk add --no-cache ca-certificates tzdata

RUN mkdir /app

COPY --from=builder /app/api-key-rotator /app/api-key-rotator

WORKDIR /app

EXPOSE 8000

ENV TZ=Asia/Shanghai

RUN cp /usr/share/zoneinfo/${TZ} /etc/localtime && echo "${TZ}" > /etc/timezone

CMD ["./api-key-rotator"]