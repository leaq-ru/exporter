package consumer

import (
	"github.com/nats-io/stan.go"
	"time"
)

const exportSubjectName = "export"

func (c Consumer) Subscribe() (err error) {
	c.state.sub, err = c.stanConn.QueueSubscribe(
		exportSubjectName,
		c.serviceName,
		c.cb,
		stan.DurableName(exportSubjectName),
		stan.SetManualAckMode(),
		stan.MaxInflight(5),
		stan.AckWait(7*time.Minute),
	)
	if err != nil {
		return
	}

	c.state.subscribeCalledOK = true
	return
}
