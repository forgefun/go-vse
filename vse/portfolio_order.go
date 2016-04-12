package vse

import (
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

type OrderRequest []Order

type PendingOrder struct {
	Id     string
	Symbol string
	Shares uint64
	Type   string
	Date   string
}

type PendingOrders []*PendingOrder

type CancelOrderParams struct {
	Id string `url:id,omitempty`
}

func (p *Portfolio) ListOrders() (PendingOrders, error) {
	path := fmt.Sprintf("/game/%s/portfolio/orders", p.game)

	req := p.c.newRequest("GET", path, nil)

	resp, err := p.c.doRequest(req)
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

// Available prefixes: STOCK-XNAS-, EXCHANGETRADEDFUND-XASQ-
// Example request body: [{Fuid: "STOCK-XNAS-AAPL", Shares: "1", Type: "Buy", Term: "Cancelled"}]
// TODO: Handle case where there is not enough buying power to process order
func (p *Portfolio) SubmitOrder(order Order) error {
	path := fmt.Sprintf("/game/%s/trade/submitorder", p.game)

	reqObj := append(OrderRequest{}, order)
	log.Debug(reqObj)

	req := p.c.newRequest("POST", path, nil)
	req.obj = reqObj

	log.Debug("%v", req)

	resp, err := p.c.doRequest(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (p *Portfolio) CancelOrder(id string) error {
	path := fmt.Sprintf("/game/%s/trade/cancelorder", p.game)

	// TODO: Constructing params can be refactored
	params := &url.Values{}
	params.Set("id", id)

	req := p.c.newRequest("GET", path, nil)
	req.params = *params

	log.Debug(req.params)

	resp, err := p.c.doRequest(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
