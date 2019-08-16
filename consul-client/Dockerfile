FROM ubuntu:16.04

# Install.
RUN \
  sed -i 's/# \(.*multiverse$\)/\1/g' /etc/apt/sources.list && \
  apt-get update && \
  apt-get -y upgrade && \
  apt-get install -y build-essential && \
  apt-get install -y software-properties-common && \
  apt-get install -y byobu curl git htop man unzip vim wget nano rsyslog && \
  curl https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz | tar xzf - -C / && \
  mv /go /goroot && \
  mkdir /app && \
  rm -rf /var/cache/apk/* && \
  rm -rf /var/lib/apt/lists/*

ADD consul /usr/local/bin/

RUN consul --version

RUN mkdir -p /consul/data && \
    mkdir -p /consul/config && \
    mkdir -p /scripts && \
    ln -sf /var/run/log /var/run/syslog

ADD scripts/. /scripts

ADD config.json /consul/config

ENTRYPOINT ["/scripts/run-script.sh"]

EXPOSE 7080 8300 8301 8301/udp 8302 8302/udp 8500 8600 8600/udp 8000
