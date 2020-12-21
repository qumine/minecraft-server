package papermc

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/qumine/qumine-server-java/internal/server/common"
	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	serverVersion     = "latest"
	serverForceUpdate = false
	serverPapermcAPI  = "https://papermc.io/api/v2/projects/paper/"
)

// Server is the struct for papermc servers.
type Server struct {
	serverVersion     string
	serverForceUpdate bool
	papermcAPI        string
}

// NewPaperMCServer creates a new papermc server.
func NewPaperMCServer() *Server {
	return &Server{
		serverVersion:     utils.GetEnvString("SERVER_VERSION", serverVersion),
		serverForceUpdate: utils.GetEnvBool("SERVER_FORCE_UPDATE", serverForceUpdate),
		papermcAPI:        utils.GetEnvString("SERVER_PAPERMC_API", serverPapermcAPI),
	}
}

// Configure configures the server.
func (s *Server) Configure() error {
	logrus.WithFields(logrus.Fields{
		"type":        "PAPERMC",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"papermcAPI":  s.papermcAPI,
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
		"type":        "PAPERMC",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"papermcAPI":  s.papermcAPI,
	}).Info("server updating")

	version := ""
	if match, _ := regexp.MatchString("\\d*\\.\\d*\\.\\d", s.serverVersion); match {
		version = s.serverVersion
	} else if match, _ := regexp.MatchString("\\d*\\.\\d*", s.serverVersion); match {
		versionGroupDetailsDownloadURL := s.papermcAPI + "version_group/" + s.serverVersion
		versionGroupDetails, err := getVersionGroupDetails(versionGroupDetailsDownloadURL)
		if err != nil {
			return err
		}
		version = versionGroupDetails.Versions[len(versionGroupDetails.Versions)-1]
	} else if s.serverVersion == "latest" {
		// TODO: Implement latest version resolver
	} else {
		return errors.New("Unsupported version")
	}

	versionDetailsDownloadURL := s.papermcAPI + "versions/" + version
	versionDetails, err := getVersionDetails(versionDetailsDownloadURL)
	if err != nil {
		return err
	}

	buildDetailsURL := versionDetailsDownloadURL + "/builds/" + strconv.Itoa(versionDetails.Builds[len(versionDetails.Builds)-1])
	buildDetails, err := getBuildDetails(versionDetailsDownloadURL + "/builds/" + strconv.Itoa(versionDetails.Builds[len(versionDetails.Builds)-1]))
	if err != nil {
		return err
	}

	if common.CompareHash(s.serverForceUpdate, buildDetails.Downloads.Application.Sha256) {
		logrus.Info("updated server")
		return nil
	}

	if err := utils.DownloadToFile(buildDetailsURL+"/downloads/"+buildDetails.Downloads.Application.Name, "papermc.jar"); err != nil {
		return err
	}

	if err := common.SaveHash(buildDetails.Downloads.Application.Sha256); err != nil {
		return err
	}

	logrus.Info("updated server")
	return nil
}

// StartupCommand retuns the command and arguments used to startup the server.
func (s *Server) StartupCommand() (string, []string) {
	return "java", []string{"-jar", "papermc.jar", "nogui"}
}
