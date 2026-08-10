package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)


// TODO: User 구조체 정의
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  []byte    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TODO: RegisterRequest 구조체 정의
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


// TODO: UserResponse 구조체 정의
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
// TODO: AuthResponse 구조체 정의
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
// TODO: JWT Claims 구조체 정의
type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// TODO: 인메모리 사용자 저장소
var users = make(map[string]User)

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "dev-only-secret-change-me"))

// TODO: JWT 시크릿 키
// 하드코딩 대신 환경변수로 읽는 걸 권장합니다 (os.Getenv("JWT_SECRET"), 없으면 개발용 기본값).
// var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "dev-only-secret-change-me"))

// registerUser handles POST /auth/register
// 학습 포인트: 사용자 등록, 비밀번호 해싱, 중복 검사, 입력 검증
func registerUser(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "요청 본문을 읽을 수 없습니다.")
	}

	if !strings.Contains(req.Email, "@") {
		return echo.NewHTTPError(http.StatusBadRequest, "올바른 이메일 형식이 아닙니다.")
	}
	if len(req.Password) < 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "비밀번호는 8자 이상이어야 합니다.")
	}
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name은 필수입니다.")
	}

	if _, exists := users[req.Email]; exists {
		return echo.NewHTTPError(http.StatusConflict, "이미 등록된 이메일입니다.")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "비밀번호 처리 중 오류가 발생했습니다.")
	}

	now := time.Now().UTC()
	user := User{
		ID:        uuid.NewString(),
		Email:     req.Email,
		Name:      req.Name,
		Password:  hashed,
		CreatedAt: now,
		UpdatedAt: now,
	}
	users[user.Email] = user

	return c.JSON(http.StatusCreated, toUserResponse(user))
}


// loginUser handles POST /auth/login
// 학습 포인트: 인증 토큰 기반 인증, JWT 생성/검증
func loginUser(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "요청 본문을 읽을 수 없습니다.")
	}
	if req.Email == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "이메일과 비밀번호는 필수입니다.")
	}

	user, ok := users[req.Email]
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "이메일 또는 비밀번호가 올바르지 않습니다.")
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "이메일 또는 비밀번호가 올바르지 않습니다.")
	}

	token, err := generateJWT(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "토큰 생성 중 오류가 발생했습니다.")
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	})
}

func generateJWT(u User) (string, error) {
	claims := Claims{
		UserID: u.ID,
		Email:  u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func toUserResponse(u User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// TODO: toUserResponse 함수 구현
// - User를 UserResponse로 변환 (비밀번호 제외)
