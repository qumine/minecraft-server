package api

import "net/http"

func (a *API) handleSendCommand(w http.ResponseWriter, req *http.Request) {
	if err := a.w.SendCommand(req.URL.Query().Get("command")); err != nil {
		w.Write([]byte("ok"))
	} else {
		w.Write([]byte("nok"))
	}
}
