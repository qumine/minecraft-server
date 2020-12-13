package whitelist

import (
	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

// User is a user on the whitelist
type User struct {
	// UUID is the uuid of the user
	UUID string `json:"uuid"`
	// Name is the name of the user
	Name string `json:"name"`
}

// Configure the whitelist.json
func Configure() {
	if !utils.GetEnvBool("SERVER_WHITE_LIST_OVERRIDE", false) && utils.FileExists("whitelist.json") {
		logrus.Info("whitelist.json already exist, skipping configuration")
		return
	}

	whitelist := newWhitelist(utils.GetEnvStringList("SERVER_WHITE_LIST", ""))
	logrus.WithField("whitelist", whitelist).Info("whitelist.json not found, configuring it now")

	if err := utils.WriteFileAsJSON("whitelist.json", &whitelist); err != nil {
		logrus.WithError(err).Error("failed to configure whitelist.json")
	}
	logrus.Info("whitelist.json configured")
}

func newWhitelist(users []string) []User {
	var whitelist []User
	for _, n := range users {
		whitelist = append(whitelist, newUser(n))
	}
	return whitelist
}

func newUser(name string) User {
	return User{
		Name: name,
	}
}
