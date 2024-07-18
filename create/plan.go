package create

import (
	"fmt"
	"io/ioutil"
	"os"

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
		fmt.Printf("title or yaml filename must be set")
		os.Exit(1)
	} else if planTitle != "" && fileName != "" {
		fmt.Printf("title and yaml cannot both be set")
		os.Exit(1)
	}

	// Process the command
	plan := domain.Plan{}
	if fileName != "" {
		yamlData, err = ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Printf("error reading file: %v", err)
			os.Exit(1)
		}
		// Check the yaml is valid vs domain.Plan
		_, err = yaml.Marshal(&plan)
		if err != nil {
			fmt.Printf("Error marshalling to YAML: %v\n", err)
			os.Exit(1)
		}
	} else if planTitle != "" {
		plan.Title = planTitle
		yamlData, err = yaml.Marshal(&plan)
		if err != nil {
			fmt.Printf("Error marshalling to YAML: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("RunCreatePlan should not get here")
		os.Exit(1)
	}
	response, err = common.PostYAMLDocument(string(yamlData), common.CurrentContext.URL + "plan")
	if err != nil {
		fmt.Printf("Error posting: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(response)
}
