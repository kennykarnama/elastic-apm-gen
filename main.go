package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.elastic.co/apm/module/apmechov4/v2"
)

func main() {
	e := echo.New()
	e.Use(apmechov4.Middleware())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
