package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// https://yourbasic.org/golang/format-parse-string-time-date-example/
const (
	layoutUS = "January 2, 2006"
)

func DateFormatAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		_, err := time.Parse(layoutUS, authHeader)
		if err != nil {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}

func RequestLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		c.Set("customValueFromContext", time.Now().String())
		err := next(c)
		fmt.Printf("REQUEST: uri: %v, status: %v, datetime: %v\n", req.RequestURI, res.Status, c.Get("customValueFromContext"))
		return err
	}
}
