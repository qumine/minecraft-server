package custom

import (
	"os"
	"strings"

	"github.com/qumine/minecraft-server/internal/server/common"
	"github.com/qumine/minecraft-server/internal/utils"
	"github.com/sirupsen/logrus"
)

// Server is the updater for custom servers.
type Server struct {
	customURL string

	filename string
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
	}).Debug("configuring server")

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

	logrus.WithFields(logrus.Fields{
		"type":      "CUSTOM",
		"customURL": s.customURL,
	}).Info("configured server")
	return nil
}

// Update updates the resource, if supported uses cache.
func (s *Server) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":      "CUSTOM",
		"customURL": s.customURL,
	}).Debug("updating server")

	parts := strings.Split(s.customURL, "/")
	s.filename = parts[len(parts)-1]
	if err := utils.DownloadToFile(s.customURL, s.filename); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"type":      "CUSTOM",
		"customURL": s.customURL,
	}).Info("updated server")
	return nil
}

// UpdatePlugins updates the plugins, if supported.
func (s *Server) UpdatePlugins() error {
	logrus.WithFields(logrus.Fields{
		"type":      "CUSTOM",
		"customURL": s.customURL,
	}).Info("updating plugins failed, not supported by SERVER_TYPE")
	return nil
}

// StartupCommand returns the command and arguments used to startup the server.
func (s *Server) StartupCommand() (string, []string) {
	return utils.GetEnvString("SERVER_CUSTOM_COMMAND", "java"), utils.GetEnvStringList("SERVER_CUSTOM_ARGS", "-jar,"+s.filename+",nogui")
}
