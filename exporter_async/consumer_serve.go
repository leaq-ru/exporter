package exporter_async

import "errors"

func (c Consumer) Serve() (err error) {
	if !c.state.subscribeCalledOK {
		err = errors.New("you must call Subscribe before Serve")
		return
	}

	return c.pollSubIsValid()
}
