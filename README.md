# bell 
bell is a small api to play sound via an API. 

| Endpoint               | Method | Description                               |
| ---------------------- | ------ | ----------------------------------------- |
| /api/v1/               | GET    | list registered sound that can be played  |
| /api/v1/play/{sound}   | GET    | play a sound                              |
| /api/v1/tts            | GET    | retrieve an html gui to play text         |
| /api/v1/tts            | POST   | send text to play                         |
| /api/v1/sounds         | GET    | list registered sounds that can be played |
| /api/v1/sounds         | POST   | add new sound to bell                     |
| /api/v1/sounds/{sound} | DELETE | remove sound from bell                    |

# dependencies 
This program needs `mpv` to play sound.

The text to speach functionnality need an aws pairs of key to work. It uses Polly service. 
[see here](https://console.aws.amazon.com/iam/home#/security_credential) to create services to access it.


# docker run 
```bash
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
  --group-add $(getent group audio | cut -d: -f3) \
  restanrm/bell:latest
```
