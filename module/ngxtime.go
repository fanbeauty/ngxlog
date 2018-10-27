package module

import (
	"regexp"
	"fmt"
	"time"
	"ngxlog"
)

var monthMap = map[string]int{"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4, "May": 5, "Jun": 6, "Jul": 7, "Aug": 8, "Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12}

func formatNgxTime(ngxTime string, myHour bool) string {
	ret := regexp.MustCompile(`(\d+)\/(\w+)\/(\d{4}):(\d{2}):(\d{2}):(\d{2})`)
	collection := ret.FindStringSubmatch(ngxTime)
	day := collection[1]
	month := monthMap[collection[2]]
	year := collection[3]
	hour := collection[4]
	minute := collection[5]
	second := collection[6]
	if myHour == true {
		hour = "03"
		minute = "25"
		second = "00"
	}
	timeStr := fmt.Sprintf("%v-%v-%v %v:%v:%v", year, month, day, hour, minute, second)
	return timeStr
}

func TransferNgxTs2UnixTs(entry *ngxlog.Entry) int64 {
	timeLocal, _ := entry.GetField("time_local")
	strTime := formatNgxTime(timeLocal, false)
	p, _ := time.Parse("2006-01-02 15:04:05", strTime)
	return p.Unix()
}

func TransferNgxStartTs2UnixTs(entry *ngxlog.Entry) int64 {
	timeLocal, _ := entry.GetField("time_local")
	strTime := formatNgxTime(timeLocal, true)
	p, _ := time.Parse("2006-01-02 15:04:05", strTime)
	return p.Unix()
}

//获取前一天的时间
func GetYesDate() string {
	nTime := time.Now()
	yesTime := nTime.AddDate(0, 0, -1)
	return fmt.Sprintf("%v-%v-%v", yesTime.Year(), yesTime.Month(), yesTime.Day())
}
