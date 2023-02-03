package proxy

import (
	"fmt"
	"net"
	"time"

	"github.com/andig/mbserver"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/modbus"
)

func StartProxy(port int, config modbus.Settings, readOnly bool) error {
	conn, err := modbus.NewConnection(config.URI, config.Device, config.Comset, config.Baudrate, modbus.ProtocolFromRTU(config.RTU), config.ID)
	if err != nil {
		return err
	}

	conn.ConnectDelay(2 * time.Second)
	conn.Delay(50 * time.Millisecond)

	h := &handler{
		log:            util.NewLogger(fmt.Sprintf("proxy-%d", port)),
		readOnly:       readOnly,
		RequestHandler: new(mbserver.DummyHandler), // supplies HandleDiscreteInputs
		conn:           conn,
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	h.log.DEBUG.Printf("modbus proxy for %s listening at :%d", config.String(), port)

	srv, err := mbserver.New(h)

	if err == nil {
		err = srv.Start(l)
	}

	return err
}
