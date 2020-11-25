package operators

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

// User is a user on the operator list
type User struct {
	// UUID is the uuid of the user
	UUID string `json:"uuid"`
	// Name is the name of the user
	Name string `json:"name"`
	// Level is the permission level of the user
	Level int `json:"level"`
	// BypassesPlayerLimit defines if the user can bypass the server player limit
	BypassesPlayerLimit bool `json:"bypassesPlayerLimit"`
}

// Configure the ops.json
func Configure() {
	if _, err := os.Stat("ops.json"); !os.IsNotExist(err) {
		logrus.Info("ops.json does already exist, skipping configuration")
		return
	}
	logrus.Info("ops.json does not exist, configuring it now")

	var whitelist []User
	for _, n := range strings.Split(utils.GetEnvString("OPS", ""), ",") {
		whitelist = append(whitelist, User{
			Name:                n,
			Level:               4,
			BypassesPlayerLimit: true,
		})
	}
	b, err := json.Marshal(&whitelist)
	if err != nil {
		logrus.WithError(err).Error("failed to marshall ops.json")
		return
	}
	if err := ioutil.WriteFile("ops.json", []byte(b), 0); err != nil {
		logrus.WithError(err).Error("failed to write ops.json")
	}
}
