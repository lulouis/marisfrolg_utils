package marisfrolg_utils

import (
	"fmt"
	"time"
)

/*
时间相关
*/
const (
	LongDateFormat  = "2006-01-02 15:04:05"
	ShortDateFormat = "2006-01-02"
)

//获取日期格式
func GetLongDateString(date string, Hours int64) (dateString string, err error) {
	if len(date) <= 0 {
		return "", fmt.Errorf("时间不能为空")
	}
	inputDate, err := time.Parse(LongDateFormat, date)
	if err == nil {
		h, _ := time.ParseDuration("1h")
		d := inputDate.Add(time.Duration(Hours) * h)
		return d.Format(LongDateFormat), err
	} else {
		return "", fmt.Errorf("时间格式错误")
	}
}

//获取日期格式
func GetShortDateString(date string, Hours int64) (dateString string, err error) {
	if len(date) <= 0 {
		return "", fmt.Errorf("时间为空")
	}
	inputDate, err := time.Parse(ShortDateFormat, date)
	if err == nil {
		h, _ := time.ParseDuration("1h")
		d := inputDate.Add(time.Duration(Hours) * h)
		return d.Format(LongDateFormat), err
	} else {
		return "", fmt.Errorf("时间格式错误")
	}
}

//获取相差时间
func GetMinuteDiffer(start_time, end_time string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 60
		return hour
	} else {
		return hour
	}
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

//获取某一天的23:59:59点时间
func GetZeroTimeEnd(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 59, d.Location())
}

