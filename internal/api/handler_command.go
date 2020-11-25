package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func (a *API) command(w http.ResponseWriter, req *http.Request) {
	if err := a.Wrapper.Console.SendCommand(req.URL.Query().Get("command")); err != nil {
		logrus.WithError(err)
		w.Write([]byte("nok"))
	} else {
		w.Write([]byte("ok"))
	}
}
