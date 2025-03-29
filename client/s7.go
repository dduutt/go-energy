package client

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/robinson/gos7"
)

type S7 struct {
	Code    string
	Address string
	Start   int
	Size    int
	Area    string
	AreaNo  int
	Bytes   []byte
	h       gos7.ClientHandler
	Error   *string
}

type SyncS7Group struct {
	Addr       string
	Rack       int
	Slot       int
	S7s        []*S7
	Retry      int
	RetryDelay time.Duration
	Timeout    time.Duration
	Error      *string
}

func (s *SyncS7Group) Read() {

	h := gos7.NewTCPClientHandler(s.Addr, s.Rack, s.Slot)
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
	for _, s7 := range s.S7s {
		s7.h = h
		err := s7.Read()
		if err != nil {
			*s7.Error = err.Error()
		}
		time.Sleep(s.RetryDelay)
	}
}

func (s *S7) Read() error {

	err := s.ParseAddress()
	if err != nil {
		return err
	}

	switch s.Area {
	case "DB":
		return s.ReadDB()
	default:
		return fmt.Errorf("[unknown area] %s %s", s.Code, s.Area)
	}
}

func (s *S7) ParseAddress() error {
	address := strings.Split(strings.ToUpper(s.Address), ":")
	if len(address) == 2 {
		i, err := strconv.Atoi(address[1])
		if err != nil {
			return err
		}
		s.AreaNo = i
	}
	s.Area = address[0]
	return nil
}

func (s *S7) ReadDB() error {
	c := gos7.NewClient(s.h)
	return c.AGReadDB(s.AreaNo, s.Start, s.Size, s.Bytes)
}
