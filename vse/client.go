package vse

import (
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

type AuthResponse struct {
  Username string `json:"username"`
  Url      string `json:"url"`
  Uuid     string `json:"uuid"`
  FName    string `json:"fname"`
  LName    string `json:"lname"`
  Result   string `json:"result"`
}

type Config struct {
  HttpClient *http.Client
  Username   string
  Password   string
}

type Client struct {
  config Config
}

// Default configuration for the client
func DefaultConfig() *Config {
  // Create cookie jar`
  jar, err := cookiejar.New(nil)
  checkError(err)

  // Set transport to accept unsecure cert due to SSL on the site
  transport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }

  config := &Config{
    HttpClient: &http.Client{
      Jar: jar,
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
  // Bootstrap default config
  defConfig := DefaultConfig()

  if config.Username == "" {
    config.Username = defConfig.Username
  }

  if config.Password == "" {
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

// TODO: Make it also work with POST requests
func (c *Client) doRequest(method string, path string, body io.Reader) (*http.Response, error) {
  url := &url.URL{
    Scheme: "https",
    Host:   "www.marketwatch.com",
    Path:   path,
  }

  req, err := http.NewRequest(method, url.RequestURI(), body)
  if err != nil {
    return nil, err
  }

  // Set header if body if is not nil
  // if body != nil {
  //   req.Header.Set("Content-Type", "application/json")
  //   req.Header.Set("X-Requested-With", "XMLHttpRequest")
  // }

  req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Host = url.Host

  resp, err := c.config.HttpClient.Do(req)
  return resp, err
}

// Perform authentication to get back auth cookies
func authenticate(config *Config) {
  log.Debug("Authenticating...")

  client := config.HttpClient
  uri := fmt.Sprintf("https://id.marketwatch.com/auth/submitlogin.json?username=%s&password=%s", config.Username, config.Password)
  respBody := AuthResponse{}

  // Hit login to get url back from the response
  resp, err := client.Get(uri)
  checkError(err)

  // Decode the body into struct
  err = json.NewDecoder(resp.Body).Decode(&respBody)
  checkError(err)

  // GET request to the url in the JSON response from id.marketwatch.com
  resp, err = client.Get(respBody.Url)
  checkError(err)

  defer resp.Body.Close()
  defer log.Debug("Authenticated!")
}
