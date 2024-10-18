package web

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomRenderer struct {
	templates *template.Template
}

func (r *CustomRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

//go:embed templates/*
var templates embed.FS

func Init() (*echo.Echo, error) {
	e := echo.New()

	t, err := template.ParseFS(templates, "templates/*")
	if err != nil {
		return e, err
	}

	e.Renderer = &CustomRenderer{templates: t}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	return username == "user" && password == "pass", nil
	// }))

	// Routes
	e.GET("/", showIndex)
	e.GET("/item-move", showItemMoveForm)
	e.POST("/item-move-submit", submitItemMoveForm)

	return e, nil
}

func showIndex(c echo.Context) error {
	fmt.Println("showIndex")

	return render(http.StatusOK, "index.html", nil, c)
}

func render(code int, templateName string, data map[string]interface{}, c echo.Context) error {
	var builder strings.Builder

	fmt.Println("Rendering: ", templateName)

	c.Echo().Renderer.Render(&builder, templateName, data, c)

	return c.Render(code, "base.html", map[string]interface{}{
		"content": builder.String(),
	})
}
