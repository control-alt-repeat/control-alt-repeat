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

// Import listing subcommand: "car freeagent get-contact"
var cmdFreeagentGetContact = &cobra.Command{
	Use:   "get-contact",
	Short: "Gets a contact from freeagent",
	Run:   freeagentGetContact,
}

func freeagentGetContact(cmd *cobra.Command, args []string) {
	if _, err := strconv.Atoi(contactID); err != nil {
		fmt.Println("Error: contact ID must be a valid integer")
		os.Exit(1)
	}
	fmt.Println("Fetching contact from Freeagent with ID:", contactID)

	contact, err := freeagent.GetContact(cmd.Context(), contactID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(contact.FirstName)
}
