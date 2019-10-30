FROM ubuntu:latest

# Install.
RUN \
  sed -i 's/# \(.*multiverse$\)/\1/g' /etc/apt/sources.list && \
  apt-get update && \
  apt-get -y upgrade && \
  apt-get install -y build-essential && \
  apt-get install -y software-properties-common && \
  apt-get install -y byobu curl && \
  curl https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz | tar xzf - -C / && \
  mv /go /goroot && \
  rm -rf /var/cache/apk/* && \
  rm -rf /var/lib/apt/lists/*

# Set environment variables.
ENV GOROOT=/goroot \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PATH=${PATH}:/goroot/bin:/gopath/bin

WORKDIR app
COPY . /app/

COPY config.json /var/log/config.json

CMD ["/app/app"]

EXPOSE 80 7080 6080
