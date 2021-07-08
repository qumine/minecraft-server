package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/qumine/minecraft-server/internal/server/custom"
	"github.com/qumine/minecraft-server/internal/server/papermc"
	"github.com/qumine/minecraft-server/internal/server/starter"
	"github.com/qumine/minecraft-server/internal/server/travertine"
	"github.com/qumine/minecraft-server/internal/server/vanilla"
	"github.com/qumine/minecraft-server/internal/server/waterfall"
)

// Server is the interface for servers.
type Server interface {
	// Configure configures the server.
	Configure() error
	// Update updates the server, if supported uses cache.
	Update() error
	// UpdatePlugins updates the plugins, if supported.
	UpdatePlugins() error
	// StartupCommand returns the command and arguments used to startup the server.
	StartupCommand() (string, []string)
}

// NewServer creates a new server for the serverType provided in the env variable SERVER_TYPE.
func NewServer() (Server, error) {
	switch strings.ToUpper(os.Getenv("SERVER_TYPE")) {
	case "CUSTOM":
		return custom.NewCustomServer(), nil
	case "PAPERMC":
		return papermc.NewPaperMCServer(), nil
	case "STARTER":
		return starter.NewStarterServer(), nil
	case "TRAVERTINE":
		return travertine.NewTravertineServer(), nil
	case "VANILLA":
		return vanilla.NewVanillaServer(), nil
	case "WATERFALL":
		return waterfall.NewWaterfallServer(), nil
	default:
		return nil, fmt.Errorf("serverType(%s) not supported", strings.ToUpper(os.Getenv("SERVER_TYPE")))
	}
}
