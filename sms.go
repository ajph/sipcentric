package sipcentric

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SmsHistoryItem struct {
	Uri            string  `json:"uri"`
	Created        string  `json:"created"`
	Direction      string  `json:"direction"`
	From           string  `json:"from"`
	To             string  `json:"to"`
	Body           string  `json:"body"`
	SendStatus     string  `json:"sendStatus,omitempty"`
	DeliveryStatus int     `json:"deliveryStatus,omitempty"`
	Cost           float32 `json:"cost"`
}

type SmsHistoryResult struct {
	TotalItems int              `json:"totalItems"`
	PageSize   int              `json:"pageSize"`
	Page       int              `json:"page"`
	Items      []SmsHistoryItem `json:"items"`
}

func (api *Api) SmsHistory(page int, pageSize int) (*SmsHistoryResult, error) {
	resp, err := api.apiRequest("GET", fmt.Sprintf("/customers/me/sms?page=%d&pageSize=%d", page, pageSize), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse result
	switch resp.StatusCode {
	case 200: // OK
		r := &SmsHistoryResult{}
		err := json.NewDecoder(resp.Body).Decode(r)
		if err != nil {
			return nil, err
		}
		return r, err

	default:
		return nil, fmt.Errorf("Invalid response code %d", resp.StatusCode)
	}
}

func (api *Api) SendSms(from string, to string, message string) error {
	// build post data
	post := `
		{
			"type": "smsmessage",
			"from": "` + from + `",
			"to": "` + to + `",
			"body": "` + message + `"
		}
	`

	// make request
	s := strings.NewReader(post)
	resp, err := api.apiRequest("POST", "/customers/me/sms", s)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 201 = OK, other responses undocumented
	switch resp.StatusCode {
	case 201: // CREATED
		return nil

	default:
		err := fmt.Errorf("SMS send failure; %d", resp.Status)
		return err
	}

}
