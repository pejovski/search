# Search

## Requirements

```bash
go version 1.13
docker
docker-compose
golang-statik: sudo apt install golang-statik
```

## Installation

Use [go mod](https://blog.golang.org/using-go-modules) to install dependencies.

```bash
go mod tidy
```

Run [docker-compose](https://docs.docker.com/compose/) to build docker images and run containers.

```bash
docker-compose up -d
```

if the command above returns this error

```bash
es01 exited with code 78
```

that is because

```bash
elasticsearch     | [1]: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]
```

Therefore we need to increase the vm.max_map_count limit:

```bash
sudo sysctl -w vm.max_map_count=524288
```

Now we need to edit /etc/sysctl.conf so the setting will also be in effect after a reboot.

Look for any vm.max_map_count line in /etc/sysctl.conf. If you find one, set its value to 524288. If there is no such line present, add the line

```bash
vm.max_map_count=524288
```
to the end of /etc/sysctl.conf

### Usage

```bash
docker-compose up -d
```
wait some time until elasticsearch is ready
- open [Search API](http://localhost:8207)
- use Api-Key: e9435160-1df0-4d27-b0bf-9b3edd9d5737 for authentication

## Swagger update
- use http://editor.swagger.io
- modify app/swagger/swagger.yaml
- run: statik -src=./app/swagger -dest=./app

# Architecture and Design

The project code follows the design principles from the resources bellow

### Microsoft Micro-Services

https://docs.microsoft.com/en-us/dotnet/architecture/microservices/index

### Uber Go Code Structure

https://www.youtube.com/watch?v=nLskCRJOdxM
- extended with receivers and emitters for working with events

### Rest HTTP Server by Go veteran

https://www.youtube.com/watch?v=rWBSMsLG8po

## License
[MIT](https://choosealicense.com/licenses/mit/)