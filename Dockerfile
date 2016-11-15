FROM golang:1.7
MAINTAINER kc merrill <kcmerrill@gmail.com>

RUN apt-get -y update
RUN apt-get -y install curl iproute2 netbase

RUN go get -u github.com/kcmerrill/automagicproxy

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["automagicproxy"]
CMD ["--dockerized", "--containerized"]
