package common

import (
	"os"
	"strings"

	props "github.com/magiconair/properties"
	"github.com/sirupsen/logrus"
)

const serverPropertiesPath = "server.properties"

// ConfigureServerProperties configures the server.properties
func ConfigureServerProperties() error {
	pp, err := props.LoadFile("", props.UTF8)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Infof("%s does not exist yet, will create it", serverPropertiesPath)
		} else {
			logrus.WithError(err).Error("Error loading %s", serverPropertiesPath)
		}
		pp = props.NewProperties()
	}
	p := pp.Map()
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		if strings.HasPrefix(variable[0], "SERVER_PROPERTIES_") {
			variable[0] = strings.ReplaceAll(variable[0], "SERVER_PROPERTIES_", "")
			variable[0] = strings.ReplaceAll(variable[0], "_", "-")
			variable[0] = strings.ToLower(variable[0])
			p[variable[0]] = variable[1]
		}
	}

	f, err := os.Create(serverPropertiesPath)
	if err != nil {
		return err
	}
	w, err := props.LoadMap(p).Write(f, props.UTF8)
	if err != nil {
		return err
	}
	logrus.WithField("bytesWritte", w).Debugf("%s configured", serverPropertiesPath)
	return nil
}
