# API Specifications

## 변경 내역

  - 2018-09-23 : 최초 작성


## GrpahQL API

GraphQL 인터페이스를 이용한 API입니다. 파일 관련 동작을 제외한 모든 동작은 GraphQL API를 통해 이루어집니다.

### 요청 방법

요청은 GET, POST 두 가지 방식을 지원합니다. 아래는 각 방식으로 자신의 UUID를 조회하는 예제입니다.

```sh
# GET 방식
curl -H 'Content-Type: application/graphql' \
  -H 'Authorization: Bearer your_access_token' \
  'https://server_host/graphql?query=query{me{uuid}}'

# POST 방식
curl -X POST -H 'Content-Type: application/graphql' \
  -H 'Authorization: Bearer your_access_token' \
  -d 'query{me{uuid}}' \
  'https://server_host/graphql'
```

createMember, createAccessToken mutation을 제외한 모든 API 요청에는 인증이 필요합니다. 인증 방법은 [RFC 6750](https://tools.ietf.org/html/rfc6750) 표준을 따릅니다. [createAccessToken](#createaccesstoken) mutation으로 토큰을 발급받은 뒤, 발급받은 토큰을 예제와 같이 Authorization API 헤더에 넣어 요청해야합니다.

인증에 성공한 경우, 200 OK 상태 코드와 함께 `application/json; charset=utf-8` 포맷의 응답이 제공됩니다.

```js
// 오류가 발생하지 않은 경우
{
  "data": { "uuid": "00000000-0000-0000-0000-000000000000" }
}

// 오류가 발생한 경우
{
  "data": null,
  "errors": [{ "message": "some error occurred :(" }]
}
```

인증에 실패한 경우, 401 Unauthorized 응답이 반환됩니다.

### 자료형

이하 문서는 [GraphQL 자료형](https://graphql.org/learn/schema/)의 표현식을 따릅니다.

  - ID : [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier) 알고리즘으로 생성된 문자열
  - String : 문자열
  - Int : 64비트 정수
  - Boolean : 참/거짓
  - DateTime : [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) 형태로 포매팅 된 날짜/시각 정보 문자열

### 스펙

[GraphQL Playground](http://nagase.lynlab.co.kr/graphql)에서 스펙을 확인할 수 있습니다.


## FIle API

파일을 업로드, 다운로드 하기 위한 API 입니다.

### 요청 방법

요청은 일반적인 RESTful API의 형태로 이루어집니다.

  - `GET  /files/{fileName}`
  - `POST /files/{fileName}`

POST 요청에는 인증이 필요합니다. 인증 방식은 GraphQL API와 동일합니다.

#### GET /files/{fileName}

#### POST /files/{fileName}

Formdata의 multipart 업로드를 지원합니다. 업로드 할 파일의 form name은 `upload`로 지정해야합니다.
