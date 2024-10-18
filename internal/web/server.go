package web

import (
	"io"
	"net/http"
	"path/filepath"
	"runtime"
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

func Init() *echo.Echo {
	e := echo.New()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		e.Logger.Fatal("Failed to get caller information")
	}
	basePath := filepath.Dir(filename)

	templatesPath := filepath.Join(basePath, "templates/*.html")

	renderer := &CustomRenderer{
		templates: template.Must(template.ParseGlob(templatesPath)),
	}
	e.Renderer = renderer

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return username == "user" && password == "pass", nil
	}))

	// Routes
	e.GET("/", showIndex)
	e.GET("/item-move", showItemMoveForm)
	e.POST("/item-move-submit", submitItemMoveForm)

	return e
}

func showIndex(c echo.Context) error {
	return render(http.StatusOK, "index.html", nil, c)
}

func render(code int, templateName string, data map[string]interface{}, c echo.Context) error {
	var builder strings.Builder

	c.Echo().Renderer.Render(&builder, templateName, data, c)

	return c.Render(code, "base.html", map[string]interface{}{
		"content": builder.String(),
	})
}
