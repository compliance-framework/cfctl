package common

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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
	URL string `yaml:"url"`
}

var RunConfig Config
var CLIConfigVar CLIConfig
var CurrentContext Context

func PostYAMLDocument(doc string, url string) (string, error) {
	// Create the request body
	body := bytes.NewBuffer([]byte(doc))

	// Create the request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
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
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %w\n", err)
		os.Exit(1)
	}
	//fmt.Printf("Response Status: %s\n", resp.Status)
	//fmt.Printf("Response Body: %s\n", responseBody)

	return string(responseBody), nil
}

func PutURL(url string) (string, error) {
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: received status code %d\n", resp.StatusCode)
		os.Exit(1)
	}

	fmt.Println("Request was successful.")
	return "200 OK", nil
}

func SanitiseContext() {
	url := string(CurrentContext.URL)

	if len(url) > 0 && url[len(url)-1] != '/' {
		CurrentContext.URL += "/"
	}
}
func GetConfigFilePath() (string, error) {
	configFile := os.Getenv("CFCONF")
	if configFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configFile = filepath.Join(homeDir, ".cfctl", "config")
	}
	return configFile, nil
}
func ReadConfigFile() {
	configFile, err := GetConfigFilePath()
	if err != nil {
		os.Exit(1)
	}
	fileContent, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		fmt.Printf(`Create a file, here is a simple example which works for development of CF:

default: dev
contexts:
  dev:
    url: http://localhost:8080/api`)
	}

	err = yaml.Unmarshal(fileContent, &CLIConfigVar)
	if err != nil {
		fmt.Printf("Error unmarshaling config file: %v\n", err)
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
		fmt.Printf("Context '%s' not found in config file.\n", contextName)
		os.Exit(1)
	}

	SanitiseContext()
}
