package vse

import (
  "crypto/tls"
  "encoding/json"
  "fmt"
  "log"
  "io"
  "io/ioutil"
  "net/http"
  "net/http/cookiejar"
  "net/url"
  "os"

  "github.com/PuerkitoBio/goquery"
  // "golang.org/x/net/publicsuffix"
)

type AuthResponse struct {
  Username string `json:"username"`
  Url      string `json:"url"`
  Uuid     string `json:"uuid"`
  FName    string `json:"fname"`
  LName    string `json:"lname"`
  Result   string `json:"result"`
}

func checkError(err error){
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }
}

func getJson(url string, target interface{}) error {
    r, err := http.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

    return json.NewDecoder(r.Body).Decode(target)
}

// Need to do GET request to the url in the JSON response from id.marketwatch.com
func Authenticate(username string, password string) {
  tr := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }

  // options := cookiejar.Options{
  //   PublicSuffixList: publicsuffix.List,
  // }

  jar, err := cookiejar.New(nil)
  checkError(err)

  uri := fmt.Sprintf("https://id.marketwatch.com/auth/submitlogin.json?username=%s&password=%s", username, password)
  client := &http.Client{Jar: jar, Transport: tr}
  respBody := AuthResponse{}

  resp, err := client.Get(uri)
  checkError(err)

  err = json.NewDecoder(resp.Body).Decode(&respBody)
  checkError(err)

  resp, err = client.Get(respBody.Url)
  checkError(err)

  for _, cookie := range resp.Cookies() {
    log.Println(cookie.Name)
  }

  defer resp.Body.Close()

  resp, err = client.Get("https://www.marketwatch.com/game/sim101/portfolio/Holdings")
  checkError(err)

  // Debug: Log Cookies from both domains
  curl, _ := url.Parse("https://www.marketwatch.com")
  durl, _ := url.Parse("https://id.www.marketwatch.com")
  log.Println(jar.Cookies(curl))
  log.Println(jar.Cookies(durl))

  contents, err := ioutil.ReadAll(resp.Body)
  checkError(err)
  fmt.Printf("%s\n", string(contents))

  doc, err := goquery.NewDocumentFromReader(io.Reader(resp.Body))

  doc.Find(".highlight").Each(func(i int, s *goquery.Selection) {
    log.Println("reached here")
    str, exists := s.Attr("data-ticker")
    if exists {
      u, err := url.Parse(str)
      checkError(err)
      m, _ := url.ParseQuery(u.RawQuery)
      fmt.Println("\033[1;35m"+s.Text()+"\033[0m", m["q"][0])
    } else {
      fmt.Println(s.Text())
    }
  })

  log.Println("reached exit")
  // data, err := ioutil.ReadAll(resp.Body)
  // fmt.Printf("%s\n", string(data))
}
