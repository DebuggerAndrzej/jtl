package backend

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/DebuggerAndrzej/jtl/backend/entities"
)

func GetTomlConfig(configPath string) *entities.Config {
	var config entities.Config

	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		panic(fmt.Sprintf("Couldn't load config file. Check if %s file exists!", configPath))
	}

	return &config
}

func AddIssueToConfig(configPath, issueNumber string, config *entities.Config) {
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Coldn't load config file: %s", configPath))
	}

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		if strings.Contains(line, "additionalIssues") {
			lines[i] = line[:len(line)-1] + fmt.Sprintf(",FCA_OAM-%s\"", issueNumber)
		}

		config.AdditionalIssues += fmt.Sprintf(",FCA_OAM-%s", issueNumber)
	}
	output := strings.Join(lines, "\n")
	err = os.WriteFile(configPath, []byte(output), 0644)
	if err != nil {
		panic(err)
	}
}

func RemoveIssueFromConfig(configPath, issueId string, config *entities.Config) {
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Coldn't load config file: %s", configPath))
	}

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		if strings.Contains(line, "additionalIssues") {
			lines[i] = strings.Replace(line, fmt.Sprintf(",%s", issueId), "", 1)
		}

		config.AdditionalIssues = strings.Replace(config.AdditionalIssues, fmt.Sprintf(",%s", issueId), "", 1)
	}
	output := strings.Join(lines, "\n")
	err = os.WriteFile(configPath, []byte(output), 0644)
	if err != nil {
		panic(err)
	}
}
func getDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Couldn't determine user's home dir!")
	}
	return path.Join(homeDir, ".config/jtl.toml")
}
