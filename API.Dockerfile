# Builder image
FROM golang:1.17-alpine as builder

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories && \
    apk add --no-cache \
    wget \
    git

RUN mkdir -p /home/build && \
    mkdir -p /home/api

ARG build_dir=/home/build
ARG api_dir=/home/api

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
COPY wait-for-api.sh output/bin/wait-for-api.sh

RUN find conf/ -type f ! -name "*local*" -print0 | xargs -0 -I{} cp {} output/conf/
RUN cp scripts/run_api.sh output/

RUN CGO_ENABLED=0 GO111MODULE=on go build -o output/bin/${ServiceName} ./cmd/api

RUN cp -rf output/* $api_dir

# --------------------------------------------------------------------------------- #
# Executable image
FROM alpine:3.14

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories

RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone
ENV TZ Asia/Shanghai

RUN apk add --no-cache bash

ENV ServiceName=gf.bridgx.api
ENV SpecifiedConfig=prod

COPY --from=builder /home/api /home/tiger/api
RUN addgroup -S tiger && adduser -S tiger -G tiger
WORKDIR /home/tiger/api
RUN chown -R tiger:tiger /home/tiger && chmod +x run_api.sh && chmod +x bin/wait-for-api.sh

USER tiger
EXPOSE 9090
CMD ["/bin/sh","/home/tiger/api/run_api.sh"]