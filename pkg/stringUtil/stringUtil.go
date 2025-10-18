package stringUtil

import (
	"errors"
	"strconv"
)

func StringToUint64(str string) (uint64, error) {
	if str == "" {
		return 0, errors.New("空字符不能被转换")
	}
	// 将字符串转换为uint
	id64, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, errors.New("传入的字符必须是数字")
	}
	return id64, nil
}
