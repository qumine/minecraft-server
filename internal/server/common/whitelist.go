package common

import (
	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const whitelistPath = "whitelist.json"

// User is a user on the whitelist
type User struct {
	// UUID is the uuid of the user
	UUID string `json:"uuid"`
	// Name is the name of the user
	Name string `json:"name"`
}

// ConfigureWhitelist configures the whitelist.json
func ConfigureWhitelist() error {
	if !utils.GetEnvBool("SERVER_WHITE_LIST_OVERRIDE", false) && utils.FileExists(whitelistPath) {
		logrus.Info("whitelist.json already exist, skipping configuration")
		return nil
	}

	whitelist := newWhitelist(utils.GetEnvStringList("SERVER_WHITE_LIST", ""))
	logrus.WithField("whitelist", whitelist).Infof("%s not found, configuring it now", whitelistPath)

	if err := utils.WriteFileAsJSON(whitelistPath, &whitelist); err != nil {
		return err
	}
	logrus.Debugf("%s configured", whitelistPath)
	return nil
}

func newWhitelist(users []string) []User {
	var whitelist []User
	for _, n := range users {
		if n == "" {
			continue
		}
		whitelist = append(whitelist, newUser(n))
	}
	return whitelist
}

func newUser(name string) User {
	return User{
		Name: name,
	}
}
