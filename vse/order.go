package vse

import(
  "bytes"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "net/url"
  "strconv"
  "strings"

  "github.com/PuerkitoBio/goquery"
  log "github.com/Sirupsen/logrus"
)

type Order struct {
  Fuid   string
  Symbol string
  Shares string
  Limit  string
  Stop   string
  Term   string
  Type   string
}

type OrderRequest struct {
  Collection []Order
}

type PendingOrder struct {
  Id     string
  Symbol string
  Shares uint64
  Type   string
  Date   string
}

type PendingOrders []*PendingOrder

func(p *Portfolio) ListOrders() (PendingOrders, error) {
  uri := fmt.Sprintf("https://www.marketwatch.com/game/%s/portfolio/orders", p.game)

  resp, err := p.c.config.HttpClient.Get(uri)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()

  doc, err := goquery.NewDocumentFromReader(io.Reader(resp.Body))
  if err != nil {
    return nil, err
  }

  // var po []*PendingOrders
  var pendingOrders []*PendingOrder

  // TODO: Gracefully handle errors
  doc.Find("section.portfolio.tabular tr td").Each(func(i int, s *goquery.Selection) {
    po := &PendingOrder{}
    switch i {
    case 0:
      symbol := strings.TrimSpace(s.Text())
      po.Symbol = symbol
      log.Debug(symbol)
    case 1:
      shares := strings.TrimSpace(s.Text())
      po.Shares, _ = strconv.ParseUint(shares, 10, 64)
      log.Debug(shares)
    case 2:
      // NOTE: Non-optimal at the moment
      rawOrderType := stringMinifier(s.Text())
      po.Type = rawOrderType
      log.Debug(fmt.Sprintf("%s", rawOrderType))
    case 3:
      date := strings.TrimSpace(s.Text())
      po.Date = date
      log.Debug(date)
    case 4:
      href, _ := s.Find("a").Attr("href")
      u, _ := url.Parse(href)
      q, _ := url.ParseQuery(u.RawQuery)
      id := q["id"][0]
      po.Id = id
      log.Debug(id)
    }
    pendingOrders = append(pendingOrders, po)
  })
  return pendingOrders, nil
}

func(p *Portfolio) SubmitOrder(order Order) error {
  url := fmt.Sprintf("https://www.marketwatch.com/game/%s/trade/submitorder?week=1", p.game)

  // Convert Order to JSON
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
