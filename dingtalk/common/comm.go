package common

import (
	"encoding/json"
	"math/rand"
	"time"
)

func JsonEncode(data interface{}) string {
	s, err := json.Marshal(data)
	if err == nil {
		return string(s)
	}
	return ""
}
func JsonEncodeByte(data interface{}) []byte {
	buf, err := json.Marshal(data)
	if err == nil {
		return buf
	}
	return nil
}

func JsonDecode(s string, data interface{}) error {
	return json.Unmarshal([]byte(s), data)
}
func JsonDecodeByte(s []byte, data interface{}) error {
	return json.Unmarshal([]byte(s), data)
}

// 随机字符串
func RandConstString(length int, cstr string) string {
	s := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	max := len(cstr)
	for i := 0; i < length; i++ {
		idx := r.Intn(max)
		s = s + cstr[idx:idx+1]
	}
	return s
}

// 随机生成数字字符
func RandNumStr(length int) string {
	return RandConstString(length, "0123456789")
}

// 随机生成字符
func RandString(length int) string {
	return RandConstString(length, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
}

func RandStrWithUppAndNum(length int) string {
	return RandConstString(length, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// 随机生成字符/符号
func RandSignStr(length int) string {
	return RandConstString(length, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz~!@#$%^&*()_+`-={}[]:;|?,.")
}
