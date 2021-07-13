package utils

import (
	"strings"

	"github.com/artdarek/go-unzip/pkg/unzip"
	"github.com/sirupsen/logrus"
)

// Update updates the resource, if supported uses cache.
func DownloadAdditionalFiles() error {
	additionalFiles := GetEnvStringList("ADDITIONAL_FILES", "")
	logrus.WithFields(logrus.Fields{
		"additionalFiles": additionalFiles,
	}).Debug("downloading additonal files")

	if len(additionalFiles) == 0 {
		logrus.WithFields(logrus.Fields{
			"additionalFiles": additionalFiles,
		}).Info("downloading additonal files skipped")
		return nil
	}

	for _, additionalFile := range additionalFiles {
		parts := strings.Split(additionalFile, "/")
		filename := parts[len(parts)-1]
		if err := DownloadToFile(additionalFile, filename); err != nil {
			return err
		}

		if _, err := unzip.New().Extract(filename, "/data"); err != nil {
			return err
		}
	}

	logrus.WithFields(logrus.Fields{
		"additionalFiles": additionalFiles,
	}).Info("downloaded additonal files")
	return nil
}
