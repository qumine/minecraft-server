package common

import (
	"os"
	"strings"

	"github.com/qumine/minecraft-server/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	pluginsPath = "plugins/"
)

// UpdatePlugins updates the plugins provided with SERVER_PLUGINS, currently only supports bukkit based plugins.
func UpdatePlugins(plugins []string) {
	os.MkdirAll(pluginsPath, 0777)
	for _, plugin := range plugins {
		if strings.HasPrefix(plugin, "http://") || strings.HasPrefix(plugin, "https://") {
			updatePluginFromURL(plugin)
		}
	}
}

func updatePluginFromURL(url string) {
	urlParts := strings.Split(url, "/")
	plugin := urlParts[len(urlParts)-1]
	logrus.WithFields(logrus.Fields{
		"url":    url,
		"plugin": plugin,
	}).Trace("updating plugin")
	if err := utils.DownloadToFile(url, pluginsPath+plugin); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"url":    url,
			"plugin": plugin,
		}).Warn("updating plugin failed")
	}
	logrus.WithFields(logrus.Fields{
		"url":    url,
		"plugin": plugin,
	}).Debug("updated plugin")
}
