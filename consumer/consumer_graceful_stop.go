package consumer

import "time"

func (c Consumer) GracefulStop() {
	c.state.drain = true

	err := c.stanConn.Close()
	if err != nil {
		c.logger.Error().Err(err).Send()
	}
	close(c.state.done)

	for {
		if c.state.gracefulOK {
			return
		}
		time.Sleep(time.Second)
	}
}
