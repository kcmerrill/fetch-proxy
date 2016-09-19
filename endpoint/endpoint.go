package endpoint

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var Healthy bool=true

/* Everything we need to house our endpoints */
type Endpoint struct {
	Active     bool
	Url        *url.URL
	Proxy      *httputil.ReverseProxy
	Registered string
	Available  time.Time
}

func (e Endpoint) isActive() bool {
	return e.Active
}

func (e *Endpoint) HealthCheck() {
	previous_status := e.isActive()
	status_code := 500
	if resp, err := http.Get(e.Url.String()); err != nil {
		/* Something is up ... disable this endpoint */
		e.Active = false
	} else {
		/* Woot! Good to go ... */
		status_code = resp.StatusCode
		if Healthy && resp.StatusCode >= 500 {
			e.Active = false
		} else {
			e.Active = true
		}
	}
	log.WithFields(log.Fields{
		"previous": previous_status,
		"current":  e.Active,
	}).Debug(e.Registered)
	if e.Active != previous_status {
		if e.Active {
			/* Whew, we came back online */
			log.WithFields(
				log.Fields{
					"URL": e.Url.String(),
				}).Info("Up")
			e.Available = time.Now()
		} else {
			/* BOO HISS! */
			log.WithFields(
				log.Fields{
					"URL":    e.Url.String(),
					"Status": status_code,
				}).Error("Down")
		}
	}
}

/* Create a new endpoint! */
func New(base, endpoint_url string) (*Endpoint, error) {
	if u, err := url.Parse(endpoint_url); err != nil {
		log.WithFields(log.Fields{"url": u}).Error("Problem parsing URL")
		return nil, errors.New("Problem parsing URL " + endpoint_url)
	} else {
		e := &Endpoint{
			Url:        u,
			Proxy:      httputil.NewSingleHostReverseProxy(u),
			Active:     true,
			Available:  time.Now(),
			Registered: base,
		}
		e.HealthCheck()
		return e, nil
	}
}
