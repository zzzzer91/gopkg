package timex

import "time"

var (
	CstZone = time.FixedZone("CST", int((8 * time.Hour).Seconds()))
)

func RangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, start.Location())
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, end.Location())

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

// GetDateCst 获取 cst 时区今天时间为 0 的 Time
func GetDateCst() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, CstZone)
}

// ParseDateInCst 解析 `2022-08-10` 这样的字符串为 cst 时区
func ParseDateInCst(date string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", date, CstZone)
}

func GetMondayOfWeek(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	return t.AddDate(0, 0, offset)
}

func GetMondayOfNextWeek(t time.Time) time.Time {
	return GetMondayOfWeek(t).AddDate(0, 0, 7)
}

// 计算日期相差多少月
func SubMonth(t1, t2 time.Time) (month int) {
	y1 := t1.Year()
	y2 := t2.Year()
	m1 := int(t1.Month())
	m2 := int(t2.Month())
	if y1 > y2 {
		if m1 >= m2 {
			return (y1-y2)*12 + m1 - m2
		}
		return (y1-y2-1)*12 + (12 + m1 - m2)
	} else {
		return m1 - m2
	}
}

// 获取某个时间下个月第一天
func GetFirstDayOfNextMonth(t time.Time) time.Time {
	t = t.AddDate(0, 1, 0)
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, CstZone)
}
