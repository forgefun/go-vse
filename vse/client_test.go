package vse

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

var (
	testUsername = "test"
	testPassword = "password"
)

// Set up test mux server and initiate a client
func setup(t *testing.T) (*http.ServeMux, func()) {
	// Create multiplexer and server for mock API response
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	// Create a new default config, and use the test server URL
	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	url, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	client.config.Scheme = url.Scheme
	client.config.Address = url.Host
	client.config.Username = testUsername
	client.config.Password = testPassword

	// Closure on server
	return mux, func() { server.Close() }
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestDefaultConfig_env(t *testing.T) {
	os.Setenv("USERNAME", testUsername)
	os.Setenv("PASSWORD", testPassword)
	defer os.Unsetenv("USERNAME")
	defer os.Unsetenv("PASSWORD")

	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	if strings.Compare(client.config.Username, testUsername) != 0 {
		t.Fatalf("Expected: %s, Got: %s", client.config.Username, testUsername)
	}

	if strings.Compare(client.config.Password, testPassword) != 0 {
		t.Fatalf("Expected: %s, Got: %s", client.config.Password, testPassword)
	}
}

func TestNewClientDefaultConfig(t *testing.T) {
	// New client with default config
	c, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	// Default config
	config := DefaultConfig()

	if !reflect.DeepEqual(c.config.HttpClient, config.HttpClient) {
		t.Errorf("Client's default config: %+v, default config: %+v", c.config, config)
	}
}
