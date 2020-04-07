package goflight

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const flightsPrefix string = "/api/flights/"

type flightService struct {
	client *Client
}

// Flight represents a flight in an /api/flights/* response
type Flight struct {
	ICAO24                                string  `json:"icao24"`                           // Unique ICAO 24-bit address of the transponder in hex string representation. All letters are lower case.
	FirstSeen                             int64   `json:"firstSeen"`                        // Estimated time of departure for the flight as Unix time (seconds since epoch).
	EstDepartureAirport                   *string `json:"estDepartureAirport"`              // ICAO code of the estimated departure airport. Can be null if the airport could not be identified.
	LastSeen                              int64   `json:"lastSeen"`                         // Estimated time of arrival for the flight as Unix time (seconds since epoch)
	EstArrivalAirport                     *string `json:"estArrivalAirport"`                // ICAO code of the estimated arrival airport. Can be null if the airport could not be identified.
	CallSign                              *string `json:"callsign"`                         // Callsign of the vehicle (8 chars). Can be null if no callsign has been received. If the vehicle transmits multiple callsigns during the flight, we take the one seen most frequently
	EstDepartureAirportHorizontalDistance *int64  `json:"estDepartureAirportHorizDistance"` // Horizontal distance of the last received airborne position to the estimated departure airport in meters
	EstDepartureAirportVerticalDistance   *int64  `json:"estDepartureAirportVertDistance"`  // Vertical distance of the last received airborne position to the estimated departure airport in meters
	EstArrivalAirportHorizontalDistance   *int64  `json:"estArrivalAirportHorizDistance"`   // Horizontal distance of the last received airborne position to the estimated arrival airport in meters
	EstArrivalAirportVerticalDistance     *int64  `json:"estArrivalAirportVertDistance"`    // Vertical distance of the last received airborne position to the estimated arrival airport in meters
	DepartureAirportCandidatesCount       int64   `json:"departureAirportCandidatesCount"`  // Number of other possible departure airports. These are airports in short distance to estDepartureAirport.
	ArrivalAirportCandidatesCount         int64   `json:"arrivalAirportCandidatesCount"`    // Number of other possible departure airports. These are airports in short distance to estArrivalAirport.
}

func (f *flightService) GetFlightsInTime(begin, end time.Time) ([]Flight, error) {
	endpoint, err := url.Parse(flightsPrefix + "all")

	if err != nil {
		return nil, err
	}

	if end.Before(begin) {
		return nil, ErrEndBeforeBegin
	}

	if end.Sub(begin) > time.Hour*2 {
		return nil, ErrTimeRangeTooBig
	}

	u := f.client.baseURL.ResolveReference(endpoint)
	params := url.Values{}
	params.Add("begin", strconv.FormatInt(begin.Unix(), 10))
	params.Add("end", strconv.FormatInt(end.Unix(), 10))
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)

	if err != nil {
		return nil, err
	}

	username, okUser := checkString(f.client.username)
	password, okPassword := checkString(f.client.password)

	if okUser && okPassword {
		req.SetBasicAuth(username, password)
	}

	resp, err := f.client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		if resp.StatusCode == 404 {
			return []Flight{}, nil
		}

		return nil, fmt.Errorf(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result []Flight

	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}
