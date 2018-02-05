FROM ubuntu:17.10

RUN apt-get update && apt-get -y install \
    alsa-base alsa-utils pulseaudio \
    golang git \
    flite mpv
RUN export GOPATH=/go && \
    go get github.com/restanrm/bell

WORKDIR /data
VOLUME /data

COPY store.json /data/store.json
COPY sounds /data/sounds

EXPOSE 8080

ENTRYPOINT ["/go/bin/bell"]
