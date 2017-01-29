package proxy

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/kcmerrill/automagicproxy/endpoint"
	"github.com/kcmerrill/shutdown.go"
	"net/http"
	"rsc.io/letsencrypt"
	"sort"
	"strings"
	"time"
)

/* Store all of our endpoints */
var endpoints map[string]*endpoint.Endpoint
var endpointkeys sort.StringSlice

/* Meat and potatoes right here ... */
func passThrough(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-AutomagicProxy", "v1.0")
	usekey := "_notfound"

	/* Grab the first key in the list that matches */
	for _, key := range endpointkeys {
		b := endpoints[key].Registered

		/* Allow for multiple containers with the same url */
		if strings.Contains(b, "_") {
			b = b[0:strings.Index(b, "_")]
		}

		if strings.HasPrefix(r.Host, b) && endpoints[key].Active {
			usekey = key
			break
		}
	}

	log.WithFields(
		log.Fields{
			"Request":   r.Host,
			"IP":        r.RemoteAddr,
			"Forwarded": usekey,
		}).Info("New Request")

	/* One quick sanity check before sending it on it's way */
	if _, exists := endpoints[usekey]; exists {
		endpoints[usekey].Proxy.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Error 502 - Bad Gateway"))
	}
}

/* Starts our proxy .. */
func Start(http_port int, secured bool) {
	log.WithFields(
		log.Fields{
			"port": http_port,
		}).Info("Starting automagic proxy")

	/* Start our healthchecks */
	go HealthChecks()

	http.HandleFunc("/", passThrough)

	if !secured {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", http_port), nil); err != nil {
			log.Fatal(err.Error())
			shutdown.Now()
		}
	} else {
		var m letsencrypt.Manager
		if err := m.CacheFile("letsencrypt.cache"); err != nil {
			log.Fatal(err)
			shutdown.Now()
		}
		log.Fatal(m.Serve())
	}

}

/* Add an endpoint to our proxy */
func Add(base, endpoint_url string) error {
	/* Check if endpoint already exists */
	for _, item := range endpoints {
		if item.Registered == base && item.Url.String() == endpoint_url {
			return nil
		}
	}

	/* Construct the key so that you can sort by url base and time added */
	urlbase := base

	/* Remove any thing after the _ from the url */
	if strings.Contains(urlbase, "_") {
		urlbase = urlbase[0:strings.Index(urlbase, "_")]
	}

	key := urlbase + "-" + time.Now().Format("2006-01-02T15:04:05.000")

	/* Add new endpoint */
	if ep, err := endpoint.New(base, endpoint_url); err == nil {
		/* If it doesn't exist ... */
		log.WithFields(log.Fields{
			"url":        endpoint_url,
			"registered": base,
			"urlbase":    urlbase,
		}).Info("Registered endpoint")
		endpoints[key] = ep
		endpointkeys = append(endpointkeys, key)

		sort.Sort(sort.Reverse(endpointkeys))

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
