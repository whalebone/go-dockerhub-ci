# go-dockerhub-ci

Service for hooking to DockerHub build webhook - sends dockerhub "build finished" message to slack channel

ENV variables:
- WEBHOOK - URL for sending slack message
- PATH_PREFIX (default /) - what URL path the service is listening at
- PORT (default 8080)
