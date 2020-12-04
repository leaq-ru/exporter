package exporter_async

func (c Consumer) GracefulStop() {
	err := c.stanConn.Close()
	if err != nil {
		c.logger.Error().Err(err).Send()
	}
	close(c.state.done)
}
