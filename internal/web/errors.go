package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func handleError(c echo.Context, err error) error {
	fmt.Println(err)
	return c.JSON(http.StatusInternalServerError, echo.Map{
		"error": err.Error(),
	})
}
