package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"golang.org/x/exp/errors"
	yaml "gopkg.in/yaml.v2"
)

// helper function to produce UTC time prefixed output
func logMsg(data []byte) string {
	return fmt.Sprintf("[" + time.Now().String() + "] " + string(data))
}

// custom rotate logger
type RotateLogWriter struct {
	RotateLogs *rotatelogs.RotateLogs
}

func (w RotateLogWriter) Write(data []byte) (int, error) {
	return w.RotateLogs.Write([]byte(logMsg(data)))
}

// Configuration represents configuration options for ECM application
type Configuration struct {
	LogFile string `json:"log_file" yaml:"LogFile"`
	Verbose int    `json:"verbose" yaml:"Verbose"`
}

// Config represents global configuration object
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

// ParseConfig parses given configuration file and initialize Config object
func ParseConfig(configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("unable to read config file", configFile, err)
		return err
	}
	// try json first
	err = json.Unmarshal(data, &Config)
	if err != nil {
		jsonErr := err
		// if fail try yaml
		err = yaml.Unmarshal(data, &Config)
		if err != nil {
			log.Println("unable to parse config file", configFile)
			log.Println("JSON failure", jsonErr)
			log.Println("YAML failure", err)
			return errors.New("Fail to parse config file")
		}
	}

	log.SetFlags(0)
	if Config.Verbose > 0 {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	if Config.LogFile != "" {
		rl, err := rotatelogs.New(Config.LogFile + "-%Y%m%d")
		if err == nil {
			rotlogs := RotateLogWriter{RotateLogs: rl}
			log.SetOutput(rotlogs)
		} else {
			log.Println("ERROR: unable to get rotatelogs", err)
		}
	}
	return nil
}
