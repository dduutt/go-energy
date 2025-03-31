package client

import (
	"time"

	"github.com/goburrow/modbus"
)

type ModbusTCP struct {
	Code  int
	Start uint16
	Size  uint16
	Bytes []byte
	// SlaveId   hex
	SlaveId byte
	Error   *string
}

type SyncModbusTCPGroup struct {
	Addr       string
	ModbusTCPs []*ModbusTCP
	Retry      int
	RetryDelay time.Duration
	Timeout    time.Duration
	Error      *string
}

func (s *SyncModbusTCPGroup) Read() {

	h := modbus.NewTCPClientHandler(s.Addr)
	h.Timeout = s.Timeout
	var timeoutError error
	for range s.Retry {

		if timeoutError = h.Connect(); timeoutError == nil {
			break
		}
		time.Sleep(s.RetryDelay)
	}
	if timeoutError != nil {
		*s.Error = timeoutError.Error()
		return
	}

	defer h.Close()
	for _, mb := range s.ModbusTCPs {
		h.SlaveId = mb.SlaveId
		c := modbus.NewClient(h)
		b, err := c.ReadHoldingRegisters(mb.Start, mb.Size)
		if err != nil {
			*mb.Error = err.Error()
		}
		copy(mb.Bytes, b)
		time.Sleep(s.RetryDelay)
	}
}
