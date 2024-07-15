package create

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"github.com/compliance-framework/configuration-service/domain"
	"cfcli/common"
)

type CreateTask struct {
	Title       string             `yaml:"title"`
	Description string             `yaml:"description"`
	Type        string             `yaml:"type"`
	Schedule    string             `yaml:"schedule"`
	PlanID      string             `yaml:"plan-id"`
}

var CreateTaskVar     CreateTask

func RunCreateTask(cmd *cobra.Command, args []string) {
	var yamlData []byte
	var err error

	fileName        := common.RunConfig.FilePath
	taskTitle       := CreateTaskVar.Title
	taskDescription := CreateTaskVar.Description
	taskSchedule    := CreateTaskVar.Schedule
	taskPlanID      := CreateTaskVar.PlanID
	taskType        := CreateTaskVar.Type

	// Check variables are correct
	// one of title or filepath must be non-empty, and one of them must be empty
	if taskTitle == "" && fileName == "" {
		log.Fatalf("title, description, schedule and planid, or yaml filename must be set")
	} else if (taskTitle != "" || taskDescription != "" || taskSchedule != "" || taskPlanID != "") && fileName != "" {
		log.Fatalf("either: title, description, schedule and planid should all be unset; or yaml filename should be unset")
	}
	// Process the command
	task := domain.Task{}
	if fileName != "" {
		yamlData, err = ioutil.ReadFile(common.RunConfig.FilePath)
		if err != nil {
			log.Fatalf("error reading file: %v", err)
		}
		// Check the yaml is valid vs domain.Task
		_, err = yaml.Marshal(&task)
    	if err != nil {
    	    log.Fatalf("Error marshalling to YAML: %v\n", err)
    	}
	} else if taskTitle != "" {
		task.Title       = taskTitle
		task.Description = taskDescription
		task.Schedule    = taskSchedule
		task.Type        = domain.TaskType(taskType)
		yamlData, err = yaml.Marshal(&task)
    	if err != nil {
    	    log.Fatalf("Error marshalling to YAML: %v\n", err)
    	}
	} else {
		log.Fatalf("RunCreateTask should not get here")
	}
	err = common.PostYAMLDocument(string(yamlData), common.CurrentContext.URL + "plan/" + taskPlanID + "/tasks")
    if err != nil {
        fmt.Printf("Error posting: %v\n", err)
        return
	}
}

