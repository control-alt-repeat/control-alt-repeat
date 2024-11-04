package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/control-alt-repeat/control-alt-repeat/internal"
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

	warehouseItem, ebayInternalItems, err := internal.LookupItem(c.Request().Context(), itemID)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	item := Map(warehouseItem)

	item.ImageURL = ebayInternalItems[0].PictureURL

	for _, ebayReference := range ebayInternalItems {
		item.EbayReferences = append(item.EbayReferences, EbayReference{
			Title:      ebayReference.Title,
			ImageURL:   ebayReference.PictureURL,
			ListingURL: ebayReference.ViewItemURL,
		})
	}

	return c.JSON(http.StatusOK, item)
}

func printLabel(c echo.Context) error {
	itemID := c.FormValue("id")

	err := internal.ItemPrintShelfLabel(c.Request().Context(), itemID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Shelf label printed successfully"})
}
