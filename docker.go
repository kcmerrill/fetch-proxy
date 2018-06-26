package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

var portPriority = map[uint16]int{
	80:   1, // HTTP
	81:   2, // HTTP
	8080: 3, // HTTP
	3000: 4, // HTTP
	5000: 5, // HTTP
}

// ContainerWatch checks for new containers and if they exist, add the sites and it's endpoints
func ContainerWatch(containerized, healthchecks bool, healthCheckURL string, myPort int) {
	myPort64 := uint16(myPort)
	address := "localhost"
	if containerized {
		address = containerizedIP()
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Error(err.Error())
		return
	}

	for {
		if containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{}); err == nil {
			for _, container := range containers {
				name := container.Names[0][1:]
				var svcPort types.Port
				for _, port := range container.Ports {
					// Prioritize port selection
					if port.PublicPort != myPort64 && port.PublicPort != 443 && port.Type == "tcp" && port.PrivatePort >= 80 {
						if 0 == svcPort.PrivatePort {
							log.WithFields(log.Fields{
								"Container": name,
								"Public":    port.PublicPort,
								"Private":   port.PrivatePort,
							}).Debug("initial port selected")
							svcPort = port
						}

						priority1, ok := portPriority[port.PrivatePort]
						priority2, _ := portPriority[svcPort.PrivatePort]
						if ok && priority1 > 0 && (0 == priority2 || priority1 < priority2) {
							log.WithFields(log.Fields{
								"Container": name,
								"Public":    port.PublicPort,
								"Private":   port.PrivatePort,
							}).Debug("prioritized port selected")
							svcPort = port
						}
					}
				}
				AddSite(name, fmt.Sprintf("http://%s:%d", address, svcPort.PublicPort), healthchecks, healthCheckURL)
			}
		} else {
			log.Error("Unable to connect to docker")
			log.Error(err.Error())
		}

		// Every 5 seconds, check for new containers
		<-time.After(5 * time.Second)
	}
}

// containerizedIP returns a string with the ip address of the docker host. Localhost else ...
func containerizedIP() string {
	// Do we need to start our docker service?
	cmd := exec.Command("bash", "-c", "/sbin/ip route|awk '/default/ { print $3 }'")
	if output, err := cmd.Output(); err == nil {
		log.WithFields(log.Fields{
			"IP": strings.TrimSpace(string(output)),
		}).Info("Auto detecting docker host IP Address")
		return fmt.Sprintf("%s", strings.TrimSpace(string(output)))
	}
	return "localhost"
}
