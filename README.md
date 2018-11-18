# bell
`bell` is a sound box with an API. You can send sounds and play them with some API calls. It has also a Text to Speech feature. By default the tts functionnality is supported by flite, but with some credentials defined you can use the AWS api.

Initially, sounds could only be played on the server. Now you can register specific clients and order them to play sound in place of the server.

This tool is mainly created to add notifications to openspace or noisy notifications to everyone.

To get started:
```bash
go get github.com/restanrm/bell/...
bell
```
A docker is also available on `restanrm/bell`.

# bellctl
bellctl is the CLI that allows to interact with the bell server. You can upload, play sound, make backup, etc.
here is the help command:

```bash
You can control a bell server. To choose your bell server use the env variable BELL_ADDRESS.addCmd

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
  register    Register allows to connect to websocket of bell server. It will receive play orders and run them with `mpv`.
  restore     restore command help to put an archive sounds list back into a bell server
  say         say target use tts to say what you wrote

Flags:
  -h, --help      help for bellctl
  -v, --verbose   Increase verbosity

Use "bellctl [command] --help" for more information about a command.
```

## API
bell is a small api to play sound via an API.

| Endpoint               | Method | Description                               |
| ---------------------- | ------ | ----------------------------------------- |
| /api/v1/               | GET    | list registered sound that can be played  |
| /api/v1/play/{sound}   | GET    | play a sound                              |
| /api/v1/tts            | GET    | retrieve an html gui to play text         |
| /api/v1/tts            | POST   | send text to play                         |
| /api/v1/tts/retrieve   | POST   | retrieve mp3 of said text                 |
| /api/v1/sounds         | GET    | list registered sounds that can be played |
| /api/v1/sounds         | POST   | add new sound to bell                     |
| /api/v1/sounds/{sound} | DELETE | remove sound from bell                    |
| /api/v1/mattermost     | POST   | allow slash commands on mattermost        |


## Play on client
The API offer possibility to list the connected clients that can play music.

| Endpoint                 | Method | Description                               |
| ----------------------   | ------ | ----------------------------------------- |
| /api/v1/clients          | GET    | list clients that can play music          |
| /api/v1/clients/register | GET    | register to the websocket endpoint        |

A client can register itself in order to receive order to play some music.
This canal allow to communicate the name of the sound. The data still need to
be retrieved by the `registered` endpoint and played locally.

### Register a client
Message to register a new client. The name could be omitted, it will be replaced with an uuidV4 value. `bellctl` use the hostname by default.
```json
{
  "name":"name_of_the_client"
}
```
The received response contains the name of the client that will be used by the server.
```json
{
  "name":"string"
}
```

A client that want to play some sound an a registered endpoint make a query on the play endpoint `/api/v1/play/{sound}?destination={dest-name}` with the query parameter `destination` positionned to the name of the registered endpoint.

### Types of messages
Different kind of messages can be received:
- tts
- sound
- errors (not implemented yet).

The json format of a play order is the following:
```json
{
  "type":"error|tts|sound",
  "data": "payload. can be an error message, something to say, or a sound to retrieve."
}
```

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
  -e MATTERMOST_SLASH_TOKEN=${MATTERMOST_SLASH_TOKEN}
  -p 10101:10101 \
  --device /dev/snd \
  -e PULSE_SERVER=unix:${XDG_RUNTIME_DIR}/pulse/native \
  -v ${XDG_RUNTIME_DIR}/pulse/native:${XDG_RUNTIME_DIR}/pulse/native \
  -v ~/.config/pulse/cookie:/root/.config/pulse/cookie \
  --group-add $(getent group audio | cut -d: -f3) \
  restanrm/bell:latest
```


TODO:
- [x] implement tags in front
- [x] add dest selection
- [x] add registering selection
- [x] add tag selection
- [x] implement the registering part on the front
- [ ] improve upload in front
  - [ ] add tags listing
  - [ ] allow to create new tags
  - [ ] upload sound with selected tags
- [x] put everything in one page on the front
- [ ] create a good design to support it
- [ ] add TTS client
  - [ ] add new type of opbject in websockets
  - [ ] implement it on bellctl
  - [ ] implement it on frontend
  - [ ] update documentation


