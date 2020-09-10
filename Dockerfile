FROM golang:1.14

WORKDIR /go/src/github.com/certikfoundation/shentu

COPY go.* /go/src/github.com/certikfoundation/shentu/

LABEL name="CertiK Chain"
LABEL maintainer="CertiK"
LABEL repository="https://github.com/certikfoundation/shentu"

RUN apt-get update && apt-get install nodejs npm -y

WORKDIR /go/src/github.com/certikfoundation/shentu

ADD . /shentu
WORKDIR /shentu

RUN make build

WORKDIR /root
COPY /go/src/github.com/certikfoundation/shentu/build/certikd /usr/local/bin/certikd
COPY /go/src/github.com/certikfoundation/shentu/build/certikcli /usr/local/bin/certikcli

CMD ["certikd"]
