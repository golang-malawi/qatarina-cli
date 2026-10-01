package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wakisa/qatarina-cli/internal/client"
)

var inboxCmd = &cobra.Command{
	Use:     "inbox",
	Aliases: []string{"i"},
	Short:   "list all test cases in authenticated user's inbox",
	Example: "qatarina-cli inbox list",
	RunE: func(cmd *cobra.Command, args []string) error {
		return inboxListCmd.RunE(cmd, args)
	},
}

var inboxListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list all test cases in authenticated user's inbox",
	Example: "qatarina-cli inbox list",
	RunE: func(cmd *cobra.Command, args []string) error {
		testCases, err := client.FetchInbox(context.Background(), &client.InboxFilterParams{})
		if err != nil {
			return err
		}

		if len(testCases) == 0 {
			fmt.Println("No test cases found.")
			return nil
		}

		fmt.Printf("Found %d Test Cases in your inbox", len(testCases))
		for _, tc := range testCases {
			fmt.Printf("• %s\n Code: %s\n Kind: %s\n ID: %s\n\n", tc.Title, tc.Code, tc.Kind, tc.ID)
		}
		return nil
	},
}

func init() {
	inboxCmd.AddCommand(inboxListCmd)
	rootCmd.AddCommand(inboxCmd)
}
