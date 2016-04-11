package docker

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	d "github.com/fsouza/go-dockerclient"
	"github.com/kcmerrill/automagicproxy/proxy"
	"os/exec"
	"strings"
	"time"
)

var url = ""
var endpoint = "unix:///var/run/docker.sock"

func Start(auto bool, my_port int) {
	my_port_64 := int64(my_port)
	if auto {
		url = containerized()
	} else {
		url = "localhost"
	}

	if client, err := d.NewClient(endpoint); err != nil {
		log.Error(err.Error())
		return
	} else {
		for {
			<-time.After(5 * time.Second)
			if containers, err := client.ListContainers(d.ListContainersOptions{All: true}); err == nil {
				for _, container := range containers {
					name := container.Names[0][1:]
					for _, port := range container.Ports {
						if port.PublicPort != my_port_64 && port.PrivatePort == 80 && port.Type == "tcp" {
							log.WithFields(log.Fields{
								"Public":  port.PublicPort,
								"Private": port.PrivatePort,
							}).Debug(name)
							proxy.Add(name, fmt.Sprintf("http://%s:%d", url, port.PublicPort))
						}
					}
				}
			} else {
				log.WithFields(log.Fields{
					"docker": endpoint,
				}).Error("Unable to connect to docker")
			}
		}
	}
}

func containerized() string {
	/* Do we need to start our docker service? */
	cmd := exec.Command("bash", "-c", "/sbin/ip route|awk '/default/ { print $3 }'")
	if output, err := cmd.Output(); err == nil {
		log.WithFields(log.Fields{
			"IP": strings.TrimSpace(string(output)),
		}).Info("Auto detecting docker host IP Address")
		return fmt.Sprintf("%s", strings.TrimSpace(string(output)))
	}

	return "localhost"
}
