# Nagase

> PoolC 홈페이지 API 서버

## Prerequisites

- NodeJS
- Docker, Docker Compose

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
# 로컬 DB 실행
docker-compose -f compose/local/docker-compose.yaml up -d

# 의존성 설치
yarn

# 테스트
yarn lint
yarn test:all

# 로컬 서버 실행
yarn start
```
