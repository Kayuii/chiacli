FROM golang:alpine as builder

RUN apk update \
  && apk --no-cache add --virtual build-dependencies \
  zlib-dev build-base linux-headers coreutils

ENV GOPROXY=https://goproxy.io,direct

COPY . /opt/
WORKDIR /opt

RUN cd /opt \
  && ls -al \
  && make build-static

FROM kayuii/chia-plotter:chia-v1.1.5
COPY --from=builder /opt/chiacli-static /root/chiacli
