package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/dduutt/go-energy/client"
	"github.com/dduutt/go-energy/meter"
)

func main() {
	path := "数据采集配置.xlsx"
	ms, err := meter.FromExcel(path)
	if err != nil {
		panic(err)
	}
	cms := meter.CodeMap(ms)
	g := meter.GroupByAddrFromExcel(ms)
	now := time.Now()
	rc := SyncRead(g)

	for rg := range rc {
		if rg.Error != nil {
			fmt.Println(rg.Error)
			continue
		}
		for _, r := range rg.Results {
			if r.Error != nil {
				fmt.Println(r.Error)
				continue
			}
			m := cms[r.ID]
			fmt.Println(m, r.Result)
			err := m.ParseFloat()
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(m.Value)

		}

	}
	fmt.Println(time.Since(now))
}

func SyncRead(m map[string]*meter.AddrGroup) chan *client.GroupSyncReadResult {
	rc := make(chan *client.GroupSyncReadResult, 1)
	var wg sync.WaitGroup
	go func() {
		wg.Wait()
		close(rc)
	}()
	for _, g := range m {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.Read(rc)

		}()
	}
	return rc
}
