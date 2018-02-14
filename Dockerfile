FROM node:9.5.0 as builder
RUN git clone https://github.com/restanrm/bell /bell && \
    cd /bell && \
    git checkout test/frontApplication && \
    cd /bell/front && \
    npm install && \
    npm run build

FROM ubuntu:17.10

RUN apt-get update && apt-get -y install \
    alsa-base alsa-utils pulseaudio \
    golang git \
    flite mpv
RUN export GOPATH=/go && \
    go get github.com/restanrm/bell

WORKDIR /data
VOLUME /data

COPY --from=builder /bell/front/dist /data/front/dist
COPY store.json /data/store.json
COPY sounds /data/sounds

EXPOSE 8080

ENTRYPOINT ["/go/bin/bell"]
