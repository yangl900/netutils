FROM golang:1.9.2 as builder_stdnet
WORKDIR  /go/src/stdnet/
COPY ./stdnet /go/src/stdnet/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o stdnet .

FROM golang:1.9.2 as builder_start
WORKDIR  /go/src/start/
COPY ./start /go/src/start/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o start .

FROM golang:1.9.2 as builder_echo
WORKDIR  /go/src/tcp-echo/
COPY ./tcp-echo /go/src/tcp-echo/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o tcp-echo .

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
COPY --from=builder_echo /go/src/tcp-echo/tcp-echo /

CMD [ "/bin/bash", "-c", "tail -f /dev/null" ]