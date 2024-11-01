package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/control-alt-repeat/control-alt-repeat/internal/freeagent"
)

// freeagent command: "car freeagent"
var cmdFreeagent = &cobra.Command{
	Use:   "freeagent",
	Short: "Freeagent related operations",
}

// Import freeagent subcommand: "car freeagent get-contact"
var cmdFreeagentGetContact = &cobra.Command{
	Use:   "get-contact",
	Short: "Gets a contact from freeagent",
	Run:   freeagentGetContact,
}

// Import freeagent subcommand: "car freeagent list-contacts"
var cmdFreeagentListContacts = &cobra.Command{
	Use:   "list-contacts",
	Short: "Lists contacts from freeagent",
	Run:   freeagentListContacts,
}

func freeagentGetContact(cmd *cobra.Command, args []string) {
	if _, err := strconv.Atoi(contactID); err != nil {
		fmt.Println("Error: contact ID must be a valid integer")
		os.Exit(1)
	}
	fmt.Println("Fetching contact from Freeagent with ID:", contactID)

	contact, err := freeagent.GetContact(cmd.Context(), freeagent.GetContactOptions{ContactID: contactID})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(contact.FirstName)
}

func freeagentListContacts(cmd *cobra.Command, args []string) {
	fmt.Println("Fetching contacts from Freeagent", contactID)

	contacts, err := freeagent.ListContacts(cmd.Context(), freeagent.ListContactsOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, contact := range contacts {
		fmt.Printf("%s %s\n", contact.ID(), contact.DisplayName())
	}
}
