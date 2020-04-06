# Goflight
![](https://goreportcard.com/badge/github.com/marcelblijleven/goflight)

Go client for the [OpenSky REST API](http://opensky-network.org).

## Example
Get all states

```go
package main

import (
	"fmt"
	"github.com/marcelblijleven/goflight"
	"net/http"
	"time"
)

func main() {
	httpClient := http.Client{Timeout: time.Second * 30}
	client, _ := goflight.NewClient(
		"user@email.com",
		"tops3cr3t",
		&httpClient,
        nil,
	)

	timeParam := time.Now()
	icao24 := "3c6444"

	resp, _ := client.States.GetAllStates(timeParam, icao24)
    fmt.Println(resp.States[0].ICAO24)
}
```

## Disclaimer
This project is not affiliated with the Opensky Network
