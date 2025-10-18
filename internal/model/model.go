package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Book{}, &Chapter{}, &ReadingRecord{})
}

// JSONTime 是一个包装 time.Time 的自定义类型
type JSONTime time.Time

// Value 数据库写入时调用
func (t JSONTime) Value() (driver.Value, error) {
	return time.Time(t), nil // 返回 time.Time 类型，GORM 就可以识别
}

// 数据库读取时调用
func (t *JSONTime) Scan(v interface{}) error {
	if v == nil {
		*t = JSONTime(time.Time{})
		return nil
	}
	switch val := v.(type) {
	case time.Time:
		*t = JSONTime(val)
		return nil
	case []byte:
		tt, err := time.Parse("2006-01-02 15:04:05", string(val))
		if err != nil {
			return err
		}
		*t = JSONTime(tt)
		return nil
	case string:
		tt, err := time.Parse("2006-01-02 15:04:05", val)
		if err != nil {
			return err
		}
		*t = JSONTime(tt)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into JSONTime", v)
	}
}

// JSON 序列化
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// JSON 反序列化
func (t *JSONTime) UnmarshalJSON(b []byte) error {
	parse, err := time.Parse(`"2006-01-02 15:04:05"`, string(b))
	if err != nil {
		return err
	}
	*t = JSONTime(parse)
	return nil
}
