package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Configuration struct {
	MinPasswordLength int    `json:"min_password_length"` // min length of generated passwords
	LogFile           string `json:"log_file"`            // full path to pwm log file
}

// global variables
var Config Configuration

// String returns string representation of dbs Config
func (c *Configuration) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("ERROR: fail to marshal configuration", err)
		return ""
	}
	return string(data)
}

func ParseConfig(configFile string) error {
	// if config file does not exists we'll create one
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// make config dir
		dir, _ := path.Split(configFile)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("unable to create PWM area %s for config file, error %v", dir, err)
		}
		lfile := fmt.Sprintf("%s/pwm.log", pwmHome())
		config := Configuration{MinPasswordLength: 24, LogFile: lfile}
		data, err := json.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(configFile, data, 0755)
		if err != nil {
			log.Fatalf("unable to create PWM config file, error %v", err)
		}
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("unable to read config file", configFile, err)
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Println("unable to parse config file", configFile, err)
		return err
	}
	if Config.MinPasswordLength == 0 {
		Config.MinPasswordLength = 24
	}
	if Config.LogFile == "" {
		Config.LogFile = fmt.Sprintf("%s/pwm.log", pwmHome())
	}
	return nil
}
