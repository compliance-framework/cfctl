package validate

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
)

type Validate struct {
	FilePath string `json:"file-path" json:"file-path"`
}

var ValidateVar Validate

func RunValidate(cmd *cobra.Command, args []string) {
    fileContent, err := ioutil.ReadFile(ValidateVar.FilePath)
    if err != nil {
        fmt.Printf("error reading file: %v", err)
        os.Exit(1)
    }

    yamlDocuments := strings.Split(string(fileContent), "---")

    for _, doc := range yamlDocuments {
        doc = strings.TrimSpace(doc)
        if doc != "" {
            err := ValidateYAMLDocument(doc)
            if err != nil {
                fmt.Printf("Invalid YAML document: %v", err)
                os.Exit(1)
            } else {
                fmt.Println("Valid YAML document")
            }
        }
    }
}

func ValidateYAMLDocument(doc string) error {
    var payload map[string]interface{}
    return yaml.Unmarshal([]byte(doc), &payload)
}
