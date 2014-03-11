package sipcentric

// https://apigee.com/sipcentric/embed/console/beta

import (
	"fmt"
	"io"
	"net/http"
)

const (
	API_URL = "http://pbx.sipcentric.com/api/v1"
)

type Api struct {
	Username string
	Password string
}

func (api *Api) apiRequest(method string, path string, iord io.Reader) (*http.Response, error) {
	// set up a HEAD request to validate auth
	url := API_URL + path
	req, err := http.NewRequest(method, url, iord)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(api.Username, api.Password)
	req.Header.Add("Content-Type", "application/json")

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (api *Api) ValidateLogin() error {
	resp, err := api.apiRequest("HEAD", "/customers/me", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 200 = OK, other responses undocumented
	// so assume != 200 is a failure
	switch resp.StatusCode {
	case 200: // OK
		return nil

	default:
		err := fmt.Errorf("Login failure; %s", resp.Status)
		return err
	}
}

func New(username string, password string) (*Api, error) {
	// new api
	api := &Api{
		Username: username,
		Password: password,
	}

	// validate login credentials
	if err := api.ValidateLogin(); err != nil {
		return nil, err
	}

	// success
	return api, nil
}