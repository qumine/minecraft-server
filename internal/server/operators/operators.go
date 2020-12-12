package operators

import (
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
	if !utils.GetEnvBool("SERVER_OPS_FORCE", false) && utils.FileExists("ops.json") {
		logrus.Info("ops.json already exist, skipping configuration")
		return
	}

	ops := newOps(utils.GetEnvStringList("SERVER_OPS", ""))
	logrus.WithField("ops", ops).Info("ops.json not found, configuring it now")

	if err := utils.WriteFileAsJSON("ops.json", &ops); err != nil {
		logrus.WithError(err).Error("failed to configure ops.json")
	}
	logrus.Info("ops.json configured")
}

func newOps(users []string) []User {
	var ops []User
	for _, n := range users {
		ops = append(ops, newUser(n))
	}
	return ops
}

func newUser(name string) User {
	return User{
		Name:                name,
		Level:               4,
		BypassesPlayerLimit: true,
	}
}
