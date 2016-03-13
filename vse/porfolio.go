package vse

import (
  "fmt"
  // "net/url"
  "io"
  "strings"

  "github.com/PuerkitoBio/goquery"
)

type Holding struct {
  Symbol    string
  Last      uint64
  MrktValue uint64
  Shares    uint64
  Positiion string
}

type Holdings []*Holding

type Portfolio struct {
  c *Client
  game string
}

func (c *Client) Portfolio(game string) *Portfolio {
  return &Portfolio{
    c: c,
    game: game,
  }
}

func (p *Portfolio) GetHoldings() {
  uri := fmt.Sprintf("https://www.marketwatch.com/game/%s/portfolio/Holdings", p.game)

  resp, err := p.c.config.HttpClient.Get(uri)
  checkError(err)

  defer resp.Body.Close()

  doc, err := goquery.NewDocumentFromReader(io.Reader(resp.Body))
  checkError(err)

  doc.Find("table.highlight tbody").First().Find("tr").Each(func(i int, s *goquery.Selection) {
    symbol := s.Find("td h2 a").First().Text()
    fmt.Printf("Symbol: %s | ", symbol)
    last := s.Find("td.numeric p.last.primaryfield").Text()
    fmt.Printf("Last: %s | ", last)
    value := s.Find("td.numeric p.marketvalue.primaryfield").Text()
    fmt.Printf("Value: %s | ", strings.TrimSpace(value))
    position := s.Find("td.equity p.secondaryfield").Text()
    fmt.Printf("Position: %s\n", strings.TrimSpace(position))
  })
}
