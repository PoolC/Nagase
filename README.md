# Nagase ![Docker Build Status](https://img.shields.io/docker/build/poolc/nagase.svg)

> PoolC 홈페이지 API 서버

## Prerequisites

  - Go 1.11
  - PostgreSQL
  - Docker


## 개발

### 환경변수 설정

환경변수 설정을 위해 [direnv](http://direnv.net)를 사용할 수 있습니다.

```sh
cp .envrc.example .envrc
vi .envrc

direnv allow
```

### 실행

```sh
# Test
go test -cover ./...

# Run server
go run main.go
```

## 배포

```sh
# Build a docker image
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o bin/nagase main.go
docker build -t poolc/nagase .

# Run a docker container
docker run 
```
