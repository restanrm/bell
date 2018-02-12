build: 
#	GOOS=linux go build -o app 
	docker build ${OPTS} -t restanrm/bell . 
#	rm app

run: build
	docker run --rm -it \
  -e POLLY_ACCESS_KEY=${POLLY_ACCESS_KEY} \
  -e POLLY_SECRET_KEY=${POLLY_SECRET_KEY} \
  -e FLITE=${FLITE} \
  -e POLLY_VOICE=${POLLY_VOICE} \
  -p 10101:10101 \
  --device /dev/snd \
  -e PULSE_SERVER=unix:${XDG_RUNTIME_DIR}/pulse/native \
  -v ${XDG_RUNTIME_DIR}/pulse/native:${XDG_RUNTIME_DIR}/pulse/native \
  -v ~/.config/pulse/cookie:/root/.config/pulse/cookie \
  --group-add $(shell getent group audio | cut -d: -f3) \
  restanrm/bell:latest
	
.PHONY: build run
