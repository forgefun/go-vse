package vse

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	log "github.com/Sirupsen/logrus"
)

const (
	defaultScheme  = "https"
	defaultAddress = "www.marketwatch.com"
)

type AuthResponse struct {
	Username string `json:"username"`
	Url      string `json:"url"`
	Uuid     string `json:"uuid"`
	FName    string `json:"fname"`
	LName    string `json:"lname"`
	Result   string `json:"result"`
}

type Config struct {
	Scheme     string
	Address    string
	HttpClient *http.Client
	Username   string
	Password   string
}

type Client struct {
	config Config
}

type request struct {
	method string
	url    *url.URL
	params url.Values
	body   io.Reader
	obj    interface{}
}

// Default configuration for the client
func DefaultConfig() *Config {
	// Create cookie jar`
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Warn("Cannot create cookie jar: %s", err)
	}

	// Set transport to accept unsecure cert due to SSL on the site
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	config := &Config{
		Scheme:  defaultScheme,
		Address: defaultAddress,
		HttpClient: &http.Client{
			Jar:       jar,
			Transport: transport,
		},
	}

	if username := os.Getenv("USERNAME"); username != "" {
		config.Username = username
	}

	if password := os.Getenv("PASSWORD"); password != "" {
		config.Password = password
	}

	return config
}

// Creates new VSE client that is authenticated
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("No configuration provided")
	}

	// Bootstrap default config
	defConfig := DefaultConfig()

	if len(config.Username) == 0 {
		config.Username = defConfig.Username
	}

	if len(config.Password) == 0 {
		config.Password = defConfig.Password
	}

	if config.HttpClient == nil {
		config.HttpClient = defConfig.HttpClient
	}

	//Perform authentication with the config
	authenticate(config)

	client := &Client{
		config: *config,
	}

	return client, nil
}

func (c *Client) newRequest(method string, path string, params map[string][]string) *request {
	r := &request{
		method: method,
		url: &url.URL{
			Scheme: c.config.Scheme,
			Host:   c.config.Address,
			Path:   path,
		},
		params: params,
	}
	return r
}

func (c *Client) doRequest(r *request) (*http.Response, error) {
	// Encode query parameters
	r.url.RawQuery = r.params.Encode()

	// Encode body if object exists and not encoded yet
	if r.body == nil && r.obj != nil {
		buf := bytes.NewBuffer(nil)
		enc := json.NewEncoder(buf)
		if err := enc.Encode(r.obj); err != nil {
			return nil, err
		}
		r.body = buf
	}

	// Create HTTP request
	req, err := http.NewRequest(r.method, r.url.RequestURI(), r.body)
	if err != nil {
		return nil, err
	}

	req.URL.Scheme = r.url.Scheme
	req.URL.Host = r.url.Host
	req.Host = r.url.Host

	// Add application/json header if body is not empty
	if r.body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	log.Debug(r.url.RequestURI())

	resp, err := c.config.HttpClient.Do(req)
	return resp, err
}

// Perform authentication to get back auth cookies
func authenticate(config *Config) error {
	log.Debug("Authenticating...")

	client := config.HttpClient
	uri := fmt.Sprintf("https://id.marketwatch.com/auth/submitlogin.json?username=%s&password=%s", config.Username, config.Password)
	respBody := AuthResponse{}

	// Hit login to get url back from the response
	resp, err := client.Get(uri)
	if err != nil {
		return err
	}

	// Decode the body into struct
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}

	// GET request to the url in the JSON response from id.marketwatch.com
	resp, err = client.Get(respBody.Url)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()
	defer log.Debug("Authenticated!")

	return nil
}
