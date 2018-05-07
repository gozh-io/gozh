package common

import (
	"fmt"
	"time"
)

// TimeIn returns the time in UTC if the name is "" or "UTC".
// It returns the local time if the name is "Local".
// Otherwise, the name is taken to be a location name in
// the IANA Time Zone database, such as "Africa/Lagos".
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

//获取到当前日期和时间,按照北京时间计算
func GetBeiJingDT() string {
	t, _ := TimeIn(time.Now(), "Asia/Shanghai")
	return fmt.Sprintf("%s", t.Format("2006-01-02 15:04:05"))
}
