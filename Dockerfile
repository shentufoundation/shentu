FROM golang:1.15

LABEL name="CertiK Chain"
LABEL maintainer="CertiK"
LABEL repository="https://github.com/certikfoundation/shentu"
LABEL org.opencontainers.image.source=https://github.com/certikfoundation/shentu

EXPOSE 26656 26657 26658

RUN apt-get update && apt-get install nodejs npm -y

ADD . /shentu
WORKDIR /shentu

RUN make install

CMD ["certik"]
