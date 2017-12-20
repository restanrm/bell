# bell 
bell is a small api to play sound via an API. 

Endpoint|Method|Description
--|--|--
/api/v1/ | GET | list registered sound that can be played
/api/v1/play/{sound} | GET | play a sound
/api/v1/tts | GET | retrieve an html gui to play text
/api/v1/tts | POST | send text to play

# dependencies 
This program needs `mpv` and `flite` to produce sound and to make text to speach
