package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dduutt/go-energy/meter"
)

func main() {
	path := "数据采集配置.xlsx"
	ms, err := meter.FromExcel(path)
	if err != nil {
		log.Fatalf("[excel error]打开excel失败:%v\n", err)
		return
	}
	_, err = meter.InitDB()
	if err != nil {
		log.Fatalf("[db error]打开数据库失败:%v\n", err)
		return
	}
	g := meter.GroupByAddrFromExcel(ms)
	now := time.Now()

	rc := SyncRead(g)

	for mag := range rc {
		if err := *mag.Error; err != "" {
			log.Printf("[group error]%s:%v\n", mag.Addr, err)
			continue
		}
		for _, m := range mag.Meters {
			if err := *m.Error; err != "" {
				log.Printf("[meter error]%d:%v\n", m.Code, err)
				continue
			}
			err := m.ParseFloat()
			if err != nil {
				log.Printf("[parse error]%d:%v\n", m.Code, err)
				continue
			}
			// 写入数据库
			err = meter.InsertEnergy(m)
			if err != nil {
				log.Printf("[insert err]%d:%v\n", m.Code, err)
			}
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
