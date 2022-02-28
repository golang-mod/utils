package version

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
)

func Check(minVersion string) (err error) {
	var _goVersion string = runtime.Version()

	var goVersion string = strings.Replace(_goVersion, "go", "", -1)
	var goArray []string = strings.Split(goVersion, ".")

	var go1 int
	var go2 int
	var goYes int
	for _, value := range goArray {
		theValue, _ := strconv.ParseInt(value, 10, 64)
		if go1 == 0 {
			go1 = int(theValue)
		} else {
			if go2 == 0 {
				go2 = int(theValue)
				if go1 == 1 {
					goYes = goYes + 1
					if go2 >= 15 {
						goYes = goYes + 1
					}
				} else if go1 >= 2 {
					goYes = goYes + 2
				}
			} else {
				break
			}
		}
	}
	if goYes != 2 {
		return errors.New("当前Go版本：" + _goVersion + " 最低要求版本：" + minVersion)
	}
	return nil
}
