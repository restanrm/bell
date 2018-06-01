# bell
## API
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

## dependencies
This program needs `mpv` to play sound.

The text to speach functionnality need an aws pairs of key to work. It uses Polly service.
[see here](https://console.aws.amazon.com/iam/home#/security_credential) to create services to access it.


## docker run
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




# bellctl
bellctl is the CLI that allows to interact with the sounds. You can upload, play sound, make backup, etc.
here is the help command:

```bash
You can controal a bell server. To choose your bell server use the env variable BELL_ADDRESS.addCmd
Example:
	export BELL_ADDRESS=http://localhost:10101
	bellctl list

Usage:
  bellctl [command]

Available Commands:
  add         Add new sounds to library
  backup      backup the list of sounds into an archive
  delete      delete allows to remove sounds from library
  get         retrieve sound and store it locally
  help        Help about any command
  list        List available sounds to play
  play        Play sound on the host that run the server command
  restore     restore command help to put an archive sounds list back into a bell server
  say         say target use tts to say what you wrote

Flags:
  -h, --help   help for bellctl

Use "bellctl [command] --help" for more information about a command.
```