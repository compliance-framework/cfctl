package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Define a structure to hold the CLI flags
type Config struct {
	FilePath string
	URL      string
}

var config Config

func main() {
	// Define the root command
	var rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "CLI for processing multiple YAML documents",
		Run:   run,
	}

	// Define flags
	rootCmd.PersistentFlags().StringVarP(&config.FilePath, "file", "f", "", "YAML file with multiple documents")
	rootCmd.PersistentFlags().StringVarP(&config.URL, "url", "u", "", "Endpoint URL to send the YAML payloads")

	// Mark required flags
	rootCmd.MarkPersistentFlagRequired("file")
	rootCmd.MarkPersistentFlagRequired("url")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Read the YAML file
	fileContent, err := ioutil.ReadFile(config.FilePath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	// Split the file content into multiple YAML documents
	yamlDocuments := strings.Split(string(fileContent), "---")

	// Process each YAML document
	for _, doc := range yamlDocuments {
		doc = strings.TrimSpace(doc)
		if doc != "" {
			processYAMLDocument(doc)
		}
	}
}

func processYAMLDocument(doc string) {
	// Unmarshal the YAML document to check for validity
	var payload map[string]interface{}
	err := yaml.Unmarshal([]byte(doc), &payload)
	if err != nil {
		log.Printf("error unmarshaling YAML document: %v", err)
		return
	}

	// Make an HTTP POST request with the YAML document as payload
	resp, err := http.Post(config.URL, "application/yaml", bytes.NewBuffer([]byte(doc)))
	if err != nil {
		log.Printf("error making HTTP request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response: %v", err)
		return
	}
	fmt.Printf("Response: %s\n", body)
}

