package skylight

type Option func(c *Client)

func WithLevel(l Level) Option {
	return func(c *Client) {
		c.WithLevel(l)
	}
}

func WithObserver(o ...*Observer) Option {
	return func(c *Client) {
		c.WithObserver(o...)
	}
}

func WithStandardLogger() Option {
	return func(c *Client) {
		c.WithStandardLogger()
	}
}
