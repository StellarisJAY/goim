FROM golang

COPY ./ ./goim

# set go env
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR ./goim

RUN /bin/sh ./script/build_all.sh