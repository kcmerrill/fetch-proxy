FROM ubuntu:15.10
MAINTAINER kc merrill <kcmerrill@gmail.com>

RUN apt-get -y update
RUN apt-get -y install curl iproute2

COPY . /automagicproxy

CMD ["/automagicproxy/automagicproxy", "--dockerized", "--containerized"]
