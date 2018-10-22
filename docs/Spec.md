# GraphQL API Specifications

## 변경 내역

  - 2018-09-23 : 최초 작성


## 요청과 응답

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


## 자료형

이하 문서는 [GraphQL 자료형](https://graphql.org/learn/schema/)의 표현식을 따릅니다.

  - ID : [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier) 알고리즘으로 생성된 문자열
  - String : 문자열
  - Int : 64비트 정수
  - Boolean : 참/거짓
  - DateTime : [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) 형태로 포매팅 된 날짜/시각 정보 문자열


## 모델

  - AccessToken
  - Board
  - Comment
  - Member
  - Post

### AccessToken

로그인 후 발급되는 Access Token입니다.

```graphql
type AccessToken {
  key:  String!
}
```

### Board

게시판입니다.

```graphql
type Board {
  name:            String!           # UI상으로 노출되는 이름
  urlPath:         String!           # URL 상으로 표현되는 경로
  readPermission:  BoardPermission!
  writePermission: BoardPermission!
}

enum BoardPermission {
  PUBLIC
  MEMBER
  ADMIN
}
```

### Comment

게시물(Post)의 댓글입니다.

```graphql
type Comment {
  id:        Int!
  author:    Member!
  body:      String!
  createdAt: DateTime!
}
```

### Member

회원 정보입니다.

```graphql
type Member {
  uuid:        ID!
  loginID:     String!
  email:       String!
  name:        String!
  department:  String!
  studentId:   String!
  isActivated: Boolean!
  isAdmin:     Boolean!
}
```

### Post

게시물입니다.

```graphql
type Post {
  id:        Int!
  author:    User!
  title:     String!
  body:      String!
  comments:  [Comment!]!
  createdAt: DateTime!
  updatedAt: DateTime!
}
```


## Query

  - [me](#me)
  - [boards](#boards)
  - [post](#post)

### me

자신의 회원 정보를 조회합니다.

응답 모델 : [Member!](#member)


## Mutation

  - [createAccessToken](#createaccesstoken)
  - [createComment](#createcomment)
  - [createMember](#createmember)
  - [createPost](#createpost)

### createAccessToken

Access Token을 발급합니다.

요청 :

```graphql
input Login {
  loginID:  String!
  password: String!
}
```

응답 모델 : [AccessToken!](#accesstoken)

### createComment

댓글을 작성합니다.

### createMember

회원을 추가합니다.

요청:

```graphql
input Member {
  loginID:     String!
  password:    String!
  email:       String!
  name:        String!
  department:  String!
  studentID:   String!
}
```

응답 모델 : [Member!](#member)

### createPost

게시물을 작성합니다.
