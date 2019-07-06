FROM node:9.5.0 as builder

ARG GOVERSION=1.12.0

# Install golang and set GOPATH
RUN cd /tmp/ && \
    wget https://dl.google.com/go/go${GOVERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOVERSION}.linux-amd64.tar.gz && \
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

# create bare server
FROM alpine:3.9 as server
COPY --from=builder /go/bin/bell /bin/bell
COPY data/store.json /data/store.json
COPY data/sounds /data/sounds
EXPOSE 10101
CMD ["bell","-d","/data"]

# create bare client
FROM alpine:3.9 as client
COPY --from=builder /go/bin/bellctl /usr/local/bin/bellctl

## create full server that can play sounds
FROM ubuntu:17.10 as full

RUN apt-get update && apt-get -y install \
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
