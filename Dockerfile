FROM golang:alpine

RUN apk update && apk add flite mpv ca-certificates git
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN go get github.com/restanrm/bell

WORKDIR /data
VOLUME /data

COPY store.json /data/store.json
COPY sounds /data/sounds

EXPOSE 8080
ENTRYPOINT ["bell"]
