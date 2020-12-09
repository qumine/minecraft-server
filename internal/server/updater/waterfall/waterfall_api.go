package waterfall

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// VersionManifest represents the version manifest of the waterfall api
type VersionManifest struct {
	ProjectID     string   `json:"project_id"`
	ProjectName   string   `json:"project_name"`
	VersionGroups []string `json:"version_groups"`
	Versions      []string `json:"versions"`
}

// VersionGroupDetails represents the versiongroup details of the waterfall api
type VersionGroupDetails struct {
	ProjectID    string   `json:"project_id"`
	ProjectName  string   `json:"project_name"`
	VersionGroup string   `json:"version_group"`
	Versions     []string `json:"versions"`
	Builds       []struct {
		Build   int       `json:"build"`
		Time    time.Time `json:"time"`
		Changes []struct {
			Commit  string `json:"commit"`
			Summary string `json:"summary"`
			Message string `json:"message"`
		} `json:"changes"`
		Downloads struct {
			Application struct {
				Name   string `json:"name"`
				Sha256 string `json:"sha256"`
			} `json:"application"`
		} `json:"downloads"`
	} `json:"builds"`
}

// VersionDetails represents the version details of the waterfall api
type VersionDetails struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Builds      []int  `json:"builds"`
}

// BuildDetails represents the build details of the waterfall api
type BuildDetails struct {
	ProjectID   string    `json:"project_id"`
	ProjectName string    `json:"project_name"`
	Version     string    `json:"version"`
	Build       int       `json:"build"`
	Time        time.Time `json:"time"`
	Changes     []struct {
		Commit  string `json:"commit"`
		Summary string `json:"summary"`
		Message string `json:"message"`
	} `json:"changes"`
	Downloads struct {
		Application struct {
			Name   string `json:"name"`
			Sha256 string `json:"sha256"`
		} `json:"application"`
	} `json:"downloads"`
}

func getVersionGroupDetails(versionGroupDetailsURL string) (*VersionGroupDetails, error) {
	logrus.WithField("url", versionGroupDetailsURL).Debug("downloading versionGroupDetails")
	rsp, getErr := http.Get(versionGroupDetailsURL)
	if getErr != nil {
		logrus.WithError(getErr).Error("downloading versionGroupDetails failed")
		return nil, getErr
	}
	logrus.WithField("contentLength", rsp.ContentLength).Trace("downloaded versionGroupDetails")

	logrus.Debug("reading versionGroupDetails")
	if rsp.Body != nil {
		defer rsp.Body.Close()
	}
	body, readErr := ioutil.ReadAll(rsp.Body)
	if readErr != nil {
		logrus.WithError(readErr).Error("reading versionGroupDetails failed")
		return nil, readErr
	}
	logrus.WithField("body", rsp.Body).Trace("read versionGroupDetails")

	logrus.Debug("unmarshalling versionGroupDetails")
	versionGroupDetails := &VersionGroupDetails{}
	jsonErr := json.Unmarshal(body, &versionGroupDetails)
	if jsonErr != nil {
		logrus.WithError(jsonErr).Error("unmarshalling versionGroupDetails failed")
		return nil, jsonErr
	}
	logrus.WithField("versionGroupDetails", versionGroupDetails).Trace("unmarshalled versionGroupDetails")
	return versionGroupDetails, nil
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

func getBuildDetails(buildDetailsURL string) (*BuildDetails, error) {
	logrus.WithField("url", buildDetailsURL).Debug("downloading buildDetails")
	rsp, getErr := http.Get(buildDetailsURL)
	if getErr != nil {
		logrus.WithError(getErr).Error("downloading buildDetails failed")
		return nil, getErr
	}
	logrus.WithField("contentLength", rsp.ContentLength).Trace("downloaded buildDetails")

	logrus.Debug("reading buildDetails")
	if rsp.Body != nil {
		defer rsp.Body.Close()
	}
	body, readErr := ioutil.ReadAll(rsp.Body)
	if readErr != nil {
		logrus.WithError(readErr).Error("reading buildDetails failed")
		return nil, readErr
	}
	logrus.WithField("body", rsp.Body).Trace("read buildDetails")

	logrus.Debug("unmarshalling buildDetails")
	buildDetails := &BuildDetails{}
	jsonErr := json.Unmarshal(body, &buildDetails)
	if jsonErr != nil {
		logrus.WithError(jsonErr).Error("unmarshalling buildDetails failed")
		return nil, jsonErr
	}
	logrus.WithField("buildDetails", buildDetails).Trace("unmarshalled buildDetails")
	return buildDetails, nil
}
