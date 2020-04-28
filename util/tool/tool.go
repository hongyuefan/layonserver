package tool

import (
	"fmt"
	"time"
)

//获取上个月一号到昨天的时间
func GetLastMonthBeginEndUnixTime() (int64, int64, int, string, string) {
	day := int(time.Now().Day())
	month := int(time.Now().Month())
	year := int(time.Now().Year())
	if month == 1 {
		month = 12
		year--
	} else {
		month--
	}
	var monthS string
	if month < 10 {
		monthS = fmt.Sprintf("0%v", month)
	} else {
		monthS = fmt.Sprintf("%v", month)
	}

	var dayS string
	if day < 10 {
		dayS = fmt.Sprintf("0%v", day)
	} else {
		dayS = fmt.Sprintf("%v", day)
	}

	beginS := fmt.Sprintf("%v-%v-01", year, monthS)
	endS := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	beginT, _ := time.Parse("2006-01-02", beginS)
	endT, _ := time.Parse("2006-01-02", endS)
	return beginT.Unix(), endT.Unix(), year, monthS, dayS
}

//获取当月1号到昨天的unix时间
func GetMonthBeginEndUnixTime() (int64, int64, int, string, string) {
	day := int(time.Now().Day())
	month := int(time.Now().Month())
	year := int(time.Now().Year())

	var monthS string
	if month < 10 {
		monthS = fmt.Sprintf("0%v", month)
	} else {
		monthS = fmt.Sprintf("%v", month)
	}
	var dayS string
	if day < 10 {
		dayS = fmt.Sprintf("0%v", day)
	} else {
		dayS = fmt.Sprintf("%v", day)
	}

	beginS := fmt.Sprintf("%v-%v-01", year, monthS)

	endS := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	beginT, _ := time.Parse("2006-01-02", beginS)

	endT, _ := time.Parse("2006-01-02", endS)

	return beginT.Unix(), endT.Unix(), year, monthS, dayS
}
