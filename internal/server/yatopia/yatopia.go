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
	}).Info("server configuring")

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
		"type":        "YATOPIA",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"yatopiaApi":  s.yatopiaAPI,
	}).Info("checking for server updates")

	if err := common.ConfigureEula(); err != nil {
		return err
	}

	versionDetailsDownloadURL := s.yatopiaAPI
	if s.serverVersion != "latest" {
		versionDetailsDownloadURL = versionDetailsDownloadURL + "?branch=ver/" + s.serverVersion
	}

	versionDetails, err := getVersionDetails(versionDetailsDownloadURL)
	if err != nil {
		return err
	}

	if common.CompareHash(s.serverForceUpdate, versionDetails.Branch.Commit.Sha) {
		logrus.Info("updated server")
		return nil
	}

	if err := common.DownloadServerJar(versionDetails.DownloadURL); err != nil {
		return err
	}

	if err := common.SaveHash(versionDetails.Branch.Commit.Sha); err != nil {
		return err
	}

	logrus.Info("updated server")
	return nil
}

// StartupCommand retuns the command and arguments used to startup the server.
func (s *Server) StartupCommand() (string, []string) {
	return "java", []string{"-jar", "server.jar", "nogui"}
}
