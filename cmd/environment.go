package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/wakisa/qatarina-cli/internal/client"
	"github.com/wakisa/qatarina-cli/internal/schema"
)

var environmentCmd = &cobra.Command{
	Use:   "environments",
	Short: "Environments commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var createEnvironmentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new moudle",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetInt64("project-id")
		name, _ := cmd.Flags().GetString("name")
		url, _ := cmd.Flags().GetString("url")
		description, _ := cmd.Flags().GetString("description")

		payload := schema.CreateEnvironmentRequest{
			Name:        name,
			BaseURL:     url,
			Description: description,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		resp, err := client.Default().Post(fmt.Sprintf("v1/projects/%d/environments", projectID), body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		if resp.StatusCode != 201 || resp.StatusCode != 200 {
			return fmt.Errorf("API error: %s", string(bodyBytes))
		}

		fmt.Println("environment created succcessfully.")
		return nil
	},
}

var updateEnvironmentCmd = &cobra.Command{
	Use:   "update <environmentID>",
	Short: "Update a environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid environment ID: %w", err)
		}

		projectID, _ := cmd.Flags().GetInt64("project-id")
		name, _ := cmd.Flags().GetString("name")
		url, _ := cmd.Flags().GetString("url")
		description, _ := cmd.Flags().GetString("description")

		payload := schema.UpdateEnvironmentRequest{
			ID:          int64(id),
			Name:        name,
			BaseURL:     url,
			Description: description,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		resp, err := client.Default().Post(fmt.Sprintf("v1/projects/%d/environments/%s", projectID, args[0]), body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("API error: %s", string(bodyBytes))
		}

		fmt.Println("environment updated successfully.")
		return nil
	},
}

var listEnvironmentsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetInt64("project-id")
		resp, err := client.Default().Get(fmt.Sprintf("v1/projects/%d/environments", projectID))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var environmentResponse schema.EnvironmentListResponse
		if err := json.NewDecoder(resp.Body).Decode(&environmentResponse); err != nil {
			return fmt.Errorf("failed to decode reponse: %w", err)
		}

		for _, m := range environmentResponse.Environments {
			fmt.Printf("• [%d] %s URL: %s\n", m.ID, m.Name, m.BaseURL)
		}
		return nil
	},
}

var viewEnvironmentCmd = &cobra.Command{
	Use:   "view <environmentID>",
	Short: "View environment details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetInt64("project-id")
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid environment ID: %w", err)
		}
		resp, err := client.Default().Get(fmt.Sprintf("v1/projects/%d/environments/%s", projectID, args[0]))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode != 200 {
			if len(bodyBytes) == 0 {
				return fmt.Errorf("environment not found (ID: %d)", id)
			}
			return fmt.Errorf("environment not found (ID: %d): %s", id, string(bodyBytes))
		}

		var environment schema.EnvironmentResponse
		if err := json.Unmarshal(bodyBytes, &environment); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		fmt.Printf("environment: %s\nID: %d\nDescription: %s\n", environment.Name, environment.ID, environment.Description)
		return nil
	},
}

var deleteEnvironmentCmd = &cobra.Command{
	Use:   "delete <environmentID>",
	Short: "Delete an environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetInt64("project-id")
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid environment ID: %w", err)
		}
		resp, err := client.Default().Delete(fmt.Sprintf("v1/projects/%d/environments/%s", projectID, args[0]))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("failed to delete environment (ID: %d): %s", id, string(bodyBytes))
		}

		fmt.Println("environment deleted successfully.")
		return nil
	},
}

func init() {
	environmentCmd.PersistentFlags().Int64("project-id", 0, "Project ID")

	createEnvironmentCmd.Flags().String("name", "", "environment name")
	createEnvironmentCmd.Flags().String("url", "", "environment url")
	createEnvironmentCmd.Flags().String("description", "", "environment description")

	updateEnvironmentCmd.Flags().String("name", "", "environment name")
	updateEnvironmentCmd.Flags().String("url", "", "environment url")
	updateEnvironmentCmd.Flags().String("description", "", "environment description")

	environmentCmd.AddCommand(createEnvironmentCmd)
	environmentCmd.AddCommand(updateEnvironmentCmd)
	environmentCmd.AddCommand(listEnvironmentsCmd)
	environmentCmd.AddCommand(viewEnvironmentCmd)
	environmentCmd.AddCommand(deleteEnvironmentCmd)

	rootCmd.AddCommand(environmentCmd)
}
