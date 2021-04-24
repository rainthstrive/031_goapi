package employees

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routes(e *echo.Echo) {
	// Restricted group
	r := e.Group("/api")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("/employees", connectionSql)
}