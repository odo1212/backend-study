package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler

	// CORS 정책: 프론트 origin 허용
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:8001", "http://127.0.0.1:8001"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	// 2주차: 인증 API
	e.POST("/auth/register", registerUser)
	e.POST("/auth/login", loginUser)

	e.Logger.Fatal(e.Start(":8083"))
}
