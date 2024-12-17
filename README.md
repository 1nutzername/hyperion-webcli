# hyperion
- configure LED device to autostart=off

# rpi
- copy hyperion-webcli folder to rpi
- run `docker build . -t hyperion-webcli:latest` to create an arm64 image on target system
- delete the copied folder since the image exists in registry now
- now it can be used with docker

# commands
- run go local without env variables: `go run .\main.go`
- build docker image: `docker build . -t hyperion-webcli:latest`
- run docker with params: `docker run -e PORT=xxxP -e HYPERION_JSON_RPC_URL="http://xxx.xxx.xxx.xxx:8090" -p xxxP:xxxP hyperion-webcli:latest`
- pack container to tar: `docker save -o hyperion-webcli.tar hyperion-webcli`
- load tar packed container to docker registry: `docker load -i hyperion-webcli.tar`

