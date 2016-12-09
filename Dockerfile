FROM golang:latest
ADD consul-lock /usr/local/bin
WORKDIR /root
ENTRYPOINT ["/usr/local/bin/consul-lock"]
