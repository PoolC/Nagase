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
docker run -e "..."
```
