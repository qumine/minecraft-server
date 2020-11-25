package api

import (
	"net/http"
	"strings"
)

func (a *API) healthLive(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(a.statusToHttpStatus(a.Wrapper.Status))
	w.Write([]byte(a.Wrapper.Status))
}

func (a *API) healthReady(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(a.statusToHttpStatus(a.Wrapper.Status))
	w.Write([]byte(a.Wrapper.Status))
}

func (a *API) statusToHttpStatus(status string) int {
	switch strings.ToUpper(status) {
	case "STARTING":
	case "STOPPING":
	case "STOPPED":
		return 503
	case "STARTED":
		return 200
	}
	return 500
}
