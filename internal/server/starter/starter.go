package starter

import (
	"os"

	"github.com/qumine/qumine-server-java/internal/server/common"
	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

// Server is the updater for starter servers.
type Server struct {
	starterJarURL    string
	starterConfigURL string
}

// NewStarterServer creates a new starter server.
func NewStarterServer() *Server {
	return &Server{
		starterJarURL:    os.Getenv("SERVER_STARTER_JAR_URL"),
		starterConfigURL: os.Getenv("SERVER_STARTER_CONFIG_URL"),
	}
}

// Configure configures the server.
func (s *Server) Configure() error {
	logrus.WithFields(logrus.Fields{
		"type":             "STARTER",
		"starterJarURL":    s.starterJarURL,
		"starterConfigURL": s.starterConfigURL,
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
		"type":             "STARTER",
		"starterJarURL":    s.starterJarURL,
		"starterConfigURL": s.starterConfigURL,
	}).Info("configured server")
	return nil
}

// Update updates the resource, if supported uses cache.
func (s *Server) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":             "STARTER",
		"starterJarURL":    s.starterJarURL,
		"starterConfigURL": s.starterConfigURL,
	}).Debug("updating server")

	if err := utils.DownloadToFile(s.starterJarURL, "starter.jar"); err != nil {
		return err
	}
	if err := utils.DownloadToFile(s.starterConfigURL, "server-setup-config.yaml"); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"type":             "STARTER",
		"starterJarURL":    s.starterJarURL,
		"starterConfigURL": s.starterConfigURL,
	}).Info("updated server")
	return nil
}

// StartupCommand retuns the command and arguments used to startup the server.
func (s *Server) StartupCommand() (string, []string) {
	return "java", []string{"-jar", "starter.jar", "nogui"}
}
