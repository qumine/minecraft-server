package custom

import (
	"os"

	"github.com/qumine/qumine-server-java/internal/server/common"
	"github.com/sirupsen/logrus"
)

// Server is the updater for custom servers.
type Server struct {
	customURL string
}

// NewCustomServer creates a new custom server.
func NewCustomServer() *Server {
	return &Server{
		customURL: os.Getenv("SERVER_CUSTOM_URL"),
	}
}

// Configure configures the server.
func (s *Server) Configure() error {
	logrus.WithFields(logrus.Fields{
		"type":      "CUSTOM",
		"customURL": s.customURL,
	}).Info("server configuring")

	if err := common.ConfigureEula(); err != nil {
		return err
	}

	if err := common.ConfigureOps(); err != nil {
		return err
	}

	if err := common.ConfigureWhitelist(); err != nil {
		return err
	}

	if err := common.ConfigureServerProperties(); err != nil {
		return err
	}

	logrus.Debug("server configured")
	return nil
}

// Update updates the resource, if supported uses cache.
func (s *Server) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":      "CUSTOM",
		"customURL": s.customURL,
	}).Info("server updating")

	if err := common.DownloadServerJar(s.customURL); err != nil {
		return err
	}

	logrus.Debug("server updated")
	return nil
}

// StartupCommand retuns the command and arguments used to startup the server.
func (s *Server) StartupCommand() (string, []string) {
	return "java", []string{"-jar", "server.jar", "nogui"}
}
