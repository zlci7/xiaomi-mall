package parseTime

import (
	"time"
)

// ParseDateStr 把"2006-01-02"格式的字符串转成time.Time（东八区），解析失败返回错误
// dateStr：待解析的日期字符串
func ParseDateStr(dateStr string) (time.Time, error) {
	layout := "2006-01-02"
	// 基础解析
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}

	// 转换为东八区（Asia/Shanghai），可选（根据业务需求）
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t = t.In(loc)

	return t, nil
}

// ParseDateTimeStr 把"2006-01-02 15:04:05"格式的字符串转成time.Time（东八区），解析失败返回错误
// dateTimeStr：待解析的日期时间字符串
func ParseDateTimeStr(dateTimeStr string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	// 基础解析
	t, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		return time.Time{}, err
	}

	// 转换为东八区（Asia/Shanghai）
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t = t.In(loc)

	return t, nil
}
