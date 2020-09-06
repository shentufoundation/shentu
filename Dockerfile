FROM golang:1.14

LABEL name="CertiK Chain"
LABEL maintainer="CertiK"
LABEL repository="https://github.com/certikfoundation/shentu"

ADD bin/dsc-linux /bin/dsc
ADD bin/solc-linux /bin/solc

RUN apt-get update && apt-get install nodejs npm -y

ADD . /shentu
WORKDIR /shentu

RUN make install

CMD ["certikd"]
