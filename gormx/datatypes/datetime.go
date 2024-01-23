package datatypes

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

type DateTime string

// Scan 实现 sql.Scanner 接口
func (dt *DateTime) Scan(value interface{}) error {
	datetime, ok := value.(time.Time)
	if !ok {
		return errors.New(fmt.Sprint("Failed to time.Time value:", value))
	}
	*dt = DateTime(datetime.Format(time.DateTime))
	return nil
}

// Value 实现 driver.Valuer 接口
func (dt DateTime) Value() (driver.Value, error) {
	return dt, nil
}
