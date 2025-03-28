package client

import (
	"time"

	"github.com/goburrow/modbus"
)

type ModbusTCP struct {
	ID    string
	Start uint16
	Size  uint16
	Bytes []byte
	// SlaveId   hex
	SlaveId byte
	Error   error
}

type SyncModbusTCPGroup struct {
	Addr       string
	ModbusTCPs []*ModbusTCP
	Retry      int
	RetryDelay time.Duration
	Timeout    time.Duration
	Error      error
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
		s.Error = timeoutError
		return
	}

	defer h.Close()
	for _, mb := range s.ModbusTCPs {
		h.SlaveId = mb.SlaveId
		c := modbus.NewClient(h)
		b, err := c.ReadHoldingRegisters(mb.Start, mb.Size)
		mb.Error = err
		copy(mb.Bytes, b)
		time.Sleep(s.RetryDelay)
	}
}
