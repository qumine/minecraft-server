package server

import (
	"fmt"
	"strings"

	"github.com/qumine/qumine-server-java/internal/updater/server/vanilla"
	"github.com/qumine/qumine-server-java/internal/updater/server/yatopia"
)

const (
	vanillaServerDownloadAPI = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
)

// Updater is the interface for updaters of servers.
type Updater interface {
	// Update updates the resource, if supported uses cache.
	Update() error
}

// NewUpdater creates a new updater for the provided serverType.
func NewUpdater(serverType string, serverVersion string, serverDownloadAPI string, serverForceDownload bool) (Updater, error) {
	switch strings.ToUpper(serverType) {
	case "VANILLA":
		return vanilla.NewVanillaUpdater(serverVersion, serverDownloadAPI, serverForceDownload), nil
	case "YATOPIA":
		return yatopia.NewYatopiaUpdater(serverVersion, serverDownloadAPI, serverForceDownload), nil
	default:
		return nil, fmt.Errorf("serverType(%s) not supported", serverType)
	}
}
