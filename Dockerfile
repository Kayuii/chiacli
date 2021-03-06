FROM golang:alpine as builder

RUN apk update \
  && apk --no-cache add --virtual build-dependencies \
  zlib-dev build-base linux-headers coreutils

ENV GOPROXY=https://goproxy.io,direct

COPY . /opt/
WORKDIR /opt

RUN cd /opt \
  && ls -alh \
  && make build-mini

# FROM kayuii/chia-plotter:chia-v1.1.7
# COPY --from=builder /opt/chiacli-static /root/chiacli
