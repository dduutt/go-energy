package meter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dduutt/go-energy/client"
)

const (
	TIMEOUT     = 1 * time.Second
	RETRY       = 3
	RETRY_DELAY = 1 * time.Second
)

type Energy struct {
	Code        string
	WorkShop    string
	Room        string
	Name        string
	Protocol    string
	IP          string
	Port        int
	SlaveOrArea string
	Start       int
	Size        int
	DataType    string
	IsBigEndian bool // 添加字节序配置
	Bytes       []byte
	Value       float64
	BigEndian   bool
}

type AddrGroup struct {
	Addr       string
	Meters     []*Energy
	Protocol   string
	Retry      int
	RetryDelay time.Duration
	Timeout    time.Duration
}

func (a *AddrGroup) Read(rc chan *client.GroupSyncReadResult) {
	var g client.GroupSyncReader
	switch a.Protocol {
	case "modbus_tcp":
		mbg := &client.SyncModbusTCPGroup{
			Addr:       a.Addr,
			Retry:      a.Retry,
			RetryDelay: a.RetryDelay,
			Timeout:    a.Timeout,
		}
		for _, m := range a.Meters {
			slaveId, err := strconv.ParseUint(m.SlaveOrArea, 10, 8)
			if err != nil {
				continue
			}
			mbc := &client.ModbusTCP{
				ID:      m.Code,
				Start:   uint16(m.Start),
				Size:    uint16(m.Size),
				SlaveId: byte(slaveId),
			}
			mbg.ModbusTCPs = append(mbg.ModbusTCPs, mbc)

		}
		g = mbg
	case "s7_200_smart":
		addr := strings.Split(a.Addr, ":")[0]
		s7g := &client.SyncS7Group{
			Addr:       addr,
			Rack:       0,
			Slot:       1,
			Retry:      a.Retry,
			RetryDelay: a.RetryDelay,
			Timeout:    a.Timeout,
		}
		for _, m := range a.Meters {
			s7 := &client.S7{
				Code:    m.Code,
				Address: m.SlaveOrArea,
				Size:    m.Size,
				Start:   m.Start,
				Byte:    make([]byte, m.Size),
			}
			s7g.S7s = append(s7g.S7s, s7)
		}
		g = s7g
	default:
		err := fmt.Errorf("protocol %s not support", a.Protocol)
		rc <- &client.GroupSyncReadResult{Error: err}
		return
	}
	g.Read(rc)
}

func GroupByAddrFromExcel(meters []*Energy) map[string]*AddrGroup {

	g := make(map[string]*AddrGroup)
	for _, m := range meters {
		addr := fmt.Sprintf("%s:%d", m.IP, m.Port)
		if _, ok := g[addr]; !ok {
			g[addr] = &AddrGroup{
				Addr:       addr,
				Meters:     make([]*Energy, 0),
				Protocol:   m.Protocol,
				Timeout:    TIMEOUT,
				Retry:      RETRY,
				RetryDelay: RETRY_DELAY,
			}
		}
		g[addr].Meters = append(g[addr].Meters, m)
	}
	return g
}
func compareHeader(header []string) bool {
	if len(header) != len(H) {
		return false
	}
	for i, v := range header {
		if v != H[i] {
			return false
		}
	}
	return true
}

func (e *Energy) ParseFloat() (err error) {
	if e.Bytes == nil {
		return fmt.Errorf("no byte data")
	}

	switch e.DataType {
	case "float32":
		var f float32
		err = e.Read(&f, e.Bytes)
		e.Value = float64(f)
	case "float64":
		var f float64
		err = e.Read(&f, e.Bytes)
		e.Value = f
	case "int16":
		var i int16
		err = e.Read(&i, e.Bytes)
		e.Value = float64(i)
	case "int32":
		var i int32
		err = e.Read(&i, e.Bytes)
		e.Value = float64(i)
	case "int64":
		var i int64
		err = e.Read(&i, e.Bytes)
		e.Value = float64(i)
	case "uint16":
		var i uint16
		err = e.Read(&i, e.Bytes)
		e.Value = float64(i)
	case "uint32":
		var i uint32
		err = e.Read(&i, e.Bytes)
		e.Value = float64(i)
	default:
		err = fmt.Errorf("unknown data type: %s", e.DataType)
	}
	return

}

func (e *Energy) Read(v any, b []byte) error {

	r := bytes.NewReader(b)
	order := binary.ByteOrder(binary.BigEndian)
	if !e.BigEndian {
		order = binary.LittleEndian
	}
	return binary.Read(r, order, v)
}
