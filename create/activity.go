package create

//activity_id="$(curl -s localhost:8080/api/plan/"${plan_id}"/tasks/"${task_id}"/activities --header 'Content-Type: application/yaml' -d '
//title: Check server is OK
//description: This activity checks the server is OK
//provider:
//  name: ssh-cf-plugin
//  image: ghcr.io/compliance-framework/ssh-cf-plugin
//  tag: seed
//  configuration:
//    yaml: |
//      username: "'${CF_SSH_USERNAME}'"
//      password: "'${CF_SSH_PASSWORD}'"
//      host: "'${CF_SSH_HOST}'"
//      command: "'"${CF_SSH_COMMAND}"'"
//      port: "'${CF_SSH_PORT:-2227}'"
//subjects:
//  title: Server
//  description: "Server: '${CF_SSH_HOST}'"
//  labels: {}
//' | jq -r '.id')"
//
//Must be a yaml file with configuration, and -p and -t flags for the ids

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"github.com/compliance-framework/configuration-service/domain"
	"cfcli/common"
)

type CreateActivity struct {
	FilePath    string             `yaml:"file-path"`
	TaskID      string             `yaml:"task-id"`
	PlanID      string             `yaml:"plan-id"`
}

var CreateActivityVar     CreateActivity

func RunCreateActivity(cmd *cobra.Command, args []string) {
	var yamlData []byte
	var err      error
	var response string

	fileName       := CreateActivityVar.FilePath
	activityPlanID := CreateActivityVar.PlanID
	activityTaskID := CreateActivityVar.TaskID

	// Process the command
	activity := domain.Activity{}
	yamlData, err = ioutil.ReadFile(fileName)
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
	response, err = common.PostYAMLDocument(string(yamlData), common.CurrentContext.URL + "plan/" + activityPlanID + "/tasks/" + activityTaskID + "/activities")
	if err != nil {
		fmt.Printf("Error posting: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(response)
}

