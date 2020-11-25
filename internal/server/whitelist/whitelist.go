package whitelist

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

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
	if _, err := os.Stat("whitelist.json"); !os.IsNotExist(err) {
		logrus.Info("whitelist.json does already exist, skipping configuration")
		return
	}
	logrus.Info("whitelist.json does not exist, configuring it now")

	var whitelist []User
	for _, n := range strings.Split(utils.GetEnvString("WHITELIST", ""), ",") {
		whitelist = append(whitelist, User{
			Name: n,
		})
	}
	b, err := json.Marshal(&whitelist)
	if err != nil {
		logrus.WithError(err).Error("failed to marshall whitelist.json")
		return
	}
	if err := ioutil.WriteFile("whitelist.json", []byte(b), 0); err != nil {
		logrus.WithError(err).Error("failed to write whitelist.json")
	}
}
