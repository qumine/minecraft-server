package server

import (
	"fmt"
	"strings"

	"github.com/qumine/qumine-server-java/internal/server/updater/vanilla"
	"github.com/qumine/qumine-server-java/internal/server/updater/yatopia"
	"github.com/qumine/qumine-server-java/internal/utils"
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
func NewUpdater() (Updater, error) {
	switch strings.ToUpper(utils.GetEnvString("SERVER_TYPE", "VANILLA")) {
	case "VANILLA":
		return vanilla.NewVanillaUpdater(), nil
	case "YATOPIA":
		return yatopia.NewYatopiaUpdater(), nil
	default:
		return nil, fmt.Errorf("serverType(%s) not supported", strings.ToUpper(utils.GetEnvString("SERVER_TYPE", "VANILLA")))
	}
}
