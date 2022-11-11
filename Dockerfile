FROM golang as build

# copy all codes
COPY ./ ./goim
# set go env
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR ./goim
RUN chmod +x ./script/*.sh
# run build script
RUN /bin/sh ./script/build_all.sh

FROM ubuntu
# set timezone
RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && dpkg-reconfigure -f noninteractive tzdata
# copy built target binary from build image
COPY --from=build /goim/script ./goim/script
COPY --from=build /goim/target ./goim/target
# volume for config files
VOLUME "/goim/config"
WORKDIR ./goim

CMD["./target/api", "--config=config.yaml"]
