package web

import (
	"fmt"
	"net/http"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal"
	"github.com/labstack/echo/v4"
)

func initialiseItemLookup(e *echo.Echo) {
	e.GET("item-lookup", renderItemLookupPage)
	e.POST("item-lookup", findItem)
	e.POST("item-print-shelf-label", printLabel)
}

func renderItemLookupPage(c echo.Context) error {
	return render(http.StatusOK, "item-lookup.html", nil, c)
}

func findItem(c echo.Context) error {
	itemID := c.FormValue("id")

	fmt.Println(itemID)

	result, err := internal.LookupItem(itemID)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	item := &Item{
		ID:             result.ID,
		Shelf:          result.Shelf,
		EbayReferences: []EbayReference{},
	}

	for _, ebayReference := range result.EbayReferences {
		item.EbayReferences = append(item.EbayReferences, EbayReference{
			Title:      ebayReference.Title,
			ImageURL:   ebayReference.ImageURL,
			ListingURL: ebayReference.ListingURL,
		})
	}

	return c.JSON(http.StatusOK, item)
}

func printLabel(c echo.Context) error {
	itemID := c.FormValue("id")

	err := internal.ItemPrintShelfLabel(itemID)

	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Shelf label printed successfully"})
}
