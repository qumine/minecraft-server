package api

import (
	"context"
	"net/http"
	"sync"

	"github.com/qumine/qumine-server-java/internal/wrapper"
	"github.com/sirupsen/logrus"
)

// API represents the api server
type API struct {
	Wrapper *wrapper.Wrapper

	httpServer *http.Server
}

// NewAPI creates a new api instance with the given host and port
func NewAPI(w *wrapper.Wrapper) *API {
	r := http.NewServeMux()

	api := &API{
		Wrapper: w,

		httpServer: &http.Server{
			Addr:    "0.0.0.0:8080",
			Handler: r,
		},
	}
	r.HandleFunc("/health/live", api.healthLive)
	r.HandleFunc("/health/ready", api.healthReady)

	r.HandleFunc("/command", api.command)

	return api
}

// Start the Api
func (a *API) Start(ctx context.Context, wg *sync.WaitGroup) {
	logrus.WithField("addr", a.httpServer.Addr).Info("starting api")

	go func() {
		wg.Add(1)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("starting api failed")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			a.Stop(wg)
			return
		}
	}
}

// Stop the api
func (a *API) Stop(wg *sync.WaitGroup) {
	logrus.Info("stopping api")
	if err := a.httpServer.Close(); err != nil {
		logrus.WithError(err).Error("stopping api failed")
	}
	logrus.Info("stopped api")
	wg.Done()
}
