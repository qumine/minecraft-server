package server

const (
	vanillaServerDownloadAPI = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
)

// Updater is the interface for updaters of servers.
type Updater interface {
	// Update updates the resource, if supported uses cache.
	Update() error
}
