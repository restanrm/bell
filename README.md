# bell 
bell is a small api to play sound via an API. 

Endpoint|Method|Description
--|--|--
/api/v1/ | GET | list registered sound that can be played
/api/v1/play/{sound} | GET | play a sound
/api/v1/tts | GET | retrieve an html gui to play text
/api/v1/tts | POST | send text to play

# dependencies 
This program needs `mpv` to play sound.

The text to speach functionnality need an aws pairs of key to work. It uses Polly service. 
[see here](https://console.aws.amazon.com/iam/home#/security_credential) to create services to access it.


# docker run 
```bash
docker run --rm -it \
  -e POLLY_ACCESS_KEY=$POLLY_ACCESS_KEY \
  -e POLLY_SECRET_KEY=$POLLY_SECRET_KEY \
  -e FLITE=0 \
  -p 10101:10101 \
  -v /dev/snd:/dev/snd \
  -v /dev/shm:/dev/shm \
  -v /run/user/$uid/pulse:/run/user/$uid/pulse \
  -v /var/lib/dbus:/var/lib/dbus \
  --privileged \
  restanrm/bell:latest
```
