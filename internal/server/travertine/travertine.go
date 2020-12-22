package travertine

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/qumine/qumine-server-java/internal/server/common"
	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	serverVersion       = "latest"
	serverForceUpdate   = false
	serverTravertineAPI = "https://papermc.io/api/v2/projects/travertine/"
)

// Server is the struct for travertine servers.
type Server struct {
	serverVersion     string
	serverForceUpdate bool
	travertineAPI     string
}

// NewTravertineServer creates a new travertine server.
func NewTravertineServer() *Server {
	return &Server{
		serverVersion:     utils.GetEnvString("SERVER_VERSION", serverVersion),
		serverForceUpdate: utils.GetEnvBool("SERVER_FORCE_UPDATE", serverForceUpdate),
		travertineAPI:     utils.GetEnvString("SERVER_TRAVERTINE_API", serverTravertineAPI),
	}
}

// Configure configures the server.
func (s *Server) Configure() error {
	return nil
}

// Update updates the resource, if supported uses cache.
func (s *Server) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":          "TRAVERTINE",
		"version":       s.serverVersion,
		"forceUpdate":   s.serverForceUpdate,
		"travertineAPI": s.travertineAPI,
	}).Debug("updating server")

	version := ""
	if match, _ := regexp.MatchString("\\d*\\.\\d*\\.\\d", s.serverVersion); match {
		version = s.serverVersion
	} else if match, _ := regexp.MatchString("\\d*\\.\\d*", s.serverVersion); match {
		versionGroupDetailsDownloadURL := s.travertineAPI + "version_group/" + s.serverVersion
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

	versionDetailsDownloadURL := s.travertineAPI + "versions/" + version
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
		logrus.WithFields(logrus.Fields{
			"type":          "TRAVERTINE",
			"version":       s.serverVersion,
			"forceUpdate":   s.serverForceUpdate,
			"travertineAPI": s.travertineAPI,
		}).Info("updating server skipped, jar seems up to date")
		return nil
	}

	if err := utils.DownloadToFile(buildDetailsURL+"/downloads/"+buildDetails.Downloads.Application.Name, "travertine.jar"); err != nil {
		return err
	}

	if err := common.SaveHash(buildDetails.Downloads.Application.Sha256); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"type":          "TRAVERTINE",
		"version":       s.serverVersion,
		"forceUpdate":   s.serverForceUpdate,
		"travertineAPI": s.travertineAPI,
	}).Info("updated server")
	return nil
}

// StartupCommand retuns the command and arguments used to startup the server.
func (s *Server) StartupCommand() (string, []string) {
	return "java", []string{"-jar", "travertine.jar", "nogui"}
}
