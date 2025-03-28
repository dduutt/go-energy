package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/dduutt/go-energy/meter"
)

func main() {
	path := "数据采集配置.xlsx"
	ms, err := meter.FromExcel(path)
	if err != nil {
		panic(err)
	}
	g := meter.GroupByAddrFromExcel(ms)
	now := time.Now()
	rc := SyncRead(g)

	for mag := range rc {
		if mag.Error != nil {
			fmt.Println(mag.Error, mag.Addr)
			continue
		}
		for _, m := range mag.Meters {
			if err := m.Error; err != nil {
				fmt.Println(m.Code, err)
				continue
			}
			err := m.ParseFloat()
			if err != nil {
				fmt.Println(m.Code, err)
				continue
			}
			fmt.Println(m.Code, m.Bytes, math.Round(m.Value*100)/1000)
		}

	}
	fmt.Println(time.Since(now))
}

func SyncRead(m map[string]*meter.AddrGroup) chan *meter.AddrGroup {
	rc := make(chan *meter.AddrGroup, 1)
	var wg sync.WaitGroup
	go func() {
		wg.Wait()
		close(rc)
	}()
	for _, g := range m {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.Read()
			rc <- g
		}()
	}
	return rc
}
