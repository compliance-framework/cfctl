package command

import (
	"log"
	"os"

	"cfctl/activate"
	"cfctl/common"
	"cfctl/create"
	"cfctl/validate"

	"github.com/spf13/cobra"
)

func ParseCommand() {
	var rootCmd = &cobra.Command{
		Use:   "cfctl",
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

	// validate
	var validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate YAML documents",
		Run:   validate.RunValidate,
	}
	validateCmd.Flags().StringVarP(&validate.ValidateVar.FilePath, "file", "f", "", "YAML file to process")
	validateCmd.MarkFlagRequired("file")

	// create
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create CF objects ",
	}

	var createPlanCmd = &cobra.Command{
		Use:   "plan",
		Short: "Create a CF plan",
		Run:   create.RunCreatePlan,
	}
	createPlanCmd.Flags().StringVarP(&create.CreatePlanVar.FilePath, "file", "f", "", "YAML file to process")
	createPlanCmd.Flags().StringVarP(&create.CreatePlanVar.Title, "title", "t", "", "Title of Plan")

	var createTaskCmd = &cobra.Command{
		Use:   "task",
		Short: "Create a CF task",
		Run:   create.RunCreateTask,
	}
	createTaskCmd.Flags().StringVarP(&create.CreateTaskVar.FilePath, "file", "f", "", "YAML file to process")
	createTaskCmd.Flags().StringVarP(&create.CreateTaskVar.Title, "title", "t", "", "Title of Task")
	createTaskCmd.Flags().StringVarP(&create.CreateTaskVar.Description, "description", "d", "", "Description of Task")
	createTaskCmd.Flags().StringVarP(&create.CreateTaskVar.Schedule, "schedule", "s", "", "Schedule of Task (cron format)")
	createTaskCmd.Flags().StringVarP(&create.CreateTaskVar.PlanID, "plan-id", "p", "", "Plan ID")
	createTaskCmd.Flags().StringVarP(&create.CreateTaskVar.Type, "type", "y", "action", "Task type (default: action)")

	var createActivityCmd = &cobra.Command{
		Use:   "activity",
		Short: "Create a CF task",
		Run:   create.RunCreateActivity,
	}
	createActivityCmd.Flags().StringVarP(&create.CreateActivityVar.FilePath, "file", "f", "", "YAML file to process configuration of plugin")
	createActivityCmd.Flags().StringVarP(&create.CreateActivityVar.TaskID, "task-id", "t", "", "Title of Task")
	createActivityCmd.Flags().StringVarP(&create.CreateActivityVar.PlanID, "plan-id", "p", "", "Description of Task")
	createCmd.MarkFlagRequired("file")
	createCmd.MarkFlagRequired("task-id")
	createCmd.MarkFlagRequired("plan-id")

	// activate
	var activateCmd = &cobra.Command{
		Use:   "activate",
		Short: "Activate CF objects ",
	}

	var activatePlanCmd = &cobra.Command{
		Use:   "plan",
		Short: "Activate a CF plan",
		Args:  cobra.ExactArgs(1),
		Run:   activate.RunActivatePlan,
	}

	// Add subcommands to the create command
	createCmd.AddCommand(createPlanCmd, createTaskCmd, createActivityCmd)
	activateCmd.AddCommand(activatePlanCmd)

	// Add subcommands to the root command
	rootCmd.AddCommand(validateCmd, createCmd, activateCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
