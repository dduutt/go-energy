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
}

type SyncModbusTCPGroup struct {
	Addr       string
	ModbusTCPs []*ModbusTCP
	Retry      int
	RetryDelay time.Duration
	Timeout    time.Duration
}

func (s *SyncModbusTCPGroup) Read(rc chan *GroupSyncReadResult) {
	gsrr := &GroupSyncReadResult{
		Results: make([]*SyncReadResult, 0),
	}

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
		gsrr.Error = timeoutError
		rc <- gsrr
		return
	}

	defer h.Close()
	for _, mb := range s.ModbusTCPs {
		srr := &SyncReadResult{ID: mb.ID}
		h.SlaveId = mb.SlaveId
		c := modbus.NewClient(h)
		b, err := c.ReadHoldingRegisters(mb.Start, mb.Size)
		if err != nil {
			srr.Error = err
			continue
		}
		srr.Error = err
		srr.Result = b
		gsrr.Results = append(gsrr.Results, srr)
		time.Sleep(s.RetryDelay)
	}
	rc <- gsrr
}
