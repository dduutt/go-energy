package meter

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

var DB_CONF = make(map[string]string)
var H = []string{"编号", "车间", "配电室", "名称", "协议", "IP", "PORT", "从站/区域", "地址", "长度", "数据类型", "倍率", "字节序"}

func FromExcel(file string) ([]*Energy, error) {
	f, err := excelize.OpenFile(file)

	if err != nil {
		return nil, err
	}
	defer f.Close()
	crows, err := f.GetRows("配置表")
	if err != nil {
		return nil, fmt.Errorf("读取配置表错误:%v", err)
	}
	for _, row := range crows {
		DB_CONF[row[0]] = row[1]
	}

	rows, err := f.GetRows("设备表")

	if err != nil {
		return nil, err
	}

	header := rows[0]
	if !compareHeader(header) {
		return nil, fmt.Errorf("表头应为:%v读取到:\n%v", H, header)
	}
	var meters []*Energy
	for i, row := range rows[1:] {
		l, err := strconv.Atoi(row[9])
		if err != nil {
			return nil, fmt.Errorf("第%d行长度错误:%v", i+2, row[9])
		}
		p, err := strconv.Atoi(row[6])
		if err != nil {
			return nil, fmt.Errorf("第%d行端口错误:%v", i+2, row[6])
		}
		address, err := strconv.Atoi(row[8])
		if err != nil {
			return nil, fmt.Errorf("第%d行地址错误:%v", i+2, row[8])
		}
		magnifications, err := strconv.ParseFloat(row[11], 64)
		if err != nil {
			return nil, fmt.Errorf("第%d行倍数错误:%v", i+2, row[11])
		}
		code, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("第%d行编号错误:%v", i+2, row[0])
		}
		e := ""
		m := &Energy{
			Code:          code,
			WorkShop:      row[1],
			Room:          row[2],
			Name:          row[3],
			Protocol:      row[4],
			IP:            row[5],
			Port:          p,
			SlaveOrArea:   row[7],
			Start:         address,
			Size:          l,
			DataType:      row[10],
			Magnification: magnifications,
			Value:         0,
			Bytes:         make([]byte, l),
			ByteOrder:     row[12],
			Error:         &e,
		}
		meters = append(meters, m)

	}
	return meters, nil

}
