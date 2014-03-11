# sipcentric

A client implementation for Sipcentric's v1 API in Go

## example

```go
func main() {
	var err error

	// initialise api
	api, err := sipcentric.New("user", "pass")
	if err != nil {
		fmt.Printf("Error initialising API: %s\n", err)
		return
	}

	// check credit
	c, err := api.Credit()
	if err != nil {
		fmt.Printf("Error getting credit: %s\n", err)
		return
	}
	fmt.Printf("%+v\n", c)

	// send SMS
	err = api.SendSms("12345678910", "12345678910", "test!")
	if err != nil {
		fmt.Printf("Error sending sms: %s\n", err)
		return
	}

	// SMS history
	s, err := api.SmsHistory(1, 10)
	if err != nil {
		fmt.Printf("Error retrieving sms history: %s\n", err)
		return
	}
	fmt.Printf("%+v\n", s)

	// consume event stream
	for {
		ch, err := api.Stream()
		if err != nil {
			fmt.Printf("Error starting stream: %s\n", err)
			return
		}

		// process events as they come in
		for ev := range ch {
			switch ev.Event {
			case "heartbeat":
				fmt.Println("got heartbeat!")
			default:
				fmt.Printf("%+v\n", ev)
			}
		}

		fmt.Println("Disconnected from stream, retrying...")
		time.Sleep(5 * time.Second)
	}

}
```

## license
Do what you want with it