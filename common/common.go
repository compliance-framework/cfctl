package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"net/http"
    "strings"

    "github.com/spf13/cobra"

	"gopkg.in/yaml.v2"
)

type Config struct {
	FilePath string
    URL      string
    Context  string
    Output   string
}

type CLIConfig struct {
    Default  string             `yaml:"default"`
    Contexts map[string]Context `yaml:"contexts"`
}

type Context struct {
    URL      string `yaml:"url"`
}

var RunConfig      Config
var CLIConfigVar   CLIConfig
var CurrentContext Context

func ValidateYAMLDocument(doc string) error {
	var payload map[string]interface{}
	return yaml.Unmarshal([]byte(doc), &payload)
}

func PostYAMLDocument(doc string, url string) error {
	// Create the request body
	body := bytes.NewBuffer([]byte(doc))

	// Create the request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/yaml")

	// Log the request details
	log.Printf("Sending %s\n", url)
	log.Printf("Headers: %v\n", req.Header)
	log.Printf("Body: %s\n", doc)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read and log the response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}
	log.Printf("Response Status: %s\n", resp.Status)
	log.Printf("Response Body: %s\n", responseBody)

	return nil
}

func SanitiseContext() {
    url := string(CurrentContext.URL)

    if len(url) > 0 && url[len(url)-1] != '/' {
        CurrentContext.URL += "/"
    }
}

func ReadConfigFile() {
    configFile := os.Getenv("CFCONF")
    if configFile == "" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            log.Fatalf("error getting home directory: %v", err)
        }
        configFile = filepath.Join(homeDir, ".cfcli", "config")
    }

    fileContent, err := ioutil.ReadFile(configFile)
    if err != nil {
        log.Printf("error reading config file: %v", err)
        log.Fatalf(`Create a file, here is a simple example which works for development of CF:

default: dev
contexts:
  dev:
    url: http://localhost:8080`)
    }

    err = yaml.Unmarshal(fileContent, &CLIConfigVar)
    if err != nil {
        log.Fatalf("error unmarshaling config file: %v", err)
    }
}

func ApplyContext() {
    contextName := RunConfig.Context
    if contextName == "" {
        contextName = CLIConfigVar.Default
    }

    var exists bool
    CurrentContext, exists = CLIConfigVar.Contexts[contextName]
    if !exists {
        log.Fatalf("context %s not found in config file", contextName)
    }

    SanitiseContext()
}

func RunValidate(cmd *cobra.Command, args []string) {
    fileContent, err := ioutil.ReadFile(RunConfig.FilePath)
    if err != nil {
        log.Fatalf("error reading file: %v", err)
    }

    yamlDocuments := strings.Split(string(fileContent), "---")

    for _, doc := range yamlDocuments {
        doc = strings.TrimSpace(doc)
        if doc != "" {
            err := ValidateYAMLDocument(doc)
            if err != nil {
                log.Printf("Invalid YAML document: %v", err)
            } else {
                fmt.Println("Valid YAML document")
            }
        }
    }
}
