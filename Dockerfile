FROM golang:1.21 AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn

WORKDIR /build
COPY . .
RUN go mod tidy && go build -o /usr/local/bin/app main.go
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN apt-get update && apt-get -y install sqlite3 && apt-get -y install libsqlite3-dev
EXPOSE 8011
CMD ["app"]




