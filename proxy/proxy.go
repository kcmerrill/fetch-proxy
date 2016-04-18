package proxy

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/kcmerrill/automagicproxy/endpoint"
	"github.com/kcmerrill/shutdown.go"
	"net/http"
	"strings"
	"time"
)

/* Store all of our endpoints */
var endpoints map[string]*endpoint.Endpoint

/* Meat and potatoes right here ... */
func passThrough(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-AutomagicProxy", "v1.0")
	use := "_default"

	/* Grab the last key in the list that matches */
	for base, _ := range endpoints {
		b := base
		if strings.Contains(b, "_") {
			b = base[0:strings.Index(b, "_")]
		}
		if strings.HasPrefix(r.Host, b) && endpoints[base].Active {
			if endpoints[base].Available.After(endpoints[use].Available) {
				use = base
			}
		}
	}

	log.WithFields(
		log.Fields{
			"Request":   r.Host,
			"IP":        r.RemoteAddr,
			"Forwarded": use,
		}).Info("New Request")

	/* One quick sanity check before sending it on it's way */
	if _, exists := endpoints[use]; exists {
		endpoints[use].Proxy.ServeHTTP(w, r)
	}
}

/* Starts our proxy .. */
func Start(http_port int) {
	log.WithFields(
		log.Fields{
			"port": http_port,
		}).Info("Starting automagic proxy")

	/* Add our default */
	Add("_default", fmt.Sprintf("www.localhost:%d", http_port))

	/* Start our healthchecks */
	go HealthChecks()

	http.HandleFunc("/", passThrough)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", http_port), nil); err != nil {
		log.Error(err.Error())
		shutdown.Now()
	}
}

/* Add an endpoint to our proxy */
func Add(base, endpoint_url string) error {
	if _, exists := endpoints[base]; exists {
		if endpoints[base].Registered == base && endpoints[base].Url.String() == endpoint_url {
			return nil
		}
	}
	if ep, err := endpoint.New(base, endpoint_url); err == nil {
		/* If it doesn't exist ... */
		log.WithFields(log.Fields{
			"url":        endpoint_url,
			"registered": base,
		}).Info("Registered endpoint")
		endpoints[base] = ep
		return nil
	} else {
		return err
	}
}

func HealthChecks() {
	for {
		<-time.After(10 * time.Second)
		for key, _ := range endpoints {
			go endpoints[key].HealthCheck()
		}
	}
}

/* Get our inits out of the way ... */
func init() {
	endpoints = make(map[string]*endpoint.Endpoint)
}
