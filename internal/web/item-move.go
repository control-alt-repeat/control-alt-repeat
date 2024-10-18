package web

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal"
	"github.com/labstack/echo/v4"
)

type ItemAssignment struct {
	ItemID string `form:"item_id"`
	Shelf  string `form:"shelf"`
}

func showItemMoveForm(c echo.Context) error {
	return render(http.StatusOK, "item-move-form.html", nil, c)
}

func submitItemMoveForm(c echo.Context) error {
	item := new(ItemAssignment)
	if err := c.Bind(item); err != nil {
		return c.String(http.StatusBadRequest, "Invalid form submission")
	}

	matched, err := regexp.MatchString(`^[A-Z]{3}-[0-9]{3}$`, item.ItemID)
	if err != nil || !matched {
		return showItemMoveError(c, errors.New("item ID must be in the format A-Z-0-9 (e.g., A-123)"))
	}
	fmt.Println("Moving item with ID:", item.ItemID, "to shelf:", item.Shelf)

	err = internal.MoveItem(item.ItemID, item.Shelf)

	if err != nil {
		return showItemMoveError(c, err)
	}

	return render(http.StatusOK, "item-move-ok.html", nil, c)
}

func showItemMoveError(c echo.Context, error error) error {
	fmt.Println(error.Error())

	return render(http.StatusOK, "item-move-error.html", map[string]interface{}{
		"error": error.Error(),
	}, c)
}
