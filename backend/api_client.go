package backend

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

type APIClient struct {
	*resty.Client
}

func (client *APIClient) Configure() {
	client.SetBaseURL("https://api.ynab.com/v1")
	client.SetHeader("Accept", "application/json")
	client.SetAuthToken("") // Set YNAB access token
}

func (client *APIClient) ValidateResponse(response *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if response.IsError() {
		return errors.New(string(response.Body()))
	}

	return nil
}
