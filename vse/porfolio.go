package vse

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
  "io"
  "strconv"
  "strings"

  "github.com/PuerkitoBio/goquery"
  log "github.com/Sirupsen/logrus"
)

type Portfolio struct {
  c *Client
  game string
}

// TODO: Correctly parse currency, DO NOT use float64
type Holding struct {
  Symbol    string
  Last      string
  MrktValue string
  Shares    uint64
  Position  string
}

type Holdings []*Holding

type Order struct {
  Fuid   string
  Shares string
  Limit  string
  Stop   string
  Term   string
  Type   string
}

type OrderRequest struct {
  Collection []Order
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

func(p *Portfolio) SubmitOrder(order Order) error {
  url := fmt.Sprintf("https://www.marketwatch.com/game/%s/trade/submitorder?week=1", p.game)

  // Convert Order to JSON
  // reqBody := &OrderRequest{}
  reqBody := append([]Order(nil), order)
  jsonStr, err := json.Marshal(reqBody)
  if err != nil {
    log.Error(err)
    return err
  }

  log.Debug(string(jsonStr))
  // var jsonStr = []byte(`[{Fuid: "STOCK-XNAS-AAPL", Shares: "1", Type: "Buy", Term: "Cancelled"}]`)

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  req.Header.Set("X-Requested-With", "XMLHttpRequest")
  req.Header.Set("Content-Type", "application/json")

  resp, err := p.c.config.HttpClient.Do(req)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  log.Debug(resp.Status)

  return nil
}

func(p *Portfolio) CancelOrder(symbol string) error {
  return nil
}
