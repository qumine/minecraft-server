package yatopia

import (
	"github.com/qumine/qumine-server-java/internal/server/common"
	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	serverVersion     = "latest"
	serverForceUpdate = false
	serverYatopiaAPI  = "https://api.yatopiamc.org/v2/latestBuild"
)

// Server is the struct for yatopia servers.
type Server struct {
	serverVersion     string
	serverForceUpdate bool
	yatopiaAPI        string
}

// NewYatopiaServer creates a new yatopia server.
func NewYatopiaServer() *Server {
	return &Server{
		serverVersion:     utils.GetEnvString("SERVER_VERSION", serverVersion),
		serverForceUpdate: utils.GetEnvBool("SERVER_FORCE_UPDATE", serverForceUpdate),
		yatopiaAPI:        utils.GetEnvString("SERVER_YATOPIA_API", serverYatopiaAPI),
	}
}

// Configure configures the server.
func (s *Server) Configure() error {
	logrus.WithFields(logrus.Fields{
		"type":        "YATOPIA",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"yatopiaApi":  s.yatopiaAPI,
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
		"type":        "YATOPIA",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"yatopiaApi":  s.yatopiaAPI,
	}).Info("configured server")
	return nil
}

// Update updates the resource, if supported uses cache.
func (s *Server) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":        "YATOPIA",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"yatopiaApi":  s.yatopiaAPI,
	}).Debug("updating server")

	versionDetailsDownloadURL := s.yatopiaAPI
	if s.serverVersion != "latest" {
		versionDetailsDownloadURL = versionDetailsDownloadURL + "?branch=ver/" + s.serverVersion
	}

	versionDetails, err := getVersionDetails(versionDetailsDownloadURL)
	if err != nil {
		return err
	}

	if common.CompareHash(s.serverForceUpdate, versionDetails.Branch.Commit.Sha) {
		logrus.WithFields(logrus.Fields{
			"type":        "YATOPIA",
			"version":     s.serverVersion,
			"forceUpdate": s.serverForceUpdate,
			"yatopiaApi":  s.yatopiaAPI,
		}).Info("updating server skipped, jar seems up to date")
		return nil
	}

	if err := utils.DownloadToFile(versionDetails.DownloadURL, "yatopia.jar"); err != nil {
		return err
	}

	if err := common.SaveHash(versionDetails.Branch.Commit.Sha); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"type":        "YATOPIA",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"yatopiaApi":  s.yatopiaAPI,
	}).Info("updated server")
	return nil
}

// StartupCommand returns the command and arguments used to startup the server.
func (s *Server) StartupCommand() (string, []string) {
	return "java", []string{"-jar", "yatopia.jar", "nogui"}
}
