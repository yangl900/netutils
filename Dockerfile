FROM golang:1.9.2 as builder_stdnet
WORKDIR  /go/src/stdnet/
COPY ./stdnet /go/src/stdnet/
RUN go test -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -o stdnet .

FROM golang:1.9.2 as builder_start
WORKDIR  /go/src/start/
COPY ./start /go/src/start/
RUN go test -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -o start .

FROM debian:jessie

RUN apt-get update \
    && apt-get install -y \
        traceroute \
        curl \
        dnsutils \
        netcat-openbsd \
        jq \
        nmap \ 
        net-tools \
        openssh-client \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder_stdnet /go/src/stdnet/stdnet /
COPY --from=builder_start /go/src/start/start /

CMD [ "/bin/bash", "-c", "tail -f /dev/null" ]