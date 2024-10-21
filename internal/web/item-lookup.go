package web

import (
	"fmt"
	"net/http"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal"
	"github.com/labstack/echo/v4"
)

type LookupAssignment struct {
	Term string `form:"lookup_term"`
}

func showItemLookup(c echo.Context) error {
	fmt.Println("showItemLookup")

	return render(http.StatusOK, "item-lookup.html", nil, c)
}

func showItemLookupSubmit(c echo.Context) error {
	fmt.Println("showItemLookupSubmit")

	lookup := new(LookupAssignment)
	if err := c.Bind(lookup); err != nil {
		return c.String(http.StatusBadRequest, "Invalid form submission")
	}

	result, err := internal.LookupItem(lookup.Term)
	if err != nil {
		return showItemLookupError(c, err)
	}

	lookupMap := map[string]interface{}{
		"id":       result.ID,
		"shelf":    result.Shelf,
		"title":    result.Title,
		"imageURL": result.ImageURL,
		"eBayURL":  result.EbayURL,
	}

	fmt.Println("Details loaded")
	fmt.Println(lookupMap)

	return render(http.StatusOK, "item-lookup.html", lookupMap, c)
}

func showItemLookupError(c echo.Context, err error) error {
	fmt.Println("showItemLookupError")

	fmt.Println(err.Error())

	return render(http.StatusOK, "item-lookup.html", map[string]interface{}{
		"error": err.Error(),
	}, c)
}
