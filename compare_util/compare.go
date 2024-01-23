package compare_util

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Compare(newObj, oldObj, filter interface{}, name, field string) (list [][]string, err error) {
	if reflect.TypeOf(newObj).Kind() != reflect.Ptr || reflect.TypeOf(oldObj).Kind() != reflect.Ptr {
		return list, errors.New("不是指针类型，没法进行操作")
	}
	if reflect.TypeOf(filter).Kind() != reflect.Struct {
		return list, errors.New("filter不是Struct类型，没法进行操作")
	}

	vn := reflect.ValueOf(newObj).Elem()
	vo := reflect.ValueOf(oldObj).Elem()
	flt := reflect.TypeOf(filter)
	for i := 0; i < flt.NumField(); i++ {
		key := flt.Field(i)
		nField := vn.FieldByName(key.Name)
		oField := vo.FieldByName(key.Name)
		if !nField.IsValid() || !oField.IsValid() {
			continue
		}
		nif := nField.Interface()
		oif := oField.Interface()
		var nStr string
		var oStr string
		var fieldType string
		nStr, fieldType, err = toString(nif)
		oStr, fieldType, err = toString(oif)

		if nStr != oStr {
			list = append(list, []string{
				key.Tag.Get(name),
				key.Tag.Get(field),
				nStr,
				oStr,
				fieldType,
			})
		}
	}
	return
}

func toString(nif interface{}) (nStr, fieldType string, err error) {
	switch nif.(type) {
	case float64:
		fieldType = "float64"
		nStr = strconv.FormatFloat(nif.(float64), 'f', -1, 64)
	case *float64:
		fieldType = "*float64"
		if nif.(*float64) != nil {
			nStr = strconv.FormatFloat(*nif.(*float64), 'f', -1, 64)
		}
	case float32:
		fieldType = "float32"
		nStr = strconv.FormatFloat(float64(nif.(float32)), 'f', -1, 64)
	case *float32:
		fieldType = "*float32"
		nStr = strconv.FormatFloat(float64(*nif.(*float32)), 'f', -1, 64)
	case int:
		fieldType = "int"
		nStr = strconv.Itoa(nif.(int))
	case uint:
		fieldType = "uint"
		nStr = strconv.Itoa(int(nif.(uint)))
	case int8:
		fieldType = "int8"
		nStr = strconv.Itoa(int(nif.(int8)))
	case uint8:
		fieldType = "uint8"
		nStr = strconv.Itoa(int(nif.(uint8)))
	case int16:
		fieldType = "int16"
		nStr = strconv.Itoa(int(nif.(int16)))
	case uint16:
		fieldType = "uint16"
		nStr = strconv.Itoa(int(nif.(uint16)))
	case int32:
		fieldType = "int32"
		nStr = strconv.Itoa(int(nif.(int32)))
	case uint32:
		fieldType = "uint32"
		nStr = strconv.Itoa(int(nif.(uint32)))
	case int64:
		fieldType = "uint32"
		nStr = strconv.FormatInt(nif.(int64), 10)
	case uint64:
		fieldType = "uint64"
		nStr = strconv.FormatUint(nif.(uint64), 10)
	case string:
		fieldType = "string"
		nStr = nif.(string)
	case []string:
		fieldType = "[]string"
		nStr = strings.Join(nif.([]string), "")
	case *int:
		fieldType = "*int"
		if nif.(*int) != nil {
			nStr = strconv.Itoa(*nif.(*int))
		}
	case *int8:
		fieldType = "*int8"
		if nif.(*int8) != nil {
			nStr = strconv.Itoa(int(*nif.(*int8)))
		}
	case *int64:
		fieldType = "*int64"
		if nif.(*int64) != nil {
			nStr = strconv.Itoa(int(*nif.(*int64)))
		}
	case *time.Time:
		fieldType = "*time.Time"
		if nif.(*time.Time) != nil {
			nStr = nif.(*time.Time).Format(time.DateTime)
		}
	case time.Time:
		fieldType = "time.Time"
		if !nif.(time.Time).IsZero() {
			nStr = nif.(time.Time).Format(time.DateTime)
		}
	default:
		err = errors.New("未知类型")
		return
	}
	return
}
