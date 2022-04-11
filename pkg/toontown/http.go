package toontown

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	BaseURL    string
	UserAgent  string
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

func (c *HTTPClient) makeRequest(path string) (req *http.Request, err error) {
	url := strings.Join([]string{c.BaseURL, path}, "")
	req, err = http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("Content-Type", "application/json")
	return
}

func (c *HTTPClient) Get(path string, b ToontownObject) (err error) {
	req, err := c.makeRequest(path)
	if err != nil {
		return
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(body, &b)
}
