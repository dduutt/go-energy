package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/dduutt/go-energy/meter"
)

func main() {
	f, logger := InitLogger()
	defer f.Close()

	path := "数据采集配置.xlsx"
	ms, err := meter.FromExcel(path)
	if err != nil {
		logger.Fatalf("[excel error]打开excel失败:%v\n", err)
		return
	}
	_, err = meter.InitDB()
	if err != nil {
		logger.Fatalf("[db error]打开数据库失败:%v\n", err)
		return
	}
	g := meter.GroupByAddrFromExcel(ms)
	now := time.Now()

	rc := SyncRead(g)

	for mag := range rc {
		if err := *mag.Error; err != "" {
			logger.Printf("[group error]%s:%v\n", mag.Addr, err)
			continue
		}
		for _, m := range mag.Meters {
			if err := *m.Error; err != "" {
				logger.Printf("[meter error]%d %s:%v\n", m.Code, m.Name, err)
				continue
			}
			err := m.ParseFloat()
			if err != nil {
				logger.Printf("[parse error]%d %s:%v\n", m.Code, m.Name, err)
				continue
			}
			// 写入数据库
			err = meter.InsertEnergy(m)
			if err != nil {
				logger.Printf("[insert err]%d %s:%v\n", m.Code, m.Name, err)
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

func InitLogger() (*os.File, *log.Logger) {

	date := time.Now().Format("2006-01-02")
	name := date + ".txt"
	fp := path.Join("logs", name)
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		log.Fatalln("创建日志文件夹失败")
	}
	f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("[log error]打开日志文件失败:%v\n", err)
	}
	return f, log.New(f, "", log.Ldate|log.Ltime)
}
