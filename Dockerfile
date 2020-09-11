FROM golang:1.14

LABEL name="CertiK Chain"
LABEL maintainer="CertiK"
LABEL repository="https://github.com/certikfoundation/shentu"

RUN apt-get update && apt-get install nodejs npm -y

ADD . /shentu
WORKDIR /shentu

RUN make install

CMD ["certikd"]
