FROM nimmis/alpine-glibc
MAINTAINER nimmis <kjell.havneskold@gmail.com>

RUN apk update && apk upgrade && \

    # Make info file about this build
    mkdir -p /etc/BUILDS/ && \
    printf "Build of nimmis/alpine-golang:1.7.1, date: %s\n"  `date -u +"%Y-%m-%dT%H:%M:%SZ"` > /etc/BUILDS/alpine-golang && \

    # add extra package for installation  
    apk add curl && \

    #apk add rsyslog && \

    # download and upack golang
    curl https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz | tar xzf - -C / && \
    mv /go /goroot && \
 
    # remove extra packages

    # remove cached info
    rm -rf /var/cache/apk/*

ENV GOROOT=/goroot \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PATH=${PATH}:/goroot/bin:/gopath/bin

ADD consul /usr/local/bin/

RUN consul --version

RUN mkdir -p /consul/data && \
    mkdir -p /consul/config && \
    mkdir -p /scripts

ADD scripts/. /scripts

ADD config.json /consul/config

EXPOSE 7080 8300 8301 8301/udp 8302 8302/udp 8500 8600 8600/udp 8000

#RUN consul agent --config-dir=/consul/config -retry-join 172.17.0.10
