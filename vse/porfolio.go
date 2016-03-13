package vse

import (
  "fmt"
  // "net/url"
  "io"
  "strconv"
  "strings"

  "github.com/PuerkitoBio/goquery"
  log "github.com/Sirupsen/logrus"
)

// TODO: Correctly parse currency, DO NOT use float64
type Holding struct {
  Symbol    string
  Last      string
  MrktValue string
  Shares    uint64
  Position  string
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

func (p *Portfolio) GetHoldings() (Holdings, error) {
  uri := fmt.Sprintf("https://www.marketwatch.com/game/%s/portfolio/Holdings", p.game)

  resp, err := p.c.config.HttpClient.Get(uri)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()

  doc, err := goquery.NewDocumentFromReader(io.Reader(resp.Body))
  if err != nil {
    return nil, err
  }

  var holdings []*Holding

  // Parse HTML and pass into Holdings struct
  doc.Find("table.highlight tbody").First().Find("tr").Each(func(i int, s *goquery.Selection) {
    symbol := s.Find("td h2 a").First().Text()

    last := s.Find("td.numeric p.last.primaryfield").Text()

    valueRaw := s.Find("td.numeric p.marketvalue.primaryfield").Text()
    valueRaw = strings.TrimSpace(valueRaw)
    value := strings.TrimPrefix(valueRaw, "$")

    sharesPositionRaw := s.Find("td.equity p.secondaryfield").Text()
    sharesPositionRaw = strings.TrimSpace(sharesPositionRaw)

    position := strings.TrimLeft(sharesPositionRaw, "1234567890 /")

    cutset := fmt.Sprintf("%s/ ", position)
    sharesRaw := strings.TrimRight(sharesPositionRaw, cutset)
    shares, err := strconv.ParseUint(sharesRaw, 10, 64)
    if err != nil {
      log.Error(err)
    }

    // Construct Holding struct and append it to the holdings slice
    h := &Holding{
      Symbol:    symbol,
      Last:      last,
      MrktValue: value,
      Shares:    shares,
      Position:  position,
    }

    holdings = append(holdings, h)
    info := fmt.Sprintf("Symbol: %s | Last: %s | Value: %s | Shares: %d | Position: %s", h.Symbol, h.Last, h.MrktValue, h.Shares, h.Position)
    log.Debug(info)
  })

  return holdings, nil
}
