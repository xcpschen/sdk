package notice

import "strings"

type IContentType interface {
	GetTypeKey() string // 获取类型
}

func StrPad(str string, count int) string {
	arr := make([]string, 0)
	for i := 0; i < count; i++ {
		arr = append(arr, str)
	}
	return strings.Join(arr, "")
}

func InStrArray(a string, b []string) bool {
	for _, v := range b {
		if v == a {
			return true
		}
	}
	return false
}
