FROM golang:1.9.2 as builder
WORKDIR  /go/src/ncssh/
COPY . /go/src/ncssh/
RUN go test -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ncssh .

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

COPY --from=builder /go/src/ncssh/ncssh /
COPY ./ssh /
RUN chmod +x /ssh

CMD [ "/bin/bash", "-c", "sleep 10000000000000" ]