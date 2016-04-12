package vse

type Portfolio struct {
	c    *Client
	game string
}

func (c *Client) Portfolio(game string) *Portfolio {
	return &Portfolio{
		c:    c,
		game: game,
	}
}
