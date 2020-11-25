package properties

import (
	"os"
	"strings"

	props "github.com/magiconair/properties"
	"github.com/sirupsen/logrus"
)

// Configure the server.properties
func Configure() {
	pp, err := props.LoadFile("", props.UTF8)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Info("server.properties does not exist yet, will create it")
		} else {
			logrus.WithError(err).Error("Error loading server.properties")
		}
		pp = props.NewProperties()
	}
	p := pp.Map()
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		if strings.HasPrefix(variable[0], "PROPERTIES_") {
			variable[0] = strings.ReplaceAll(variable[0], "PROPERTIES_", "")
			variable[0] = strings.ReplaceAll(variable[0], "_", "-")
			variable[0] = strings.ToLower(variable[0])
			p[variable[0]] = variable[1]
		}
	}

	f, err := os.Create("server.properties")
	if err != nil {
		logrus.WithError(err).Error("Failed to create server.properties")
	}
	w, err := props.LoadMap(p).Write(f, props.UTF8)
	if err != nil {
		logrus.WithError(err).Error("Failed to write server.properties")
	}
	logrus.WithField("bytesWritte", w).Debug("written server.properties")
}
