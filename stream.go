package sipcentric

import (
	"encoding/json"
	"net"
	"net/http"
	"runtime"
	"time"
)

type StreamEvent struct {
	Event    string                 `json:"event"`
	Location string                 `json:"location"`
	Values   map[string]interface{} `json:"values"`
}

func (api *Api) Stream() (<-chan *StreamEvent, error) {

	// setup dialer
	var conn net.Conn
	dialer := func(netw, addr string) (net.Conn, error) {
		netc, err := net.DialTimeout(netw, addr, 5*time.Second)
		if err != nil {
			return nil, err
		}
		conn = netc
		return netc, nil
	}

	// setup request
	url := API_URL + "/stream"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(api.Username, api.Password)

	// make request
	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan *StreamEvent)

	go func() {
		for {
			conn.SetReadDeadline(time.Now().Add(90 * time.Second))
			e := &StreamEvent{}
			err := json.NewDecoder(resp.Body).Decode(e)
			if err == nil {
				ch <- e

			} else if n, ok := err.(net.Error); ok && n.Timeout() {
				resp.Body.Close()
				close(ch)
				return

			} else {
				// ignore

			}
			runtime.Gosched()
		}
	}()

	return ch, nil

}
