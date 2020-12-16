package common

import (
	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const opsPath = "ops.json"

// Operator is a user on the operator list
type Operator struct {
	// UUID is the uuid of the user
	UUID string `json:"uuid"`
	// Name is the name of the user
	Name string `json:"name"`
	// Level is the permission level of the user
	Level int `json:"level"`
	// BypassesPlayerLimit defines if the user can bypass the server player limit
	BypassesPlayerLimit bool `json:"bypassesPlayerLimit"`
}

// ConfigureOps configures the ops.json
func ConfigureOps() error {
	if !utils.GetEnvBool("SERVER_OPS_OVERRIDE", false) && utils.FileExists(opsPath) {
		logrus.Info("ops.json already exist, skipping configuration")
		return nil
	}

	ops := newOps(utils.GetEnvStringList("SERVER_OPS", ""))
	logrus.WithField("ops", ops).Info("ops.json not found, configuring it now")

	if err := utils.WriteFileAsJSON(opsPath, &ops); err != nil {
		return err
	}
	logrus.Debug("ops.json configured")
	return nil
}

func newOps(users []string) []Operator {
	var ops []Operator
	for _, n := range users {
		if n == "" {
			continue
		}
		ops = append(ops, newOperator(n))
	}
	return ops
}

func newOperator(name string) Operator {
	return Operator{
		Name:                name,
		Level:               4,
		BypassesPlayerLimit: true,
	}
}
