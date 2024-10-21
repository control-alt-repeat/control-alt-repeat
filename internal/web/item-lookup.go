package web

import (
	"fmt"
	"net/http"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal"
	"github.com/labstack/echo/v4"
)

type Item struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Shelf string `json:"shelf"`
	Image string `json:"image"`
	Link  string `json:"link"`
}

func initialiseItemLookup(e *echo.Echo) {
	e.GET("item-lookup", renderItemLookupPage)
	e.POST("item-lookup", findItem)
	e.POST("item-print-shelf-label", printLabel)
}

func renderItemLookupPage(c echo.Context) error {
	fmt.Println("renderItemLookupPage")

	return render(http.StatusOK, "item-lookup.html", nil, c)
}

func findItem(c echo.Context) error {
	itemID := c.FormValue("itemID")

	result, err := internal.LookupItem(itemID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	item := &Item{
		ID:    result.ID,
		Title: result.Title,
		Shelf: result.Shelf,
		Image: result.ImageURL,
		Link:  result.EbayURL,
	}

	return c.JSON(http.StatusOK, item)
}

func printLabel(c echo.Context) error {
	itemID := c.FormValue("itemID")

	err := internal.ItemPrintShelfLabel(itemID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Shelf label printed successfully"})
}
