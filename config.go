package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Config holds all of our configuration values.
type Config struct {
	ForwardMap     map[string]string `yaml:"forward"`
	containerized  bool              `yaml:"-"`
	healthChecks   bool              `yaml:"-"`
	healthCheckURL string            `yaml:"-"`
}

// Forward the mappings to the appropiate port
func (c *Config) Forward() {
	address := "localhost"
	if c.containerized {
		address = containerizedIP()
	}

	for base, port := range c.ForwardMap {
		// Go ahead and add it
		AddSite(base, fmt.Sprintf("http://%s:%s", address, port), c.healthChecks, c.healthCheckURL)
	}
}

// ConfigWatch will watch a configfile for any changes and update it's mappings
func ConfigWatch(file string, containerized, healthChecks bool, healthCheckURL string) {
	for {
		c := &Config{
			healthChecks:   healthChecks,
			healthCheckURL: healthCheckURL,
			containerized:  containerized,
		}
		// make sure the file exists first
		if _, exists := os.Stat(file); exists == nil {
			contents, readErr := ioutil.ReadFile(file)
			if readErr == nil {
				yaml.Unmarshal(contents, c)
				c.Forward()
			}
		}

		// Pause. Rinse. Repeat
		<-time.After(5 * time.Second)
	}
}
