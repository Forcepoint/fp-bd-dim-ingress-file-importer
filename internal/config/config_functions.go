package config

import (
	"fp-dim-aws-guard-duty-ingress/internal/filehandler"
	"fp-dim-aws-guard-duty-ingress/internal/structs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	configFilePath = "./config/"
	listFilePath   = "./lists/"
	ConfigFile     = configFilePath + "config.yml"
	UrlsFile       = listFilePath + "urls.txt"
	BlocklistFile  = listFilePath + "blocklist.txt"
)

func InitFiles() {
	// Check if required file exists. If not create.
	if _, err := os.Stat(ConfigFile); err != nil {
		filehandler.CreateFile(ConfigFile)
	}

	if _, err := os.Stat(UrlsFile); err != nil {
		filehandler.CreateFile(UrlsFile)
	}

	if _, err := os.Stat(BlocklistFile); err != nil {
		filehandler.CreateFile(BlocklistFile)
	}
	return
}

func InitConfig() {
	// Initialise viper config using config.yml file.
	viper.SetConfigName("config")
	viper.AddConfigPath(configFilePath)
	viper.WatchConfig()

	// If unable to read in config, exit.
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal("There was an error while trying to read the config.")
	}
}

func WriteConfig(hours int) {
	viper.Set("update-schedule", hours)
	if err := viper.WriteConfig(); err != nil {
		logrus.Error("Error writing last run time to config file")
	}
	return
}

func GetConfig() (config structs.ModuleConfig) {
	config.Fields = []structs.Element{
		{
			Label:            "URL List",
			Type:             structs.TextArea,
			ExpectedJsonName: "lists",
			Rationale:        "Enter the URL of one or more intelligence sources (separated by comma).",
			Value:            strings.Join(filehandler.ReadFileContents(filehandler.OpenFile(UrlsFile, os.O_RDONLY)), ",\n"), // value of urls in urls.txt
			PossibleValues:   nil,
			Required:         false,
		}, {
			Label:            "Automated Update Schedule",
			Type:             structs.Number,
			ExpectedJsonName: "updater_schedule",
			Rationale:        "Enter how often the importer will import from sources (hours).",
			Value:            viper.GetString("update-schedule"),
			PossibleValues:   []string{"0", "100000"},
			Required:         true,
		}, {
			Label:            "File Upload",
			Type:             structs.FileUpload,
			ExpectedJsonName: "",
			Rationale:        `Import intelligence manually if sources are not reachable via the network.`,
			Value:            "",
			PossibleValues:   nil,
			Required:         false,
		}, {
			Label:            "Requirements",
			Type:             structs.Info,
			ExpectedJsonName: "",
			Rationale:        "URL List: CSV and TXT\n\nClick the Help icon for further information on how to configure this module.",
			Value:            "",
			PossibleValues:   nil,
			Required:         false,
		}, {
			Label:            "Elements Imported",
			Type:             structs.Info,
			ExpectedJsonName: "",
			Rationale:        "IP Addresses\n IP Ranges\n URLs\n Domains.",
			Value:            "",
			PossibleValues:   nil,
			Required:         false,
		},
	}
	return
}
