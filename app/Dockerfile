FROM ubuntu:16.04

# Install.
RUN \
  sed -i 's/# \(.*multiverse$\)/\1/g' /etc/apt/sources.list && \
  apt-get update && \
  apt-get -y upgrade && \
  apt-get install -y build-essential && \
  apt-get install -y software-properties-common && \
  apt-get install -y byobu curl git htop man unzip vim wget nano && \
  curl https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz | tar xzf - -C / && \
  mv /go /goroot && \
  rm -rf /var/cache/apk/* && \
  rm -rf /var/lib/apt/lists/*

# Set environment variables.
ENV GOROOT=/goroot \
    GOPATH=/gopath \
    GOBIN=/gopath/bin \
    PATH=${PATH}:/goroot/bin:/gopath/bin

# Define working directory.
WORKDIR /root


# Define default command.
CMD ["sh", "-c", "tail -f /dev/null"]

EXPOSE 80 7080 6080