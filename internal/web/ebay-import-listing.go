package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/control-alt-repeat/control-alt-repeat/internal"
)

func initialiseEbayImportListing(e *echo.Echo) {
	e.GET("ebay-import-listing", renderEbayImportListingPage)
	e.POST("ebay-import-listing", ebayImportListing)
}

func renderEbayImportListingPage(c echo.Context) error {
	return render(http.StatusOK, "ebay-import-listing.html", nil, c)
}

func ebayImportListing(c echo.Context) error {
	listingID := c.FormValue("listingID")

	warehouseID, err := internal.ImportEbayListingByID(c.Request().Context(), listingID)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{"warehouseID": warehouseID})
}
