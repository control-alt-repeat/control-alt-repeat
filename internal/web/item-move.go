package web

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal"
)

type ItemMove struct {
	ItemID string `form:"id"`
	Shelf  string `form:"shelf"`
}

func initialiseItemMove(e *echo.Echo) {
	e.GET("item-move", renderItemMovePage)
	e.POST("item-move", itemMove)
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
	fmt.Println("Moving item with ID:", itemMove.ItemID, "to shelf:", itemMove.Shelf)

	err = internal.MoveItem(c.Request().Context(), itemMove.ItemID, itemMove.Shelf)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully moved the item"})
}
