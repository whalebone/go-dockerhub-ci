# go-dockerhub-ci

Service for hooking to DockerHub build webhook - sends dockerhub "build finished" message to slack channel

ENV variables:
- WEBHOOK - URL for sending slack message
- PATH_PREFIX (default /) - what URL path the service is listening at
- PORT (default 8080)
- DEBUG - default 0, set to 1 to enable debugging incoming data to console
