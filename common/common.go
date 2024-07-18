package common

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"net/http"
    "strings"
	"fmt"

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

func PostYAMLDocument(doc string, url string) (string, error) {
	// Create the request body
	body := bytes.NewBuffer([]byte(doc))

	// Create the request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Printf("error creating HTTP request: %v", err)
		return "", err
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/yaml")

	// fmt the request details
	//fmt.Printf("Sending %s\n", url)
	//fmt.Printf("Headers: %v\n", req.Header)
	//fmt.Printf("Body: %s\n", doc)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error making HTTP request: %w", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read and fmt the response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response: %w", err)
		os.Exit(1)
	}
	//fmt.Printf("Response Status: %s\n", resp.Status)
	//fmt.Printf("Response Body: %s\n", responseBody)

	return string(responseBody), nil
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
			os.Exit(1)
        }
        configFile = filepath.Join(homeDir, ".cfctl", "config")
    }

    fileContent, err := ioutil.ReadFile(configFile)
    if err != nil {
        fmt.Printf("error reading config file: %v", err)
        fmt.Printf(`Create a file, here is a simple example which works for development of CF:

default: dev
contexts:
  dev:
    url: http://localhost:8080`)
    }

    err = yaml.Unmarshal(fileContent, &CLIConfigVar)
    if err != nil {
        fmt.Printf("error unmarshaling config file: %v", err)
		os.Exit(1)
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
        fmt.Printf("context %s not found in config file", contextName)
		os.Exit(1)
    }

    SanitiseContext()
}
