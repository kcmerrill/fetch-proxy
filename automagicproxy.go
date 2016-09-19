package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	d "github.com/kcmerrill/automagicproxy/docker"
	"github.com/kcmerrill/automagicproxy/proxy"
	"github.com/kcmerrill/automagicproxy/endpoint"
	"github.com/kcmerrill/shutdown.go"
	"net/http"
	"syscall"
	"time"
)

func main() {
	/* Setup some command line arguments */
	port := flag.Int("port", 80, "The port in which the proxy will listen on")
	dockerized := flag.Bool("dockerized", false, "Query running containers and auto map them")
	containerized := flag.Bool("containerized", false, "Is automagicproxy running in a container?")
	insecure := flag.Bool("insecure", false, "Should use HTTP or HTTPS? HTTP works great for dev envs")
	unhealthy := flag.Bool("unhealthy", false, "Disable health checks for dev envs")
	dev := flag.Bool("dev", false, "Disable health checks and HTTPS for dev envs")
	flag.Parse()

	/* Disable ssl/tls and health checks in dev mode */
	if *dev {
		*unhealthy = true
		*insecure = true
	}

	endpoint.Healthy = !*unhealthy

	/* Set a global timeout */
	http.DefaultClient.Timeout = 10 * time.Second

	/* Start our proxy on the specified port */
	go proxy.Start(*port, !*insecure)

	if *dockerized {
		go d.Start(*containerized, *port)
	}

	/* No need to shutdown the application _UNLESS_ we catch it */
	shutdown.WaitFor(syscall.SIGINT, syscall.SIGTERM)
	log.Info("Shutting down ... ")
}
