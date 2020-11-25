package yatopia

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type versionDetails struct {
	branch struct {
		name   string `json:"name"`
		commit struct {
			sha        string `json:"sha"`
			authoredAt string `json:"authoredAt"`
			message    string `json:"message"`
			comment    string `json:"comment"`
		} `json:"commit"`
	} `json:"branch"`
	changeSets []struct {
		sha        string `json:"sha"`
		authoredAt string `json:"authoredAt"`
		message    string `json:"message"`
		comment    string `json:"comment"`
	} `json:"changeSets"`
	number         int    `json:"number"`
	jenkinsViewURL string `json:"jenkinsViewUrl"`
	status         string `json:"status"`
	downloadURL    string `json:"downloadUrl"`
}

func getVersionDetails(versionDetailsURL string) (*versionDetails, error) {
	logrus.WithField("url", versionDetailsURL).Debug("downloading versionDetails")
	rsp, getErr := http.Get(versionDetailsURL)
	if getErr != nil {
		logrus.WithError(getErr).Error("downloading versionDetails failed")
		return nil, getErr
	}
	logrus.WithField("contentLength", rsp.ContentLength).Trace("downloaded versionDetails")

	logrus.Debug("reading versionDetails")
	if rsp.Body != nil {
		defer rsp.Body.Close()
	}
	body, readErr := ioutil.ReadAll(rsp.Body)
	if readErr != nil {
		logrus.WithError(readErr).Error("reading versionDetails failed")
		return nil, readErr
	}
	logrus.WithField("body", rsp.Body).Trace("read versionDetails")

	logrus.Debug("unmarshalling versionDetails")
	versionDetails := &versionDetails{}
	jsonErr := json.Unmarshal(body, versionDetails)
	if jsonErr != nil {
		logrus.WithError(jsonErr).Error("unmarshalling versionDetails failed")
		return nil, jsonErr
	}
	logrus.WithField("versionDetails", versionDetails).Trace("unmarshalled versionDetails")
	return versionDetails, nil
}
