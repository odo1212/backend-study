# backend-template-go-echo

`openapi.yaml` / `2주차api명세.md`의 인증 명세를 채워나가는 **Go + Echo 스타터
템플릿**입니다. CORS, 공통 에러 포맷, 프로젝트 골격은 이미 완성되어 있고,
`auth.go` 안의 **TODO만 채우면 됩니다.**

## 요구 사항

- Go 1.21+ (권장: 1.26+)

## 시작하는 법

1. 이 폴더를 통째로 복사해서 `backend-<본인이름>/`으로 이름을 바꾸세요.
   ```bash
   cp -r backend-template-go-echo backend-내이름
   cd backend-내이름
   ```
2. 실행해보세요. **바로 켜집니다** (아무것도 구현 안 해도 됨).
   ```bash
   go run .   # http://localhost:8083
   ```
3. 아무 엔드포인트나 호출해보세요. 아직 구현을 안 했기 때문에 이렇게 나옵니다.
   ```bash
   curl -X POST http://localhost:8083/auth/register
   # {"error":"TODO: implement registerUser"}
   ```
4. `auth.go`를 열어서 TODO 주석을 하나씩 실제 로직으로 바꾸세요. TODO 옆에
   구체적으로 어떤 라이브러리/함수를 쓰면 되는지 적어뒀습니다 (`bcrypt`,
   `github.com/golang-jwt/jwt/v5`, `github.com/google/uuid` — 셋 다 이미
   `go.mod`에 등록되어 있으니 `import`만 추가하면 됩니다).
   채울 때마다 다시 실행해서 curl로 확인하세요.

## 이미 되어 있는 것 / 내가 채워야 하는 것

| 이미 되어 있음 (건드릴 필요 없음) | 내가 채워야 함 |
|---|---|
| `main.go` — CORS, 라우팅 등록(`/auth/register`, `/auth/login`), 서버 부트스트랩 | `auth.go`의 `registerUser`/`loginUser` 함수 본문 |
| `errors.go` — 공통 에러 포맷(`{ "error": "..." }`) 변환 | `auth.go` 상단의 `User`/`RegisterRequest`/`LoginRequest`/`UserResponse`/`AuthResponse`/`Claims` 구조체 정의 |
| `go.mod` — `golang-jwt/jwt`, `golang.org/x/crypto`(bcrypt), `google/uuid` 의존성 | `generateJWT`, `toUserResponse` 함수 |

## 막힐 때 참고할 것

- `2주차api명세.md`, `openapi.yaml` — 각 엔드포인트의 정확한 요청/응답 형태 (원본)
- `mock-server/src/index.js` — 같은 명세를 Node/Express로 구현한 예시 (bcrypt 해싱,
  JWT 발급, 에러 처리 흐름을 그대로 참고할 수 있습니다)

## Swagger UI로 명세 보면서 확인하기

`openapi.yaml`을 브라우저에서 보기 좋게 렌더링해서, 각 엔드포인트를 펼쳐 값을 입력하고
바로 실제 요청까지 날려볼 수 있는 도구입니다. TODO를 하나씩 채울 때마다 이걸로 확인하면
curl보다 눈으로 보기 편합니다.

### 1. 실행하기

**꼭 `openapi.yaml`이 있는 프로젝트 루트에서** 실행하세요 (백엔드 폴더 안에서 실행하면
파일을 못 찾습니다).

```bash
cd backend-study   # openapi.yaml이 있는 최상위 폴더
npx -y swagger-ui-watcher openapi.yaml --port 8001 --host localhost
```

터미널에
`Listening on http://localhost:8001`이 뜨면 성공, 브라우저에서 그 주소로 접속하세요.
`openapi.yaml`이 바뀌면 자동으로 새로고침되니 계속 켜둬도 됩니다.

### 2. 내 백엔드 포트로 연결하기

화면 상단 "Servers" 드롭다운에서 `http://localhost:{port}` 항목을 선택하면, 바로 아래
`port` 입력창이 나타납니다. 여기에 본인 백엔드가 실행 중인 포트(템플릿 기본값 `8083`,
바꿨다면 그 값)를 입력하세요.

### 3. TODO 채우고 → Swagger UI에서 확인, 반복

1. `auth.go`에서 TODO 하나를 실제 로직으로 채움
2. 서버 재시작 (`Ctrl+C` → `go run .`)
3. Swagger UI에서 `/auth/register` 또는 `/auth/login`을 펼치고 "Try it out" → 값
   입력 → "Execute"
4. 상태 코드/응답 바디가 명세대로 나오는지 확인 (201/200뿐 아니라 400/401/409도
   일부러 틀린 값을 넣어서 확인해보세요)
5. 다음 TODO로 이동, 반복

### 4. CORS

`main.go`의 CORS 허용 목록에 Swagger UI 테스트용 origin(`http://localhost:8001`,
`http://127.0.0.1:8001`)이 프론트 origin과 함께 이미 등록되어 있어서 별도 설정 없이
바로 "Execute"가 됩니다. **제출 전에는 이 두 줄을 지우고 `http://localhost:3000`만
남겨주세요** — 명세가 요구하는 건 그것뿐이고, 테스트용 origin을 제출물에 남길 필요는
없습니다.

## 완료 체크리스트

- [ ] `POST /auth/register` — 이메일 중복 시 409, 유효성 검사 실패 시 400
- [ ] 비밀번호 해싱 (`bcrypt`, 응답에 절대 포함하지 않기)
- [ ] `POST /auth/login` — 이메일/비밀번호 불일치 시 401 (메시지는 동일하게)
- [ ] JWT 토큰 생성 (HS256, 24시간 만료) 및 `AuthResponse`로 반환
- [ ] `createdAt`/`updatedAt`이 `2026-08-05T10:00:00Z`처럼 `Z`가 붙은 UTC로 나가는지 확인
- [ ] `frontend/.env.local`의 `NEXT_PUBLIC_API_URL`을 내 포트로 바꿔서 실제 화면에서 확인
- [ ] `main.go`의 포트를 다른 참여자와 안 겹치는 포트로 변경했는지 확인
- [ ] 이 README를 내 프로젝트에 맞게 실행법/포트/환경변수로 고쳐쓰기
