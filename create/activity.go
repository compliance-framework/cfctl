package create

import (
	"fmt"
	"os"

	"cfctl/common"

	"github.com/compliance-framework/configuration-service/domain"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type CreateActivity struct {
	FilePath string `yaml:"file-path"`
	TaskID   string `yaml:"task-id"`
	PlanID   string `yaml:"plan-id"`
}

var CreateActivityVar CreateActivity

func RunCreateActivity(cmd *cobra.Command, args []string) {
	var yamlData []byte
	var err error
	var response string

	fileName := CreateActivityVar.FilePath
	activityPlanID := CreateActivityVar.PlanID
	activityTaskID := CreateActivityVar.TaskID

	// Process the command
	activity := domain.Activity{}
	yamlData, err = os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("error reading file: %v", err)
		os.Exit(1)
	}
	// Check the yaml is valid vs domain.Activity
	_, err = yaml.Marshal(&activity)
	if err != nil {
		fmt.Printf("Error marshalling to YAML: %v\n", err)
		os.Exit(1)
	}
	response, err = common.PostYAMLDocument(string(yamlData), common.CurrentContext.URL+"plan/"+activityPlanID+"/tasks/"+activityTaskID+"/activities")
	if err != nil {
		fmt.Printf("Error posting: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(response)
}
