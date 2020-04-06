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

var stateVectorTests = []struct {
	inputFile string
}{
	{"./mocks/state_vector.json"},
	{"./mocks/state_vector_null_values.json"},
}

func TestStateVector_UnmarshalJSON(t *testing.T) {
	for _, tt := range stateVectorTests {
		t.Run(tt.inputFile, func(t *testing.T) {
			data, err := ioutil.ReadFile(tt.inputFile)

			if err != nil {
				t.Fatal(err.Error())
			}

			var vector goflight.StateVector

			if err = json.Unmarshal(data, &vector); err != nil {
				t.Errorf("unexpected error for %q: %v", tt.inputFile, err.Error())
			}

			// Unmarshal mock data onto slice of interface
			// to check values
			var tmp []interface{}

			if err = json.Unmarshal(data, &tmp); err != nil {
				t.Fatal(err.Error())
			}
		})
	}
}

func TestStateVector_UnmarshalJSON_IncorrectLength(t *testing.T) {
	data, err := ioutil.ReadFile("./mocks/state_vector_incorrect_length.json")

	if err != nil {
		t.Fatal(err.Error())
	}

	var vector goflight.StateVector

	if err = json.Unmarshal(data, &vector); err == nil {
		t.Errorf("expected error to be non nil")
	}
}

var getStatesInputs = []struct {
	icao24Input string
	timeInput   time.Time
}{
	{"", time.Time{}},
	{"c0ffee", time.Date(2020, time.July, 6, 0, 0, 0, 0, time.UTC)},
}

func TestStatesService_GetAllStates(t *testing.T) {
	mockResponseBody, err := ioutil.ReadFile("./mocks/states.json")

	if err != nil {
		t.Fatal("unexpected error in retrieving mock response body")
	}

	mockHandler := CreateTestHandler(200, mockResponseBody)
	mockClient, closeServer := HTTPTestClient(mockHandler)
	defer closeServer()

	for _, tt := range getStatesInputs {
		t.Run(tt.icao24Input, func(t *testing.T) {
			u, _ := url.Parse("http://example.com")

			client, err := goflight.NewClient("", "", mockClient, u)

			if err != nil {
				t.Fatal("unexpected error in setting up Goflight Client")
			}

			response, err := client.States.GetAllStates(tt.timeInput, tt.icao24Input)

			if err != nil {
				t.Error(err.Error())
			}

			var expectedTime int64
			expectedTime = 1586031310

			if response.Time != 1586031310 {
				t.Errorf("expected %v to equal %v", response.Time, expectedTime)
			}

			if len(response.States) != 6 {
				t.Error("expect length of states in response to equal 6")
			}
		})
	}
}

func TestStatesService_GetOwnStates(t *testing.T) {
	mockResponseBody, err := ioutil.ReadFile("./mocks/states.json")

	if err != nil {
		t.Fatal("unexpected error in retrieving mock response body")
	}

	mockHandler := CreateTestHandler(200, mockResponseBody)
	mockClient, closeServer := HTTPTestClient(mockHandler)
	defer closeServer()

	for _, tt := range getStatesInputs {
		t.Run(tt.icao24Input, func(t *testing.T) {
			u, _ := url.Parse("http://example.com")

			client, err := goflight.NewClient("user", "password", mockClient, u)

			if err != nil {
				t.Fatal("unexpected error in setting up Goflight Client")
			}

			response, err := client.States.GetAllStates(tt.timeInput, tt.icao24Input)

			if err != nil {
				t.Error(err.Error())
			}

			var expectedTime int64
			expectedTime = 1586031310

			if response.Time != 1586031310 {
				t.Errorf("expected %v to equal %v", response.Time, expectedTime)
			}

			if len(response.States) != 6 {
				t.Error("expect length of states in response to equal 6")
			}
		})
	}
}

func TestStatesService_GetOwnStates_Authentication(t *testing.T) {
	// Check for unauthorized access
	mockHandler := CreateTestHandler(403, nil)
	mockClient, closeServer := HTTPTestClient(mockHandler)
	defer closeServer()

	client, err := goflight.NewClient("unauthorizedUser", "secret", mockClient, nil)

	if err != nil {
		t.Fatal("unexpected error in setting up Goflight Client")
	}

	_, err = client.States.GetOwnStates(time.Time{}, "")

	if !errors.As(err, &goflight.UnauthorizedAccessError) {
		t.Errorf("expected error to be: %v", goflight.UnauthorizedAccessError.Error())
	}

	// Check for invalid credentials
	client, err = goflight.NewClient("", "", mockClient, nil)

	if err != nil {
		t.Fatal("unexpected error in setting up Goflight Client")
	}

	_, err = client.States.GetOwnStates(time.Time{}, "")

	if !errors.Is(err, goflight.ErrInvalidCredentials) {
		t.Errorf("expected error to be: %v", goflight.UnauthorizedAccessError.Error())
	}
}
