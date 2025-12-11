package timer

import "time"

// 便携的time工具

func GetNowTime() time.Time {
	return time.Now()
}

// 时间推算
func GetCalculateTime(currentTimer time.Time, d string) (time.Time, error) {
	duration, err := time.ParseDuration(d)
	// 在字符串中解析出持续时间
	if err != nil {
		return time.Time{}, err
	}
	return currentTimer.Add(duration), nil
}


