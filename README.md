# Nagase [![Docker Build Status](https://img.shields.io/docker/build/poolc/nagase.svg)](https://hub.docker.com/r/poolc/nagase)

> PoolC 홈페이지 API 서버

## Prerequisites

  - Go 1.11
  - PostgreSQL
  - Docker


## 개발

### 환경변수 설정

서버 구동에 필요한 환경변수를 설정합니다. [direnv](http://direnv.net)를 사용하여 편리하게 환경변수를 설정할 수 있습니다.

```sh
cp .envrc.example .envrc
vi .envrc

direnv allow
```

그 다음, `secrets` 디렉토리에 아래 시크릿 파일들을 추가합니다.

  - `service-account.json` : Firebase 관련 기능을 사용하기 위한 비공개 키입니다. [Firebase Console](https://console.firebase.google.com)에서 발급받을 수 있습니다.

### 실행

```sh
# Install dependencies
go mod tidy

# Test
go test -cover ./...

# Run server
docker run -e "POSTGRES_USER=$DB_USERNAME" -e "POSTGRES_PASSWORD=$DB_PASSWORD" -e "POSTGRES_DB=$DB_NAME" -p 5432:5432 --name nagase-db -d postgres:9.6
go run main.go
```

## 배포

```sh
# Build a docker image
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o bin/nagase main.go
docker build -t poolc/nagase .

# Run a docker container
docker run -e "..." -p 8080:8080 poolc/nagase
```
