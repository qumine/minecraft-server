package api

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/qumine/qumine-server-java/internal/wrapper"
	"github.com/sirupsen/logrus"
)

// API represents the api server
type API struct {
	Wrapper *wrapper.Wrapper

	httpServer *http.Server
}

// NewAPI creates a new api instance
func NewAPI(w *wrapper.Wrapper) *API {
	r := http.NewServeMux()

	api := &API{
		Wrapper: w,

		httpServer: &http.Server{
			Addr:    net.JoinHostPort(utils.GetEnvString("HTTP_ADDR", "0.0.0.0"), utils.GetEnvString("HTTP_PORT", "8080")),
			Handler: r,
		},
	}
	r.HandleFunc("/health/live", api.healthLive)
	r.HandleFunc("/health/ready", api.healthReady)

	return api
}

// Start the Api
func (a *API) Start(ctx context.Context, wg *sync.WaitGroup) {
	logrus.WithFields(logrus.Fields{
		"addr": a.httpServer.Addr,
	}).Debug("starting api")

	go func() {
		wg.Add(1)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithFields(logrus.Fields{
				"addr": a.httpServer.Addr,
			}).Fatal("starting api failed")
		}
	}()

	logrus.WithFields(logrus.Fields{
		"addr": a.httpServer.Addr,
	}).Info("started api")
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
	logrus.WithFields(logrus.Fields{
		"addr": a.httpServer.Addr,
	}).Info("stopping api")

	if err := a.httpServer.Close(); err != nil {
		logrus.WithFields(logrus.Fields{
			"addr": a.httpServer.Addr,
		}).Error("stopping api failed")
	}

	logrus.WithFields(logrus.Fields{
		"addr": a.httpServer.Addr,
	}).Info("stopped api")
	wg.Done()
}
