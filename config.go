package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// Configuration represents vault configuration structure
type Configuration struct {
	MinPasswordLength   int    `json:"min_password_length"`   // min length of generated passwords
	LogFile             string `json:"log_file"`              // full path to gpm log file
	TokenExpireInterval int    `json:"token_expire_interval"` // token expire interval in seconds
	TokenSecret         string `json:"token_secret"`          // token secret
}

// Config represents our vault configuration object
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

// ParseConfig provides config parsing
func ParseConfig(configFile string, verbose int) error {
	// if config file does not exists we'll create one
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// make config dir
		dir, _ := path.Split(configFile)
		log.Println("make dir", dir)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("unable to create GPM area %s for config file, error %v", dir, err)
		}
		lfile := fmt.Sprintf("%s/gpm.log", gpmHome())
		config := Configuration{MinPasswordLength: 24, LogFile: lfile}
		data, err := json.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(configFile, data, 0755)
		if err != nil {
			log.Fatalf("unable to create GPM config file, error %v", err)
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
	if Config.TokenExpireInterval == 0 {
		Config.TokenExpireInterval = 60
	}
	if Config.TokenSecret == "" {
		chars := voc + numbers + symbols
		Config.TokenSecret = generatePassword(24, chars)
	}
	if Config.LogFile == "" {
		Config.LogFile = fmt.Sprintf("%s/gpm.log", gpmHome())
	}

	// log time, filename, and line number
	if verbose > 0 {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	// setup logger
	log.SetOutput(new(LogWriter))
	if Config.LogFile != "" {
		logFile := Config.LogFile + "-%Y%m%d"
		rl, err := rotatelogs.New(logFile)
		if err == nil {
			rotlogs := RotateLogWriter{RotateLogs: rl}
			log.SetOutput(rotlogs)
		}
	}

	return nil
}
