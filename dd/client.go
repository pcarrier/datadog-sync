package dd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/meteor/datadog-sync/util"
)

const (
	monitorURL = "https://app.datadoghq.com/api/v1/monitor"
)

var (
	// APIKey should be a Datadog API key from https://app.datadoghq.com/account/settings#api
	APIKey = flag.String("api-key", "", "Datadog API key")
	// AppKey should be a Datadog application key from https://app.datadoghq.com/account/settings#api
	AppKey    = flag.String("app-key", "", "Datadog application key")
	httpDebug = flag.Bool("http-debug", false, "Debug HTTP, useful to understand API failures")
)

func authedURL(base string) string {
	return base + "?api_key=" + *APIKey + "&application_key=" + *AppKey
}

func doHTTP(client *http.Client, method, url, body string) ([]byte, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if *httpDebug {
		msg := fmt.Sprintf("=== HTTP ===\nRequest: %v\nRequest body: %v\nResponse: %v\nResponse body: %v",
			req, body, resp, string(respBody))
		msg = strings.Replace(msg, *APIKey, "HIDDEN_API_KEY", -1)
		msg = strings.Replace(msg, *AppKey, "HIDDEN_APP_KEY", -1)
		logrus.Debug(msg)
	}

	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("status code %d", resp.StatusCode)
	}

	return respBody, nil
}

// GetMonitors downloads all the Datadog monitors from an account through the Datadog API
func GetMonitors(client *http.Client) ([]Monitor, error) {
	resp, err := doHTTP(client, "GET",
		authedURL(monitorURL), "")
	if err != nil {
		return nil, err
	}

	var monitors []Monitor
	if err = json.Unmarshal(resp, &monitors); err != nil {
		return nil, err
	}

	return monitors, nil
}

func (monitor *Monitor) create(client *http.Client) error {
	repr, err := util.Marshal(monitor, util.JSON)
	if err != nil {
		return err
	}
	_, err = doHTTP(client, "POST",
		authedURL(monitorURL), repr)
	return err
}

func (monitor *Monitor) update(client *http.Client, target *Monitor) error {
	msg := target
	msg.ID = nil
	repr, err := util.Marshal(msg, util.JSON)
	if err != nil {
		return err
	}
	_, err = doHTTP(client, "PUT",
		authedURL(fmt.Sprintf("%s/%d", monitorURL, *monitor.ID)), repr)
	return err
}

func (monitor *Monitor) delete(client *http.Client) error {
	_, err := doHTTP(client, "DELETE",
		authedURL(fmt.Sprintf("%s/%d", monitorURL, *monitor.ID)), "")
	return err
}
