# FROM golang:1.15
FROM golang:alpine3.13 AS build-env

# Set up dependencies
ENV PACKAGES bash curl make git libc-dev gcc linux-headers eudev-dev python3

# ADD . /shentu
WORKDIR /shentu

COPY go.mod .
COPY go.sum .

COPY . .

RUN apk add --no-cache $PACKAGES && make install

FROM alpine:edge

LABEL name="CertiK Chain"
LABEL maintainer="CertiK"
LABEL repository="https://github.com/certikfoundation/shentu"
LABEL org.opencontainers.image.source=https://github.com/certikfoundation/shentu

RUN apk add --update ca-certificates

WORKDIR /shentu

COPY --from=build-env /go/bin/certik /usr/bin/certik

CMD ["certik"]
