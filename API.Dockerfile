# Builder image
FROM golang:1.17-alpine as builder

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories && \
    apk add --no-cache \
    wget \
    git

RUN mkdir -p /home/build && \
    mkdir -p /home/app

ARG build_dir=/home/build
ARG app_dir=/home/app

ENV ServiceName=gf.bridgx.api

WORKDIR $build_dir

COPY . .

# Cache dependencies
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

COPY go.mod go.mod
COPY go.sum go.sum
#RUN  go mod download

RUN mkdir -p output/conf output/bin

# detect mysql start
COPY wait-for-it.sh output/bin/wait-for-it.sh

RUN find conf/ -type f ! -name "*_local.*" | xargs -I{} cp {} output/conf/
RUN cp scripts/run_api.sh output/

RUN CGO_ENABLED=0 GO111MODULE=on go build -o output/bin/${ServiceName} ./cmd/api

RUN cp -rf output/* $app_dir

# --------------------------------------------------------------------------------- #
# Executable image
FROM alpine:3.14

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories && \
    apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone
ENV TZ Asia/Shanghai

RUN apk update \
        && apk upgrade \
        && apk add --no-cache bash \
        bash-doc \
        bash-completion \
        && rm -rf /var/cache/apk/* \
        && /bin/bash

ENV ServiceName=gf.bridgx.api
ENV SpecifiedConfig=prod

COPY --from=builder /home/app /home/tiger/app
RUN addgroup -S tiger && adduser -S tiger -G tiger
WORKDIR /home/tiger/app
RUN chown -R tiger:tiger /home/tiger && chmod +x run_api.sh && chmod +x bin/wait-for-it.sh

USER tiger
EXPOSE 9090
#CMD ["/bin/sh","/home/tiger/app/run_api.sh"]