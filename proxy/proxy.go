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
		if strings.HasPrefix(r.Host, base) {
			use = base
		}
	}

	/*If it exists ... use it ... */
	if _, exists := endpoints[use]; exists && endpoints[use].Active {
		log.WithFields(
			log.Fields{
				"Request":   r.Host,
				"IP":        r.RemoteAddr,
				"Forwarded": use,
			}).Info("New Request")
		endpoints[use].Proxy.ServeHTTP(w, r)
	} else {
		log.WithFields(
			log.Fields{
				"Request":   r.Host,
				"IP":        r.RemoteAddr,
				"Forwarded": use,
			}).Error("Endpoint not found")
	}
}

/* Starts our proxy .. */
func Start(http_port int) {
	go HealthChecks()
	log.WithFields(
		log.Fields{
			"port": http_port,
		}).Info("Starting automagic proxy")

	http.HandleFunc("/", passThrough)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", http_port), nil); err != nil {
		log.Error(err.Error())
		shutdown.Now()
	}
}

/* Add an endpoint to our proxy */
func Add(base, endpoint_url string) error {
	if ep, err := endpoint.New(base, endpoint_url); err == nil {
		if _, exists := endpoints[base]; !exists {
			log.WithFields(log.Fields{
				"incoming": base,
				"endpoint": endpoint_url,
			}).Info(fmt.Sprintf("%s regestered", base))
		}
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
