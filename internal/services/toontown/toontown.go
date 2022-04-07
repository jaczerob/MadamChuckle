package toontown

import (
	"fmt"
	"runtime"

	"github.com/jaczerob/madamchuckle/pkg/httpclient"
)

type ToontownClient struct {
	client *httpclient.HTTPClient
}

func New() *ToontownClient {
	baseURL := "https://toontownrewritten.com/api"
	userAgent := fmt.Sprintf("MadamChuckle (https://github.com/jaczerob/madamchuckle) %s", runtime.Version())

	return &ToontownClient{
		client: httpclient.NewHTTPClient(baseURL, userAgent),
	}
}

func (c *ToontownClient) Invasions() (invasionData *InvasionData, err error) {
	invasionData = new(InvasionData)
	err = c.client.Request("GET", "/invasions", nil, "application/json", invasionData)
	return
}

func (c *ToontownClient) Population() (populationData *PopulationData, err error) {
	populationData = new(PopulationData)
	err = c.client.Request("GET", "/population", nil, "application/json", populationData)
	return
}

func (c *ToontownClient) FieldOffices() (fieldOfficeData *FieldOfficeData, err error) {
	fieldOfficeData = new(FieldOfficeData)
	err = c.client.Request("GET", "/population", nil, "application/json", fieldOfficeData)
	return
}
