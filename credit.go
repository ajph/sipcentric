package sipcentric

import (
	"encoding/json"
	"fmt"
)

type CreditStatus struct {
	AccountType     string  `json:"accountType"`
	CreditRemaining float32 `json:"creditRemaining"`
}

func (api *Api) Credit() (*CreditStatus, error) {
	resp, err := api.apiRequest("GET", "/customers/me/creditstatus", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse result
	switch resp.StatusCode {
	case 200: // OK
		r := &CreditStatus{}
		err := json.NewDecoder(resp.Body).Decode(r)
		if err != nil {
			return nil, err
		}
		return r, err

	default:
		return nil, fmt.Errorf("Invalid response code %d", resp.StatusCode)
	}
}
