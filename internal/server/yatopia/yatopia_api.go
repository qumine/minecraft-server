package yatopia

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

// VersionDetails represents the version details of the yatopia api
type VersionDetails struct {
	// Branch represents the branch of the version details
	Branch struct {
		// Name is the name of the branch
		Name string `json:"name"`
		// Commit is the latest commit of the branch
		Commit struct {
			// Sha is the commit hash of the commit
			Sha string `json:"sha"`
			// AuthoredAt is the commit date of the commit
			AuthoredAt string `json:"authoredAt"`
			// Message is the commit message of the commit
			Message string `json:"message"`
			// Comment is the comment of the commit
			Comment string `json:"comment"`
		} `json:"commit"`
	} `json:"branch"`
	ChangeSets []struct {
		// Sha is the commit hash of the commit
		Sha string `json:"sha"`
		// AuthoredAt is the commit date of the commit
		AuthoredAt string `json:"authoredAt"`
		// Message is the commit message of the commit
		Message string `json:"message"`
		// Comment is the comment of the commit
		Comment string `json:"comment"`
	} `json:"changeSets"`
	// Number is the number of the build
	Number int `json:"number"`
	// JenkinsViewURL is the URL to the jenkins view of the build
	JenkinsViewURL string `json:"jenkinsViewUrl"`
	// Status is the build status
	Status string `json:"status"`
	// DownloadURL is the URL of the jar
	DownloadURL string `json:"downloadUrl"`
}

func getVersionDetails(versionDetailsURL string) (*VersionDetails, error) {
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
	versionDetails := &VersionDetails{}
	jsonErr := json.Unmarshal(body, &versionDetails)
	if jsonErr != nil {
		logrus.WithError(jsonErr).Error("unmarshalling versionDetails failed")
		return nil, jsonErr
	}
	logrus.WithField("versionDetails", versionDetails).Trace("unmarshalled versionDetails")
	return versionDetails, nil
}
