package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []propertySource `json:"propertySources"`
}

type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

func LoadConfigurationFromBranch(configServerURL string, appName string, profile string, branch string) {
	url := fmt.Sprintf("%s/%s/%s/%s", configServerURL, appName, profile, branch)
	logrus.Infof("Loading config from %s\n", url)
	body, err := fetchConfiguration(url)
	if err != nil {
		panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	}

	parseConfiguration(body)
}

func fetchConfiguration(url string) ([]byte, error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Infof("Recovered in %v\n", r)
		}
	}()

	logrus.Infof("Getting config from %v\n", url)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		message := fmt.Sprintf("Couldn't load configuration, cannot start. Terminating. Error: %v", err.Error())
		logrus.Errorln(message)
		panic(message)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		message := fmt.Sprintf("Error reading configuration: %v", err.Error())
		logrus.Errorln(message)
		panic(message)
	}
	return body, err
}

func parseConfiguration(body []byte) {
	var cloudConfig springCloudConfig
	err := json.Unmarshal(body, &cloudConfig)
	if err != nil {
		message := fmt.Sprintf("Cannot parse configuration, message: %v", err.Error())
		logrus.Errorln(message)
		panic(message)
	}

	for key, value := range cloudConfig.PropertySources[0].Source {
		viper.Set(key, value)
		logrus.Infof("Loading config property %v => %v\n", key, value)
	}

	if viper.IsSet("server_name") {
		logrus.Infof("Successfully loaded configuration for service %s\n", viper.GetString("server_name"))
	}
}
