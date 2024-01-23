package datatypes

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

type Args []string

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *Args) Scan(value interface{}) error {
	str, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to string value:", value))
	}
	if string(str) == "" {
		*j = []string{}
		return nil
	}
	*j = strings.Split(string(str), ",")
	return nil
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j Args) Value() (driver.Value, error) {
	if len(j) > 0 {
		var str = j[0]
		for _, v := range j[1:] {
			str += "," + v
		}
		return str, nil
	}
	return "", nil
}
