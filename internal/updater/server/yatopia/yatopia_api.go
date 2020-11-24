package yatopia

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type VersionDetails struct {
	Branch struct {
		Name   string `json:"name"`
		Commit struct {
			Sha        string `json:"sha"`
			AuthoredAt string `json:"authoredAt"`
			Message    string `json:"message"`
			Comment    string `json:"comment"`
		} `json:"commit"`
	} `json:"branch"`
	ChangeSets []struct {
		Sha        string `json:"sha"`
		AuthoredAt string `json:"authoredAt"`
		Message    string `json:"message"`
		Comment    string `json:"comment"`
	} `json:"changeSets"`
	Number         int    `json:"number"`
	JenkinsViewURL string `json:"jenkinsViewUrl"`
	Status         string `json:"status"`
	DownloadURL    string `json:"downloadUrl"`
}

func getVersionDetails(versionDetailsURL string) (*VersionDetails, error) {
	logrus.WithField("versionDetailsURL", versionDetailsURL).Debug("downloading versionDetails")
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
	versionDetails := &VersionDetails{}
	jsonErr := json.Unmarshal(body, &versionDetails)
	if jsonErr != nil {
		logrus.WithError(jsonErr).Error("unmarshalling versionDetails failed")
		return nil, jsonErr
	}
	logrus.WithField("versionDetails", versionDetails).Trace("unmarshalled versionDetails")
	return versionDetails, nil
}
