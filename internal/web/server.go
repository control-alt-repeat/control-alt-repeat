package web

import (
	"embed"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type CustomRenderer struct {
	templates *template.Template
}

func (r *CustomRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

//go:embed templates/*
var templates embed.FS

//go:embed public/*
var public embed.FS

func Init(e *echo.Echo) error {
	t, err := template.ParseFS(templates, "templates/*")
	if err != nil {
		return err
	}

	e.Renderer = &CustomRenderer{templates: t}

	// Middleware
	e.Use(middleware.Recover())

	// e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	return username == "user" && password == "pass", nil
	// }))

	// Routes
	e.GET("/", showIndex)

	initialiseItemLookup(e)
	initialiseItemMove(e)
	initialiseEbayImportListing(e)

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "public",
		Filesystem: http.FS(public),
	}))

	return nil
}

func showIndex(c echo.Context) error {
	return render(http.StatusOK, "index.html", nil, c)
}

func render(code int, templateName string, data map[string]interface{}, c echo.Context) error {
	c.Logger().Print("Echo interface")
	zerolog.Ctx(c.Request().Context()).Printf("Zerolog interface")
	zerolog.Ctx(c.Request().Context()).Info().Fields(data).
		Msgf("Rendering '%s' with code %d", templateName, code)

	var builder strings.Builder

	if err := c.Echo().Renderer.Render(&builder, templateName, data, c); err != nil {
		return err
	}

	return c.Render(code, "base.html", map[string]interface{}{
		"content": builder.String(),
	})
}
