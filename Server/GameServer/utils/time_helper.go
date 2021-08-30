package utils

import (
	"strings"
	"time"
)

const (
	offset = 8*3600  // ----------中国时区-----------
)



//获取 nowTime（时间戳）月开始时间戳
func GetMonthBeginTimeStamp(nowTime int64) int64 {
	thisTime := time.Unix(nowTime, 0).In(time.FixedZone("CST", offset))
	return time.Date(thisTime.Year(), thisTime.Month(), 1, 0, 0, 0, 0, thisTime.Location()).Unix()
}

//获取 nowTime（时间戳）月结束时间戳
func GetMonthEndTimeStamp(nowTime int64) int64 {
	thisTime := time.Unix(nowTime, 0).In(time.FixedZone("CST", offset))
	return time.Date(thisTime.Year(), thisTime.Month()+1, 0, 23, 59, 59, 59, thisTime.Location()).Unix()
}

//获取 nowTime（时间戳）日开始时间戳
func GetDayBeginTimeStamp(nowTime int64) int64 {
	thisTime := time.Unix(nowTime, 0).In(time.FixedZone("CST", offset))
	return time.Date(thisTime.Year(), thisTime.Month(), thisTime.Day(), 0, 0, 0, 0, thisTime.Location()).Unix()
}

//获取 nowTime（时间戳）日结束时间戳
func GetDayEndTimeStamp(nowTime int64) int64 {
	thisTime := time.Unix(nowTime, 0).In(time.FixedZone("CST", offset))
	return time.Date(thisTime.Year(), thisTime.Month(), thisTime.Day(), 23, 59, 59, 59, thisTime.Location()).Unix()
}

//根据输入模板 "Y-m-d H:i:s" or "Y-m-d" 来进行时间转换为时间戳
func TransformationTimeStamp(toBeCharge string, temp string) int64 {
	temp = strings.Replace(temp, " ", "", -1)
	timeLayout := "2006-01-02 15:04:05"

	if temp == "" || temp == "Y-m-d H:i:s" {
		timeLayout = "2006-01-02 15:04:05"
	} else if temp == "Y-m-d" {
		timeLayout = "2006-01-02"
	}

	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, time.FixedZone("CST", offset)) //使用模板在对应时区转化为time.time类型
	return theTime.Unix()                                                                     //转化为时间戳 类型是int64
}

//根据输入模板 "Y-m-d H:i:s" or "Y-m-d" 来进行时间戳转换为时间
func TransformationTime(timeStamp int64, temp string) string {
	temp = strings.Replace(temp, " ", "", -1)
	timeLayout := "2006-01-02 15:04:05"

	if temp == "" || temp == "Y-m-d H:i:s" {
		timeLayout = "2006-01-02 15:04:05"
	} else if temp == "Y-m-d" {
		timeLayout = "2006-01-02"
	}

	thisTime := time.Unix(timeStamp, 0).In(time.FixedZone("CST", offset))
	return thisTime.Format(timeLayout)
}

//获取当前时间戳
func GetNowTimeStamp() int64 {
	return time.Now().In(time.FixedZone("CST", offset)).Unix()
}

//获取当前时间 （根据模板 "Y-m-d H:i:s" or "Y-m-d"）
func GetNowTime(temp string) string {
	timeStamp := time.Now().In(time.FixedZone("CST", offset)).Unix()
	temp = strings.Replace(temp, " ", "", -1)
	timeLayout := "2006-01-02 15:04:05"

	if temp == "" || temp == "Y-m-d H:i:s" {
		timeLayout = "2006-01-02 15:04:05"
	} else if temp == "Y-m-d" {
		timeLayout = "2006-01-02"
	}
	thisTime := time.Unix(timeStamp, 0).In(time.FixedZone("CST", offset))
	return thisTime.Format(timeLayout)
}

//获取时间戳所在年
func GetYear(timeStamp int64) int {
	thisTime := time.Unix(timeStamp, 0).In(time.FixedZone("CST", offset))
	return thisTime.Year()
}

//获取时间戳所在月
func GetMonth(timeStamp int64) int {
	thisTime := time.Unix(timeStamp, 0).In(time.FixedZone("CST", offset))
	return int(thisTime.Month())
}

//获取时间戳所在日
func GetDay(timeStamp int64) int {
	thisTime := time.Unix(timeStamp, 0).In(time.FixedZone("CST", offset))
	return thisTime.Day()
}
