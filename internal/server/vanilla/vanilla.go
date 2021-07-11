package vanilla

import (
	"github.com/qumine/minecraft-server/internal/server/common"
	"github.com/qumine/minecraft-server/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	serverVersion     = "latest"
	serverForceUpdate = false
	serverVanillaAPI  = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
)

// Server is the struct for vanilla servers.
type Server struct {
	serverVersion     string
	serverForceUpdate bool
	serverVanillaAPI  string
}

// NewVanillaServer creates a new vanilla server.
func NewVanillaServer() *Server {
	return &Server{
		serverVersion:     utils.GetEnvString("SERVER_VERSION", serverVersion),
		serverForceUpdate: utils.GetEnvBool("SERVER_FORCE_UPDATE", serverForceUpdate),
		serverVanillaAPI:  utils.GetEnvString("SERVER_VANILLA_API", serverVanillaAPI),
	}
}

// Configure configures the server.
func (s *Server) Configure() error {
	logrus.WithFields(logrus.Fields{
		"type":        "VANILLA",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"vanillaAPI":  s.serverVanillaAPI,
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
		"type":        "VANILLA",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"vanillaAPI":  s.serverVanillaAPI,
	}).Debug("updating server")

	versionManifest, err := getVersionManifest(s.serverVanillaAPI)
	if err != nil {
		return err
	}

	if s.serverVersion == "latest" {
		s.serverVersion = versionManifest.Latest.Release
	}

	var versionDetailsDownloadURL string
	for i := 0; i < len(versionManifest.Versions); i++ {
		if s.serverVersion == versionManifest.Versions[i].ID {
			versionDetailsDownloadURL = versionManifest.Versions[i].URL
		}
	}
	versionDetails, err := getVersionDetails(versionDetailsDownloadURL)
	if err != nil {
		return err
	}

	if common.CompareHash(s.serverForceUpdate, versionDetails.Downloads.Client.Sha1) {
		logrus.WithFields(logrus.Fields{
			"type":        "VANILLA",
			"version":     s.serverVersion,
			"forceUpdate": s.serverForceUpdate,
			"vanillaAPI":  s.serverVanillaAPI,
		}).Info("updating server skipped, jar seems up to date")
		return nil
	}

	if err := utils.DownloadToFile(versionDetails.Downloads.Server.URL, "vanilla.jar"); err != nil {
		return err
	}

	if err := common.SaveHash(versionDetails.Downloads.Client.Sha1); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"type":        "VANILLA",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"vanillaAPI":  s.serverVanillaAPI,
	}).Info("updated server")
	return nil
}

// UpdatePlugins updates the plugins, if supported.
func (s *Server) UpdatePlugins() error {
	logrus.WithFields(logrus.Fields{
		"type":        "VANILLA",
		"version":     s.serverVersion,
		"forceUpdate": s.serverForceUpdate,
		"vanillaAPI":  s.serverVanillaAPI,
	}).Info("updating plugins failed, not supported by SERVER_TYPE")
	return nil
}

// StartupCommand returns the command and arguments used to startup the server.
func (s *Server) StartupCommand() (string, []string) {
	return "java", common.ServerFlags("vanilla.jar")
}
