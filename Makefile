build: 
#	GOOS=linux go build -o app 
	docker build -t restanrm/bell . 
#	rm app

run: build
	docker run --rm -it \
  -e POLLY_ACCESS_KEY=${POLLY_ACCESS_KEY} \
  -e POLLY_SECRET_KEY=${POLLY_SECRET_KEY} \
  -e FLITE=${FLITE} \
  -e POLLY_VOICE=${POLLY_VOICE} \
  -p 10101:10101 \
  -v /dev/snd:/dev/snd \
  -v /dev/shm:/dev/shm \
  -v /run/user/$uid/pulse:/run/user/$uid/pulse \
  -v /var/lib/dbus:/var/lib/dbus \
  --privileged \
  restanrm/bell:latest
	
.PHONY: build run
