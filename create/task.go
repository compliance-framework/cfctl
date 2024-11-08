package create

import (
	"fmt"
	"os"

	"cfctl/common"

	"github.com/compliance-framework/configuration-service/domain"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type CreateTask struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Schedule    string `yaml:"schedule"`
	PlanID      string `yaml:"plan-id"`
	FilePath    string `yaml:"file-path"`
}

var CreateTaskVar CreateTask

func RunCreateTask(cmd *cobra.Command, args []string) {
	var yamlData []byte
	var err error
	var response string

	fileName := common.RunConfig.FilePath
	taskTitle := CreateTaskVar.Title
	taskDescription := CreateTaskVar.Description
	taskSchedule := CreateTaskVar.Schedule
	taskPlanID := CreateTaskVar.PlanID
	taskType := CreateTaskVar.Type

	// Check variables are correct
	// one of title or filepath must be non-empty, and one of them must be empty
	if taskTitle == "" && fileName == "" {
		fmt.Printf("title, description, schedule and planid, or yaml filename must be set")
		os.Exit(1)
	} else if (taskTitle != "" || taskDescription != "" || taskSchedule != "" || taskPlanID != "") && fileName != "" {
		fmt.Printf("either: title, description, schedule and planid should all be unset; or yaml filename should be unset")
		os.Exit(1)
	}
	// Process the command
	task := domain.Task{}
	if fileName != "" {
		yamlData, err = os.ReadFile(common.RunConfig.FilePath)
		if err != nil {
			fmt.Printf("error reading file: %v", err)
			os.Exit(1)
		}
		// Check the yaml is valid vs domain.Task
		_, err = yaml.Marshal(&task)
		if err != nil {
			fmt.Printf("Error marshalling to YAML: %v\n", err)
			os.Exit(1)
		}
	} else if taskTitle != "" {
		task.Title = taskTitle
		task.Description = taskDescription
		task.Schedule = taskSchedule
		task.Type = domain.TaskType(taskType)
		yamlData, err = yaml.Marshal(&task)
		if err != nil {
			fmt.Printf("Error marshalling to YAML: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("RunCreateTask should not get here")
		os.Exit(1)
	}
	response, err = common.PostYAMLDocument(string(yamlData), common.CurrentContext.URL+"plan/"+taskPlanID+"/tasks")
	if err != nil {
		fmt.Printf("Error posting: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(response)
}
