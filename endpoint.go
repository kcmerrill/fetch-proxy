package main

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

// Endpoint struct containing everything needed for a new endpoint
type Endpoint struct {
	Active     bool
	Address    *url.URL
	Proxy      *httputil.ReverseProxy
	Registered string
	Available  time.Time
}

// isActive returns a bool if the endpiont is active or not
func (e Endpoint) isActive() bool {
	return e.Active
}

// HealthCheck performs a basic http check based on a positive(<500) status code
func (e *Endpoint) HealthCheck(healthCheckURL string) {
	previousStatus := e.isActive()
	statusCode := 500
	if resp, err := http.Get(e.Address.String() + "/" + healthCheckURL); err != nil {
		// Something is up ... disable this endpoint
		e.Active = false
	} else {
		// Woot! Good to go ...
		statusCode = resp.StatusCode
		if resp.StatusCode >= 500 {
			e.Active = false
		} else {
			e.Active = true
		}
	}
	log.WithFields(log.Fields{
		"previous": previousStatus,
		"current":  e.Active,
	}).Debug(e.Registered)
	if e.Active != previousStatus {
		if e.Active {
			// Whew, we came back online
			log.WithFields(
				log.Fields{
					"URL": e.Address.String(),
				}).Info("Up")
			e.Available = time.Now()
		} else {
			// BOO HISS!
			log.WithFields(
				log.Fields{
					"URL":    e.Address.String(),
					"Status": statusCode,
				}).Error("Down")
		}
	}
}

// NewEndpoint creates new endpoints to forward traffic to
func NewEndpoint(base, address string, checkHealth bool, healthURL string) (*Endpoint, error) {
	parsedAddress, err := url.Parse(address)
	if err != nil {
		log.WithFields(log.Fields{"url": parsedAddress}).Error("Problem parsing URL")
		return nil, errors.New("Problem parsing URL " + address)
	}
	e := &Endpoint{
		Address:    parsedAddress,
		Proxy:      httputil.NewSingleHostReverseProxy(parsedAddress),
		Active:     true,
		Available:  time.Now(),
		Registered: base,
	}
	if checkHealth {
		e.HealthCheck(healthURL)
	} else {
		// Make it active regardless. Godspeed developer :D
		e.Active = true
	}
	return e, nil
}
