package create

import (
	"bufio"
	"cfctl/common"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Create context given by argument with user input
func AddContext(cmd *cobra.Command, args []string) {
	var ConfigFileContents common.CLIConfig

	defaultUrl := "http://localhost:8080/api"
	contextName := args[0]

	configFile, err := common.GetConfigFilePath()
	if err != nil {
		os.Exit(1)
	}

	fmt.Printf("File path: %v\n", configFile)
	// Create file if doesn't exist
	f, err := os.OpenFile(configFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()
	// Get the file's content
	content, err := io.ReadAll(f)
	if err != nil {
		os.Exit(1)
	}
	err = yaml.Unmarshal(content, &ConfigFileContents)
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("File contents: %v\n", ConfigFileContents)

	ConfigFileContents.Default = GetUserInputDefaultContext(ConfigFileContents.Default)

	inputUrl, err := GetUserInput("Enter url: ")
	if inputUrl == "" || err != nil {
		fmt.Printf("Using default URL %v\n", defaultUrl)
		inputUrl = defaultUrl
	}

	ConfigFileContents = CreateContext(ConfigFileContents, inputUrl, contextName)

	d, err := yaml.Marshal(ConfigFileContents)
	if err != nil {
		fmt.Printf("Could not marshal contents %v\n", ConfigFileContents)
		os.Exit(1)
	}

	// Overwrite the file with the new contents
	f.Truncate(0)
	f.Seek(0, 0)
	f.Write(d)
	fmt.Printf("Successfully wrote new context for '%v'\n", contextName)
}

func GetUserInputDefaultContext(currentDefault string) string {
	var display string
	if currentDefault == "" {
		display = "None"
	} else {
		display = currentDefault
	}
	inputContext, err := GetUserInput(fmt.Sprintf("Enter default context [%v]: ", display))
	if inputContext == "" || err != nil {
		if currentDefault == "" {
			fmt.Printf("Default context must be given on new config file creation.\n")
			os.Exit(1)
		} else {
			inputContext = currentDefault
		}
	}
	return inputContext
}

// Prompts user for input and gets value delimited, and exclusive of, newline
func GetUserInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	userInput, err := reader.ReadString('\n')
	userInput = userInput[:len(userInput)-1]
	return userInput, err
}

// Modifies or adds to the existing config if it exists, otherwise creates a new one
func CreateContext(existingConfig common.CLIConfig, inputUrl string, contextName string) common.CLIConfig {
	createContext := common.Context{URL: inputUrl}
	if len(existingConfig.Contexts) == 0 {
		fmt.Println("No context file found - creating new.")
		m := make(map[string]common.Context)
		m[contextName] = createContext
		existingConfig.Contexts = m
	} else {
		existingConfig.Contexts[contextName] = createContext
	}
	return existingConfig
}
