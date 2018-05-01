FROM node:9.5.0 as builder

# Install golang and set GOPATH
RUN cd /tmp/ && \
    wget https://dl.google.com/go/go1.9.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.9.4.linux-amd64.tar.gz && \
    export PATH=$PATH:/usr/local/go/bin && export GOPATH=/go && \
    go get -v github.com/restanrm/bell && \
    go get -v github.com/rakyll/statik

# Build front assets
RUN cd /go/src/github.com/restanrm/bell/front && \
    npm install && \
    npm run build

RUN cd /go/src/github.com/restanrm/bell && \
    /go/bin/statik -src=./front/dist && \
    export PATH=$PATH:/usr/local/go/bin:/go/bin && export GOPATH=/go && \
    go generate && \
    go get ./... && \
    go install ./...

FROM ubuntu:17.10

RUN apt-get update && apt-get -y install \
    alsa-base alsa-utils pulseaudio \
    flite mpv

WORKDIR /data
VOLUME /data

COPY --from=builder /go/bin/bell /bell
COPY --from=builder /go/bin/bellctl /bellctl
COPY store.json /data/store.json
COPY sounds /data/sounds

EXPOSE 10101

ENTRYPOINT ["/bell"]
