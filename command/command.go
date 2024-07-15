package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"cfcli/create"
	"cfcli/common"
)

func ParseCommand() {
	var rootCmd = &cobra.Command{
		Use:   "cfcli",
		Short: "CLI for Compliance Framework (CF)",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			common.ReadConfigFile()
			common.ApplyContext()
		},
	}

	// Now the config file has been read in, apply any global overrides on the command line
	rootCmd.PersistentFlags().StringVarP(&common.RunConfig.Context, "context", "c", "", "Context to use from config file")
	rootCmd.PersistentFlags().StringVarP(&common.RunConfig.URL, "url", "u", "", "Endpoint URL to send the YAML payloads")
	rootCmd.PersistentFlags().StringVarP(&common.RunConfig.Output, "output", "o", "yaml", "Output format (TODO), defaults to json")

	// Define the validate subcommand
	var validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate YAML documents",
		Run:   common.RunValidate,
	}
	validateCmd.Flags().StringVarP(&common.RunConfig.FilePath, "file", "f", "", "YAML file to process")
	validateCmd.MarkFlagRequired("file")

	// Define the create subcommand
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create CF objects ",
	}

	// Define the create plan subcommand
	var createPlanCmd = &cobra.Command{
		Use:   "plan",
		Short: "Create a CF plan",
		Run:   create.RunCreatePlan,
	}
	createPlanCmd.Flags().StringVarP(&common.RunConfig.FilePath, "file", "f", "", "YAML file to process")
	createPlanCmd.Flags().StringVarP(&create.CreatePlanVar.Title, "title", "t", "", "Title of Plan")

	// Define the create task subcommand
	var createTaskCmd = &cobra.Command{
		Use:   "task",
		Short: "Create a CF task",
		Run:   create.RunCreateTask,
	}
	createTaskCmd.Flags().StringVarP(&common.RunConfig.FilePath,         "file",        "f", "",       "YAML file to process")
	createTaskCmd.Flags().StringVarP(&create.CreateTaskVar.Title,        "title",       "t", "",       "Title of Task")
	createTaskCmd.Flags().StringVarP(&create.CreateTaskVar.Description,  "description", "d", "",       "Description of Task")
	createTaskCmd.Flags().StringVarP(&create.CreateTaskVar.Schedule,     "schedule",    "s", "",       "Schedule of Task (cron format)")
	createTaskCmd.Flags().StringVarP(&create.CreateTaskVar.PlanID,       "plan-id",     "p", "",       "Plan ID")
	createTaskCmd.Flags().StringVarP(&create.CreateTaskVar.Type,         "type",        "y", "action", "Task type (default: action)")

	// Add subcommands to the create command
	createCmd.AddCommand(createPlanCmd, createTaskCmd)

	// Add subcommands to the root command
	rootCmd.AddCommand(validateCmd, createCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

