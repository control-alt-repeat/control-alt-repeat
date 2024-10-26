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

func Init(e *echo.Echo) error {
	t, err := template.ParseFS(templates, "templates/*")
	if err != nil {
		return err
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

	initialiseItemLookup(e)
	initialiseItemMove(e)
	initialiseEbayImportListing(e)

	return nil
}

func showIndex(c echo.Context) error {
	fmt.Println("showIndex")

	return render(http.StatusOK, "index.html", nil, c)
}

func render(code int, templateName string, data map[string]interface{}, c echo.Context) error {
	var builder strings.Builder

	if err := c.Echo().Renderer.Render(&builder, templateName, data, c); err != nil {
		return err
	}

	return c.Render(code, "base.html", map[string]interface{}{
		"content": builder.String(),
	})
}
