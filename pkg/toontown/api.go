package toontown

import (
	"fmt"
	"runtime"
)

type ToontownClient struct {
	client *HTTPClient
}

func New() *ToontownClient {
	baseURL := "https://toontownrewritten.com/api"
	userAgent := fmt.Sprintf("MadamChuckle (https://github.com/jaczerob/madamchuckle) %s", runtime.Version())

	return &ToontownClient{
		client: NewHTTPClient(baseURL, userAgent),
	}
}

func (c *ToontownClient) Invasions() (invasionData *InvasionData, err error) {
	invasionData = new(InvasionData)
	err = c.client.Get("/invasions", invasionData)
	return
}

func (c *ToontownClient) Population() (populationData *PopulationData, err error) {
	populationData = new(PopulationData)
	err = c.client.Get("/population", populationData)
	return
}

func (c *ToontownClient) FieldOffices() (fieldOfficeData *FieldOfficeData, err error) {
	fieldOfficeData = new(FieldOfficeData)
	err = c.client.Get("/fieldoffices", fieldOfficeData)
	return
}

func (c *ToontownClient) Status() (status *Status, err error) {
	status = new(Status)
	err = c.client.Get("/status", status)
	return
}

func (c *ToontownClient) SillyMeter() (sillymeter *SillyMeter, err error) {
	sillymeter = new(SillyMeter)
	err = c.client.Get("/sillymeter", sillymeter)
	return
}
