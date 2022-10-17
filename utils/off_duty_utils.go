package utils

import (
	"github.com/6tail/lunar-go/HolidayUtil"
	"github.com/6tail/lunar-go/calendar"
	"time"
)

// 节假日相关工具类
type offDuty struct{}

// OffDuty
// @description: 获取节假日相关工具类实例
// @return *offDuty
func OffDuty() *offDuty {
	return &offDuty{}
}

// CheckIsHoliday 检查是不是节假日，如果是，返回是什么节假日
func (od offDuty) CheckIsHoliday(t time.Time) (bool, string) {
	lunar := calendar.NewSolarFromYmd(t.Year(), int(t.Month()), t.Day())
	d := HolidayUtil.GetHoliday(lunar.ToYmd())
	// 不是节假日，判断是不是周末
	if d == nil {
		// 如果不是工作日，返回星期几
		if lunar.GetWeek() == 0 || lunar.GetWeek() == 6 {
			return true, "星期" + lunar.GetWeekInChinese()
		}
		return false, ""
	}
	return true, d.GetName()
}

// GetNextHolidayOrWeekend 获取最近的节假日或者周末，优先返回节假日
func (od offDuty) GetNextHolidayOrWeekend() (string, int) {
	// 获取最近的节假日
	h, hd := od.GetNextHoliday()
	// 如果节假日是当前日期，则返回
	if hd == 0 {
		return h, hd
	}
	// 获取最近的周末
	w, wd := od.GetNextWeekend()
	// 如果周末是当前日期，则返回
	if wd == 0 {
		return w, wd
	}
	// 如果节假日比周末早，则返回节假日
	if hd < wd {
		return h, hd
	}
	// 返回周末
	return w, wd
}

// GetNextHoliday 获取下一个节日并返回距离多少天
func (od offDuty) GetNextHoliday() (string, int) {
	t := time.Now()
	thisYearEndDay := time.Date(t.Year(), 12, 31, 0, 0, 0, 0, time.Local)
	endDay := int(thisYearEndDay.Sub(t).Hours() / 24)

	name := ""
	i := 0
	lunar := calendar.NewSolarFromYmd(t.Year(), int(t.Month()), t.Day())
	for {
		if i >= endDay {
			i = -1
			break
		}
		h := HolidayUtil.GetHoliday(lunar.ToYmd())
		if h == nil {
			lunar = lunar.Next(1)
			i++
			continue
		}
		if !h.IsWork() {
			name = h.GetName()
			break
		}
		lunar = lunar.Next(1)
		i++
	}
	return name, i
}

// GetNextWeekend 获取下一个周末并返回距离多少天
func (od offDuty) GetNextWeekend() (string, int) {
	t := time.Now()
	name := ""
	i := 0
	lunar := calendar.NewSolarFromYmd(t.Year(), int(t.Month()), t.Day())
	for {
		h := HolidayUtil.GetHoliday(lunar.ToYmd())
		if h == nil {
			if lunar.GetWeek() == 0 || lunar.GetWeek() == 6 {
				name = lunar.GetWeekInChinese()
				break
			}
			lunar = lunar.Next(1)
			i++
			continue
		}
		if !h.IsWork() {
			name = h.GetName()
			break
		}
		lunar = lunar.Next(1)
		i++
	}
	if len(name) != 1 {
		return name, i
	}
	return "星期" + name, i
}
