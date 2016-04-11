package endpoint

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

/* Everything we need to house our endpoints */
type Endpoint struct {
	Active bool
	Url    *url.URL
	Proxy  *httputil.ReverseProxy
}

func (e *Endpoint) HealthCheck() {
	previous_status := e.Active
	status_code := 0
	if resp, err := http.Get(e.Url.String()); err != nil {
		/* Something is up ... disable this endpoint */
		e.Active = false
		status_code = 500
	} else {
		/* Woot! Good to go ... */
		status_code = resp.StatusCode
		if resp.StatusCode < 500 {
			e.Active = true
		} else {
			e.Active = false
		}
	}

	if e.Active != previous_status {
		if e.Active {
			/* Whew, we came back online */
			log.WithFields(
				log.Fields{
					"URL": e.Url.String(),
				}).Debug("Up")
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
			Url:    u,
			Proxy:  httputil.NewSingleHostReverseProxy(u),
			Active: false,
		}
		e.HealthCheck()
		return e, nil
	}
}
