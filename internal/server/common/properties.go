package common

import (
	"os"
	"strings"

	props "github.com/magiconair/properties"
	"github.com/sirupsen/logrus"
)

const (
	serverPropertiesPath = "server.properties"

	envServerProperties = "SERVER_PROPERTIES_"
)

// ConfigureServerProperties configures the server.properties
func ConfigureServerProperties() error {
	logrus.Debugf("configuring %s", serverPropertiesPath)

	pp, err := props.LoadFile("", props.UTF8)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		pp = props.NewProperties()
	}
	p := pp.Map()
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		if strings.HasPrefix(variable[0], envServerProperties) {
			variable[0] = strings.ReplaceAll(variable[0], envServerProperties, "")
			variable[0] = strings.ReplaceAll(variable[0], "_", "-")
			variable[0] = strings.ToLower(variable[0])
			p[variable[0]] = variable[1]
		}
	}

	f, err := os.Create(serverPropertiesPath)
	if err != nil {
		return err
	}
	if _, err := props.LoadMap(p).Write(f, props.UTF8); err != nil {
		return err
	}

	logrus.Infof("configured %s", serverPropertiesPath)
	return nil
}
