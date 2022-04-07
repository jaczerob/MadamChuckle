package httpclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type HTTPClient struct {
	BaseURL   string
	UserAgent string

	httpClient *http.Client
}

func NewHTTPClient(baseURL string, userAgent string) *HTTPClient {
	return &HTTPClient{
		BaseURL:   baseURL,
		UserAgent: userAgent,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *HTTPClient) makeRequest(method string, urlPath string, parameters map[string]string, contentType string) (req *http.Request, err error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, urlPath)
	req, err = http.NewRequest(method, url, nil)

	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Content-Type", contentType)

	if parameters != nil {
		q := req.URL.Query()
		for key, element := range parameters {
			q.Add(key, element)
		}
		req.URL.RawQuery = q.Encode()
	}

	return
}

func (c *HTTPClient) Request(method string, urlPath string, parameters map[string]string, contentType string, v interface{}) (err error) {
	req, err := c.makeRequest(method, urlPath, parameters, contentType)
	if err != nil {
		return
	}

	log.WithFields(log.Fields{
		"method": method,
		"url":    req.URL,
	}).Info("attempting request")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(body, &v)
}
