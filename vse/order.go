package vse

import(
  "bytes"
  "encoding/json"
  "fmt"
  "io"
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
  path := fmt.Sprintf("/game/%s/portfolio/orders", p.game)

  resp, err := p.c.doRequest("GET", path, nil)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()

  doc, err := goquery.NewDocumentFromReader(io.Reader(resp.Body))
  if err != nil {
    return nil, err
  }

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

// Example request body: [{Fuid: "STOCK-XNAS-AAPL", Shares: "1", Type: "Buy", Term: "Cancelled"}]
func(p *Portfolio) SubmitOrder(order Order) error {
  path := fmt.Sprintf("/game/%s/trade/submitorder", p.game)

  // Convert Order to JSON
  reqBody := append([]Order(nil), order)
  jsonStr, err := json.Marshal(reqBody)
  if err != nil {
    log.Error(err)
    return err
  }

  log.Debug(string(jsonStr))

  resp, err := p.c.doRequest("POST", path, bytes.NewBuffer(jsonStr))
  if err != nil {
    return err
  }

  defer resp.Body.Close()

  return nil
}

func(p *Portfolio) CancelOrder(id string) error {
  uri := fmt.Sprintf("https://marketwatch.com/game/%s/trade/cancelorder?id=%s", p.game, id)
  resp, err := p.c.config.HttpClient.Get(uri)
  if err != nil {
    return err
  }

  defer resp.Body.Close()

  return nil
}
