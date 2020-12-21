package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/qumine/qumine-server-java/internal/server/custom"
	"github.com/qumine/qumine-server-java/internal/server/papermc"
	"github.com/qumine/qumine-server-java/internal/server/starter"
	"github.com/qumine/qumine-server-java/internal/server/travertine"
	"github.com/qumine/qumine-server-java/internal/server/vanilla"
	"github.com/qumine/qumine-server-java/internal/server/waterfall"
	"github.com/qumine/qumine-server-java/internal/server/yatopia"
)

// Server is the interface for servers.
type Server interface {
	// Configure configures the server.
	Configure() error
	// Update updates the resource, if supported uses cache.
	Update() error
	// StartupCommand retuns the command and arguments used to startup the server.
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
	case "YATOPIA":
		return yatopia.NewYatopiaServer(), nil
	default:
		return nil, fmt.Errorf("serverType(%s) not supported", strings.ToUpper(os.Getenv("SERVER_TYPE")))
	}
}
