FROM node:9.5.0 as builder

ARG GOVERSION=1.17.5
ARG VERSION=v0.2.3

# Install golang and set GOPATH
RUN cd /tmp/ && \
    wget https://dl.google.com/go/go${GOVERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOVERSION}.linux-amd64.tar.gz && \
    export PATH=$PATH:/usr/local/go/bin && export GOPATH=/go && \
    go get -d github.com/restanrm/bell && \
    go install github.com/rakyll/statik@latest

# Build front assets
RUN \
    cd /go/pkg/mod/github.com/restanrm/bell@$VERSION/front && \
    npm install && \
    npm run build

RUN \
    cd /go/pkg/mod/github.com/restanrm/bell@$VERSION && \
    /go/bin/statik -src=./front/dist && \
    export PATH=$PATH:/usr/local/go/bin:/go/bin && export GOPATH=/go && \
    go generate && \
    go get ./... && \
    go install ./...

FROM ubuntu:20.04

RUN apt-get update && DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt-get -y install \
    alsa-base alsa-utils pulseaudio \
    flite mpv

WORKDIR /data
VOLUME /data

COPY --from=builder /go/bin/bell /usr/local/bin/bell
COPY --from=builder /go/bin/bellctl /usr/local/bin/bellctl
COPY data/store.json /data/store.json
COPY data/sounds /data/sounds

EXPOSE 10101

CMD ["bell","-d","/data"]
