package consumer

import (
	"github.com/nats-io/stan.go"
)

const exportSubjectName = "export"

func (c Consumer) Subscribe() (err error) {
	c.state.sub, err = c.stanConn.QueueSubscribe(
		exportSubjectName,
		c.serviceName,
		c.cb,
		stan.DurableName(exportSubjectName),
		stan.SetManualAckMode(),
		stan.MaxInflight(1),
	)
	if err != nil {
		return
	}

	c.state.subscribeCalledOK = true
	return
}
