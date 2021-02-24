package ctrlvideo

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

var (
	_rand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// InArr check elment is in arrry b
func InArr(a string, b []string) bool {
	for _, j := range b {
		if a == j {
			return true
		}
	}
	return false
}

// RandConstString 随机字符串
func RandConstString(length int, cstr string) string {
	s := ""
	max := len(cstr)
	for i := 0; i < length; i++ {
		idx := _rand.Intn(max)
		s = s + cstr[idx:idx+1]
	}
	return s
}

// RandNumStr 随机生成数字字符
func RandNumStr(length int) string {
	return RandConstString(length, "0123456789")
}

// RandString 随机生成字符
func RandString(length int) string {
	return RandConstString(length, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
}

// RandStrWithUppAndNum s
func RandStrWithUppAndNum(length int) string {
	return RandConstString(length, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// RandSignStr 随机生成字符/符号
func RandSignStr(length int) string {
	return RandConstString(length, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz~!@#$%^&*()_+`-={}[]:;|?,.")
}

// StructToParam struct kv add into url.Values
func StructToParam(obj interface{}, data *Param) error {
	var val reflect.Value
	reflectRawValue := reflect.ValueOf(obj)
	if reflectRawValue.Kind() == reflect.Ptr {
		val = reflectRawValue.Elem()
		if val.Kind() != reflect.Struct {
			return fmt.Errorf("解析结构体错误1")
		}
	} else if reflectRawValue.Kind() == reflect.Struct {
		var ok bool
		val, ok = obj.(reflect.Value)
		if !ok {
			return fmt.Errorf("解析结构体错误2")
		}
		if !val.CanAddr() {
			val = val.Elem()
		}
	} else {
		return fmt.Errorf("解析结构体错误3")
	}
	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i)
		if !value.CanInterface() {
			continue
		}

		kind := value.Kind()

		typeFile := val.Type().Field(i)
		tag := typeFile.Tag.Get("url")
		if tag == "" {
			tag = typeFile.Name
		} else {
			tmp := strings.Split(tag, ",")
			if InArr("--", tmp) {
				continue
			}
			tag = tmp[0]
		}
		if kind == reflect.Struct {
			if tag == "inherit" {
				StructToParam(value, data)
			}
			continue
		} else {
			if kind == reflect.String && value.Len() == 0 {
				continue
			}
			data.Set(tag, fmt.Sprintf("%v", value.Interface()))
		}
	}
	return nil
}
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
