package main

import (
	"flag"
	"net/http"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kcmerrill/shutdown.go"
)

func main() {
	// Setup some command line arguments
	port := flag.Int("port", 80, "The port in which the proxy will listen on")
	containerized := flag.Bool("containerized", false, "Is fetch-proxy running in a container?")
	insecure := flag.Bool("insecure", false, "Should use HTTP or HTTPS? HTTP works great for dev envs")
	disableHealthChecks := flag.Bool("disable-healthchecks", false, "disable health checks for dev envs")
	healthCheckURL := flag.String("healthcheck", "?health", "The url to be used for healthchecks")
	dev := flag.Bool("dev", false, "Disable health checks and HTTPS for dev envs")
	timeout := flag.Int("response-timeout", 10, "The response timeout for the proxy")
	defaultEndpoint := flag.String("default", "__default", "The default endpoint fetch-proxy uses when requested endpoing isn't found")
	config := flag.String("config", "", "Location for the configuration file you want to use")

	flag.Parse()

	// Disable ssl/tls and health checks in dev mode
	if *dev {
		*disableHealthChecks = true
		*insecure = true
	}

	// Set a global timeout
	http.DefaultClient.Timeout = time.Duration(*timeout) * time.Second

	// Start our proxy on the specified port
	go FetchProxyStart(*port, !*insecure, !*disableHealthChecks, *healthCheckURL, *defaultEndpoint)

	go ContainerWatch(*containerized, !*disableHealthChecks, *healthCheckURL, *port)

	if *config != "" {
		go ConfigWatch(*config, *containerized, !*disableHealthChecks, *healthCheckURL)
	}

	// No need to shutdown the application _UNLESS_ we catch it
	shutdown.WaitFor(syscall.SIGINT, syscall.SIGTERM)
	log.Info("Shutting down ... ")
}
