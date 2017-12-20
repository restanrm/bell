FROM alpine

RUN apk update && apk add flite mpv
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY app app

WORKDIR /data
VOLUME /data

EXPOSE 8080
ENTRYPOINT ["/app"]
