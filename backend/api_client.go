package backend

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

// APIClient represents an API client for interacting with the YNAB API, embedding a Resty client to handle HTTP requests
type APIClient struct {
	*resty.Client
}

// AccessToken is the YNAB Personal Access Token
// Please ensure this token is kept secure and not exposed publicly
const AccessToken string = ""

// Configure sets up the APIClient with the necessary configurations for interacting with the YNAB API
func (client *APIClient) Configure() {
	client.SetBaseURL("https://api.ynab.com/v1")
	client.SetHeader("Accept", "application/json")
	client.SetAuthToken(AccessToken)
}

// ValidateResponse checks if the API response indicates an error
func (client *APIClient) ValidateResponse(response *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if response.IsError() {
		return errors.New(string(response.Body()))
	}

	return nil
}
