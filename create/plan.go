package create

import (
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"github.com/compliance-framework/configuration-service/domain"
	"cfcli/common"
)

type CreatePlan struct {
	Title    string `yaml:"title"`
	FilePath string `yaml:"file-path"`
}

var CreatePlanVar     CreatePlan

func RunCreatePlan(cmd *cobra.Command, args []string) {
	var yamlData []byte
	var err      error
	var response string

	fileName  := CreatePlanVar.FilePath
	planTitle := CreatePlanVar.Title

	// Check variables are correct
	// one of title or filepath must be non-empty, and one of them must be empty
	if planTitle == "" && fileName == "" {
		log.Fatalf("title or yaml filename must be set")
	} else if planTitle != "" && fileName != "" {
		log.Fatalf("title and yaml cannot both be set")
	}

	// Process the command
	plan := domain.Plan{}
	if fileName != "" {
		yamlData, err = ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatalf("error reading file: %v", err)
		}
		// Check the yaml is valid vs domain.Plan
		_, err = yaml.Marshal(&plan)
		if err != nil {
			log.Fatalf("Error marshalling to YAML: %v\n", err)
		}
	} else if planTitle != "" {
		plan.Title = planTitle
		yamlData, err = yaml.Marshal(&plan)
		if err != nil {
			log.Fatalf("Error marshalling to YAML: %v\n", err)
		}
	} else {
		log.Fatalf("RunCreatePlan should not get here")
	}
	response, err = common.PostYAMLDocument(string(yamlData), common.CurrentContext.URL + "plan")
	if err != nil {
		log.Printf("Error posting: %v\n", err)
		return
	}
	log.Printf(response)
}
