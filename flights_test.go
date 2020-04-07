package goflight_test

import (
	"encoding/json"
	"errors"
	"github.com/marcelblijleven/goflight"
	"io/ioutil"
	"net/url"
	"testing"
	"time"
)

func TestFlight_UnmarshalJSON(t *testing.T) {
	data, err := ioutil.ReadFile("./mocks/flights.json")

	if err != nil {
		t.Fatal("unexpected error while reading mock file flights.json")
	}

	var result []goflight.Flight

	if err = json.Unmarshal(data, &result); err != nil {
		t.Error(err.Error())
	}

	if len(result) != 2 {
		t.Error("expected length of flights to be 2")
	}
}

var getFlightInTimeTests = []struct {
	label       string
	begin       time.Time
	end         time.Time
	errExpected error
	statusCode  int
	mockFile    string
}{
	{
		"Correct call",
		time.Now(),
		time.Now().Add(time.Hour),
		nil,
		200,
		"./mocks/flights.json",
	},
	{
		"Time range too big",
		time.Now(),
		time.Now().Add(time.Hour * 3),
		goflight.ErrTimeRangeTooBig,
		0,
		"",
	},
	{
		"End before begin",
		time.Now().Add(time.Hour),
		time.Now(),
		goflight.ErrEndBeforeBegin,
		0,
		"",
	},
	{
		"No results found",
		time.Now(),
		time.Now().Add(time.Hour),
		nil,
		404,
		"",
	},
	{
		"Server error",
		time.Now(),
		time.Now().Add(time.Hour),
		errors.New("500 Internal Server Error"),
		500,
		"",
	},
}

func TestFlightService_GetFlightsInTime(t *testing.T) {
	for _, tt := range getFlightInTimeTests {
		t.Run(tt.label, func(t *testing.T) {
			if tt.mockFile == "" && tt.errExpected != nil {
				// Checking returns errors
				if tt.statusCode == 0 {
					// Checking errors where calls to the http client are not made
					client, err := goflight.NewClient("", "", nil)

					if err != nil {
						t.Fatal("unexpected error setting up new Goflight client")
					}

					flights, err := client.Flights.GetFlightsInTime(tt.begin, tt.end)

					if err == nil {
						t.Error("expected err to be non nil")
					}

					if !errors.As(err, &tt.errExpected) {
						t.Errorf("expected error to be as %v", tt.errExpected.Error())
					}

					if flights != nil {
						t.Error("expected flights to be nil")
					}
				}

				if tt.statusCode != 0 {
					// Checking errors where calls to the http client are made
					mockHandler := CreateTestHandler(tt.statusCode, nil)
					mockHTTPClient, closeServer := HTTPTestClient(mockHandler)
					defer closeServer()

					client, err := goflight.NewClient("", "", mockHTTPClient)
					u, _ := url.Parse("http://example.com")
					goflight.SetBaseURL(client, u)

					if err != nil {
						t.Fatal("unexpected error while creating new Gofight client")
					}

					flights, err := client.Flights.GetFlightsInTime(tt.begin, tt.end)

					if err == nil {
						t.Error("expected err to be non nil")
					}

					if err.Error() != tt.errExpected.Error() {
						t.Errorf("expected %v to equal %v", tt.errExpected.Error(), err.Error())
					}

					if flights != nil {
						t.Error("expected flights to be nil")
					}
				}
			}

			if tt.mockFile != "" {
				// Checking requests where a mock response is expected
				body, err := ioutil.ReadFile(tt.mockFile)

				if err != nil {
					t.Fatal("unexpected error while reading mock file flights.json")
				}

				mockHandler := CreateTestHandler(tt.statusCode, body)
				mockHTTPClient, closeServer := HTTPTestClient(mockHandler)
				defer closeServer()

				client, err := goflight.NewClient("", "", mockHTTPClient)
				u, _ := url.Parse("http://example.com")
				goflight.SetBaseURL(client, u)

				if err != nil {
					t.Fatal("unexpected error while creating new Gofight client")
				}

				flights, err := client.Flights.GetFlightsInTime(tt.begin, tt.end)

				if err != nil {
					t.Error("expected err to be nil")
				}

				if flights == nil {
					t.Error("expected flights to be non nil")
				}

				if len(flights) != 2 {
					t.Error("expected flights to have a length of 2")
				}
			}
		})
	}
}
