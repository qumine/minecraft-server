package vanilla

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type VersionManifest struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []struct {
		ID          string    `json:"id"`
		Type        string    `json:"type"`
		URL         string    `json:"url"`
		Time        time.Time `json:"time"`
		ReleaseTime time.Time `json:"releaseTime"`
	} `json:"versions"`
}

type VersionDetails struct {
	Arguments struct {
		Game []interface{} `json:"game"`
		Jvm  []interface{} `json:"jvm"`
	} `json:"arguments"`
	AssetIndex struct {
		ID        string `json:"id"`
		Sha1      string `json:"sha1"`
		Size      int    `json:"size"`
		TotalSize int    `json:"totalSize"`
		URL       string `json:"url"`
	} `json:"assetIndex"`
	Assets          string `json:"assets"`
	ComplianceLevel int    `json:"complianceLevel"`
	Downloads       struct {
		Client struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"client"`
		ClientMappings struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"client_mappings"`
		Server struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"server"`
		ServerMappings struct {
			Sha1 string `json:"sha1"`
			Size int    `json:"size"`
			URL  string `json:"url"`
		} `json:"server_mappings"`
	} `json:"downloads"`
	ID        string `json:"id"`
	Libraries []struct {
		Downloads struct {
			Artifact struct {
				Path string `json:"path"`
				Sha1 string `json:"sha1"`
				Size int    `json:"size"`
				URL  string `json:"url"`
			} `json:"artifact"`
			Classifiers struct {
				Javadoc struct {
					Path string `json:"path"`
					Sha1 string `json:"sha1"`
					Size int    `json:"size"`
					URL  string `json:"url"`
				} `json:"javadoc"`
				NativesLinux struct {
					Path string `json:"path"`
					Sha1 string `json:"sha1"`
					Size int    `json:"size"`
					URL  string `json:"url"`
				} `json:"natives-linux"`
				NativesMacos struct {
					Path string `json:"path"`
					Sha1 string `json:"sha1"`
					Size int    `json:"size"`
					URL  string `json:"url"`
				} `json:"natives-macos"`
				NativesWindows struct {
					Path string `json:"path"`
					Sha1 string `json:"sha1"`
					Size int    `json:"size"`
					URL  string `json:"url"`
				} `json:"natives-windows"`
				Sources struct {
					Path string `json:"path"`
					Sha1 string `json:"sha1"`
					Size int    `json:"size"`
					URL  string `json:"url"`
				} `json:"sources"`
			} `json:"classifiers"`
		} `json:"downloads,omitempty"`
		Name  string `json:"name"`
		Rules []struct {
			Action string `json:"action"`
			Os     struct {
				Name string `json:"name"`
			} `json:"os"`
		} `json:"rules,omitempty"`
		Natives struct {
			Osx     string `json:"osx,omitempty"`
			Linux   string `json:"linux,omitempty"`
			Windows string `json:"windows,omitempty"`
		} `json:"natives,omitempty"`
		Extract struct {
			Exclude []string `json:"exclude"`
		} `json:"extract,omitempty"`
	} `json:"libraries"`
	Logging struct {
		Client struct {
			Argument string `json:"argument"`
			File     struct {
				ID   string `json:"id"`
				Sha1 string `json:"sha1"`
				Size int    `json:"size"`
				URL  string `json:"url"`
			} `json:"file"`
			Type string `json:"type"`
		} `json:"client"`
	} `json:"logging"`
	MainClass              string    `json:"mainClass"`
	MinimumLauncherVersion int       `json:"minimumLauncherVersion"`
	ReleaseTime            time.Time `json:"releaseTime"`
	Time                   time.Time `json:"time"`
	Type                   string    `json:"type"`
}

func getVersionManifest(versionManifestURL string) (*VersionManifest, error) {
	logrus.WithField("url", versionManifestURL).Debug("downloading versionManifest")
	rsp, getErr := http.Get(versionManifestURL)
	if getErr != nil {
		logrus.WithError(getErr).Error("downloading versionManifest failed")
		return nil, getErr
	}
	logrus.WithField("contentLength", rsp.ContentLength).Trace("downloaded versionManifest")

	logrus.Debug("reading versionManifest")
	if rsp.Body != nil {
		defer rsp.Body.Close()
	}
	body, readErr := ioutil.ReadAll(rsp.Body)
	if readErr != nil {
		logrus.WithError(readErr).Error("reading versionManifest failed")
		return nil, readErr
	}
	logrus.WithField("body", rsp.Body).Trace("read versionManifest")

	logrus.Debug("unmarshalling versionManifest")
	versionManifest := &VersionManifest{}
	jsonErr := json.Unmarshal(body, &versionManifest)
	if jsonErr != nil {
		logrus.WithError(jsonErr).Error("unmarshalling versionManifest failed")
		return nil, jsonErr
	}
	logrus.WithField("versionManifest", versionManifest).Trace("unmarshalled versionManifest")

	return versionManifest, nil
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
