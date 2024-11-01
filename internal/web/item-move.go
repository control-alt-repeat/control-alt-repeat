package web

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"

	"github.com/control-alt-repeat/control-alt-repeat/internal"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

type ItemMove struct {
	ItemID string `form:"id"`
	Shelf  string `form:"shelf"`
}

func initialiseItemMove(e *echo.Echo) {
	e.GET("item-move", renderItemMovePage)
	e.POST("item-move", itemMove)
	e.POST("items-unshelved", itemsUnshelved)
}

func renderItemMovePage(c echo.Context) error {
	return render(http.StatusOK, "item-move.html", nil, c)
}

func itemMove(c echo.Context) error {
	itemMove := &ItemMove{
		ItemID: c.FormValue("id"),
		Shelf:  c.FormValue("shelf"),
	}

	matched, err := regexp.MatchString(`^[A-Z]{3}-[0-9]{3}$`, itemMove.ItemID)
	if err != nil || !matched {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = internal.MoveItem(c.Request().Context(), itemMove.ItemID, itemMove.Shelf)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully moved the item"})
}

func itemsUnshelved(c echo.Context) error {
	items, err := warehouse.LoadUnshelvedItems(c.Request().Context())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	result := []Item{}

	for _, item := range items {
		result = append(result, Map(item))
	}

	return c.JSON(http.StatusOK, result)
}
