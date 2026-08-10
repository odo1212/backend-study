package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// 아직 구현하지 않은 엔드포인트에서 임시로 반환하는 응답.
// TODO를 채워나가면서 이 함수 호출을 실제 로직으로 바꿔주세요.
func notImplemented(where string) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "TODO: "+where+" 를 구현하세요.")
}

// openapi.yaml 공통 에러 포맷: { "error": "메시지" }
// Echo 기본 에러 포맷({"message": "..."})을 여기서 덮어씁니다.
func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "서버 내부 오류가 발생했습니다."

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			message = msg
		}
	}

	if c.Response().Committed {
		return
	}
	if jsonErr := c.JSON(code, map[string]string{"error": message}); jsonErr != nil {
		c.Logger().Error(jsonErr)
	}
}
