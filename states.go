package goflight

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type statesService struct {
	client *Client
}

// StateVector represents the state of the aircraft at a given time
type StateVector struct {
	ICAO24         string   // Unique ICAO 24-bit address of the transponder in hex string representation.
	Callsign       *string  // Callsign of the vehicle (8 chars). Can be nil if no callsign has been received.
	OriginCountry  string   // Country name inferred from the ICAO 24-bit address.
	TimePosition   *int64   // Unix timestamp (seconds) for the last position update. Can be nil if no position report was received by OpenSky within the past 15s.
	LastContact    int64    // Unix timestamp (seconds) for the last update in general. This field is updated for any new, valid message received from the transponder.
	Longitude      *float64 // WGS-84 longitude in decimal degrees. Can be nil.
	Latitude       *float64 // WGS-84 latitude in decimal degrees. Can be nil.
	BaroAltitude   *float64 // Barometric altitude in meters. Can be nil.
	OnGround       bool     // Boolean value which indicates if the position was retrieved from a surface position report.
	Velocity       *float64 // Velocity over ground in m/s. Can be nil.
	TrueTrack      *float64 // True track in decimal degrees clockwise from north (north=0°). Can be nil.
	VerticalRate   *float64 // Vertical rate in m/s. A positive value indicates that the airplane is climbing, a negative value indicates that it descends. Can be nil.
	Sensors        *[]int   // IDs of the receivers which contributed to this state vector. Is nil if no filtering for sensor was used in the request.
	GeoAltitude    *float64 // Geometric altitude in meters. Can be nil.
	Squawk         *string  // The transponder code aka Squawk. Can be nil.
	Spi            *bool    // Whether flight status indicates special purpose indicator. Can be nil
	PositionSource int      // Origin of this state’s position: 0 = ADS-B, 1 = ASTERIX, 2 = MLAT
}

// StatesResponse is the response retrieved from the /api/states/all and /api/states/own endpoints
type StatesResponse struct {
	Time   int64         `json:"time"`
	States []StateVector `json:"states"`
}

// UnmarshalJSON unmarshals the provided []byte onto a StateVector struct
func (s *StateVector) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{
		&s.ICAO24,
		&s.Callsign,
		&s.OriginCountry,
		&s.TimePosition,
		&s.LastContact,
		&s.Longitude,
		&s.Latitude,
		&s.BaroAltitude,
		&s.OnGround,
		&s.Velocity,
		&s.TrueTrack,
		&s.VerticalRate,
		&s.Sensors,
		&s.GeoAltitude,
		&s.Squawk,
		&s.Spi,
		&s.PositionSource,
	}

	expectedLen := len(tmp)

	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}

	if actual, expected := len(tmp), expectedLen; actual != expected {
		return errors.New("incorrect number of fields in StateVector json")
	}

	return nil
}

func (s *statesService) getStatesRequest(endpoint string, timeParam time.Time, icao24 string) (*http.Request, error) {
	method := "GET"

	e, err := url.Parse(endpoint)

	if err != nil {
		return nil, err
	}

	u := s.client.baseURL.ResolveReference(e)

	req, err := http.NewRequest(method, u.String(), nil)

	if err != nil {
		return nil, err
	}

	if timeParam, ok := checkTime(timeParam); ok {
		req.URL.Query().Add("time", strconv.FormatInt(timeParam.Unix(), 10))
	}

	if icao24, ok := checkString(icao24); ok {
		req.URL.Query().Add("icao24", icao24)
	}

	req.URL.RawQuery = req.URL.Query().Encode()

	return req, nil
}

// GetAllStates returns the response of /api/states/all
func (s *statesService) GetAllStates(time time.Time, icao24 string) (StatesResponse, error) {
	endpoint := "/api/states/all"
	req, err := s.getStatesRequest(endpoint, time, icao24)

	if err != nil {
		return StatesResponse{}, err
	}

	username, okUser := checkString(s.client.username)
	password, okPass := checkString(s.client.password)

	if okUser && okPass {
		// Authentication is not required for this endpoint
		req.SetBasicAuth(username, password)
	}

	resp, err := s.client.httpClient.Do(req)
	var statesResponse StatesResponse

	if err != nil {
		return statesResponse, err
	}

	if resp.StatusCode/100 != 2 {
		if resp.StatusCode == 403 {
			return statesResponse, ErrUnauthorizedAccess
		}

		return statesResponse, fmt.Errorf("%v - %v", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return statesResponse, err
	}

	if err = json.Unmarshal(body, &statesResponse); err != nil {
		return statesResponse, err
	}

	return statesResponse, nil
}

// GetOwnStates returns the response of /api/states/own
func (s *statesService) GetOwnStates(time time.Time, icao24 string) (*StatesResponse, error) {
	endpoint := "/api/states/own"
	req, err := s.getStatesRequest(endpoint, time, icao24)

	if err != nil {
		return nil, err
	}

	username, okUser := checkString(s.client.username)
	password, okPass := checkString(s.client.password)

	if !(okUser && okPass) {
		// Authentication is required for this endpoint
		return nil, ErrInvalidCredentials
	}

	req.SetBasicAuth(username, password)

	resp, err := s.client.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response *StatesResponse
	if err = json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	return response, nil
}
