package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/control-alt-repeat/control-alt-repeat/internal"
	"github.com/control-alt-repeat/control-alt-repeat/internal/freeagent"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

func initialiseItemLookup(e *echo.Echo) {
	e.GET("item-lookup", renderItemLookupPage)
	e.POST("item-lookup", findItem)
	e.POST("item-print-shelf-label", printLabel)
	e.POST("contacts-list", listContacts)
	e.POST("owner-save", saveOwner)
}

func renderItemLookupPage(c echo.Context) error {
	return render(http.StatusOK, "item-lookup.html", nil, c)
}

func findItem(c echo.Context) error {
	itemID := c.FormValue("id")

	fmt.Println(itemID)

	warehouseItem, ebayInternalItems, err := internal.LookupItem(c.Request().Context(), itemID)
	if err != nil {
		return handleError(c, err)
	}

	item := MapToWebItem(warehouseItem)

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

type ListContactsResponse struct {
	Contacts []Contact `json:"contacts"`
}

func listContacts(c echo.Context) error {
	contacts, err := freeagent.ListContacts(c.Request().Context(), freeagent.ListContactsOptions{})
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, &ListContactsResponse{Contacts: MapSlice(contacts, MapToWebContact)})
}

type PrintLabelResponse struct {
	Message string `json:"message"`
}

func printLabel(c echo.Context) error {
	itemID := c.FormValue("id")

	err := internal.ItemPrintShelfLabel(c.Request().Context(), internal.ItemPrintShelfLabelOptions{
		ItemID: itemID,
	})
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, &PrintLabelResponse{Message: "Shelf label printed successfully"})
}

func saveOwner(c echo.Context) error {
	itemID := c.FormValue("id")
	newOwnerID := c.FormValue("ownerId")
	newOwnerName := c.FormValue("newOwnerName")

	fmt.Printf("Updating item %s with new owner %s %s\n", itemID, newOwnerID, newOwnerName)

	err := warehouse.UpdateOwner(c.Request().Context(), itemID, newOwnerID, newOwnerName)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, &PrintLabelResponse{Message: "Owner saved successfully"})
}
