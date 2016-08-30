FROM ubuntu:15.10
MAINTAINER kc merrill <kcmerrill@gmail.com>

RUN apt-get -y update
RUN apt-get -y install curl iproute2 netbase

COPY . /automagicproxy

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["/automagicproxy/bin/automagicproxy"]
CMD ["--dockerized", "--containerized"]
