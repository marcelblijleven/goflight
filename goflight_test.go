package goflight_test

import (
	"github.com/marcelblijleven/goflight"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var newClientTests = []struct {
	username   string
	password   string
	httpClient *http.Client
	host       *url.URL
}{
	{"", "", nil, nil},
	{"test", "tops3cret", &http.Client{Timeout: time.Second * 4}, nil},
}

func TestNewClient(t *testing.T) {
	for _, tt := range newClientTests {
		t.Run(tt.username, func(t *testing.T) {
			client, err := goflight.NewClient(tt.username, tt.password, tt.httpClient, tt.host)

			if err != nil {
				t.Fatal("unexpected error occurred while creating new Goflight client")
			}

			username, password := client.GetCredentials()

			if username != tt.username {
				t.Errorf("expected %v to equal %v", username, tt.username)
			}

			if password != tt.password {
				t.Errorf("expected %v to equal %v", password, tt.password)
			}

			clientBaseURLStr := client.GetBaseURL().String()
			if clientBaseURLStr != goflight.BaseURL {
				t.Errorf("expected %v to equal %v", clientBaseURLStr, goflight.BaseURL)
			}
		})
	}
}

func TestNewClient_CustomHost(t *testing.T) {
	u, err := url.Parse("https://example.com")

	if err != nil {
		t.Fatal("unexpected error occurred while creating custom host URL")
	}

	client, err := goflight.NewClient("", "", nil, u)

	if err != nil {
		t.Fatal("unexpected error occurred while creating new Goflight client")
	}

	if client.GetBaseURL() != u {
		t.Errorf("expected %v to equal %v", client.GetBaseURL(), u)
	}
}
