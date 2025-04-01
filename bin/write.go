package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
)

var H = []any{
	"编号", "名称", "耗电量", "最大值", "最大值日期", "最小值", "最小值日期",
}

type Energy struct {
	Code     int
	Name     string
	MinDate  string
	MinValue float64
	MaxDate  string
	MaxValue float64
}

func main() {
	err := Write()
	if err != nil {
		fmt.Println(err)
		fmt.Println("按任意键退出")
		fmt.Scanln()
	}

}

func InputDate() (*time.Time, error) {
	datetime := ""
	fmt.Println("请输入日期时间(格式202503):")
	_, err := fmt.Scanln(&datetime)
	if err != nil {
		return nil, err
	}
	if len(datetime) != 6 {
		return nil, fmt.Errorf("输入日期时间格式错误")
	}
	t, err := time.Parse("200601", datetime)
	if err != nil {
		return nil, fmt.Errorf("输入日期时间格式错误:%s", datetime)
	}
	return &t, nil
}

func Write() error {
	dm, err := InputDate()
	if err != nil {
		return err
	}

	tfp := fmt.Sprintf("%d年各工段设备用电起止明细表.xlsx", dm.Year())

	em, err := GetData(dm)
	if err != nil {
		return fmt.Errorf("读取数据错误：%v", err)
	}

	tmpf, err := excelize.OpenFile(tfp)
	if err != nil {
		return fmt.Errorf("打开文件[%s]错误：%v", tfp, err)
	}
	defer tmpf.Close()

	codes, err := ReadCodes(tmpf)
	if err != nil {
		return err
	}

	sw, err := tmpf.NewStreamWriter("背景数据")
	if err != nil {
		return fmt.Errorf("创建流写入器错误：%v", err)
	}
	err = SetColWith(sw)
	if err != nil {
		return fmt.Errorf("设置列宽错误：%v", err)
	}

	styleID, err := tmpf.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		return fmt.Errorf("创建样式错误：%v", err)
	}
	sw.SetRow("A1", H, excelize.RowOpts{StyleID: styleID})

	for idx, code := range codes[1:] {
		wd := make([]any, 7)
		wd[0] = code
		e, ok := em[code]
		if ok {
			wd[1] = e.Name
			wd[2] = e.MaxValue - e.MinValue
			wd[3] = e.MaxValue
			wd[4] = e.MaxDate
			wd[5] = e.MinValue
			wd[6] = e.MinDate
		}
		cell, err := excelize.CoordinatesToCellName(1, idx+2)
		if err != nil {
			return fmt.Errorf("获取单元格坐标错误：%v", err)
		}
		if err := sw.SetRow(cell, wd, excelize.RowOpts{StyleID: styleID}); err != nil {
			return fmt.Errorf("写入数据错误：%v", err)
		}

	}
	err = sw.Flush()
	if err != nil {
		return fmt.Errorf("流写入错误：%v", err)
	}
	return tmpf.Save()

}

func ReadCodes(f *excelize.File) ([]string, error) {
	cols, err := f.Cols("背景数据")
	if err != nil {
		return nil, fmt.Errorf("读取背景数据表失败：%v", err)
	}
	for cols.Next() {
		return cols.Rows()
	}
	return nil, fmt.Errorf("未读取到第一列数据：%v", err)
}

func GetData(ym *time.Time) (map[string]*Energy, error) {
	dsn := "jldgxcx:Jldg123654.@tcp(xs.jldg.com:3306)/energy?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库错误：%v", err)
	}
	defer db.Close()

	query := `
    SELECT 
        e.code,
        (SELECT name 
         FROM energy 
         WHERE code = e.code 
           AND value = e.max_value 
           AND datetime >= ?
           AND datetime < ?
         ORDER BY datetime DESC 
         LIMIT 1) AS max_name,
        e.max_value,
        (SELECT datetime 
         FROM energy 
         WHERE code = e.code 
           AND value = e.max_value 
           AND datetime >= ?
           AND datetime < ?
         ORDER BY datetime DESC 
         LIMIT 1) AS max_datetime,
        e.min_value,
        (SELECT datetime 
         FROM energy 
         WHERE code = e.code 
           AND value = e.min_value 
           AND datetime >= ?
           AND datetime < ?
         ORDER BY datetime DESC 
         LIMIT 1) AS min_datetime
    FROM (
        SELECT 
            code,
            MAX(value) AS max_value,
            MIN(value) AS min_value
        FROM energy
        WHERE datetime >= ?
          AND datetime < ?
        GROUP BY code
        ORDER BY code
    ) e;`

	bt := ym.Format("2006-01-02")
	et := ym.AddDate(0, 1, 0).Format("2006-01-02")
	rows, err := db.Query(query, bt, et, bt, et, bt, et, bt, et)
	if err != nil {
		return nil, fmt.Errorf("查询执行失败：%v", err)
	}
	defer rows.Close()

	energyMap := make(map[string]*Energy)
	for rows.Next() {
		var (
			code    int
			name    string
			maxVal  float64 // 修复：maxVal 对应 max_value
			maxDate time.Time
			minVal  float64 // 修复：minVal 对应 min_value
			minDate time.Time
		)

		// 关键修复：调整参数顺序
		if err := rows.Scan(
			&code,
			&name,
			&maxVal,
			&maxDate,
			&minVal,
			&minDate,
		); err != nil {
			return nil, fmt.Errorf("数据解析失败：%v", err)
		}

		energyMap[fmt.Sprintf("%d", code)] = &Energy{
			Code:     code,
			Name:     name,
			MinValue: minVal,
			MaxValue: maxVal,
			MinDate:  minDate.Format("2006-01-02 15:04:05"),
			MaxDate:  maxDate.Format("2006-01-02 15:04:05"),
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("结果遍历错误：%v", err)
	}
	return energyMap, nil
}

func SetColWith(sw *excelize.StreamWriter) error {

	if err := sw.SetColWidth(1, 1, 6); err != nil {
		return fmt.Errorf("设置列宽度错误：%v", err)
	}
	if err := sw.SetColWidth(2, 2, 21); err != nil {
		return fmt.Errorf("设置列宽度错误：%v", err)
	}
	if err := sw.SetColWidth(3, 4, 18); err != nil {
		return fmt.Errorf("设置列宽度错误：%v", err)
	}
	if err := sw.SetColWidth(5, 5, 22); err != nil {
		return fmt.Errorf("设置列宽度错误：%v", err)
	}
	if err := sw.SetColWidth(6, 6, 18); err != nil {
		return fmt.Errorf("设置列宽度错误：%v", err)
	}
	if err := sw.SetColWidth(7, 7, 21); err != nil {
		return fmt.Errorf("设置列宽度错误：%v", err)
	}
	return nil
}
