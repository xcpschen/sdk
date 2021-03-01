package dt_test

import (
	"fmt"
	"testing"
)

func TestSignature(t *testing.T) {
	// g := &dt.QrLogin{
	// 	TimeStamp: int64(1546084445901),
	// 	AppSerect: "testappSecret",
	// }
	// g.DoSignature()
	// fmt.Println(g.Signature)
	// u := url.Values{}
	// u.Add("signature", g.Signature)
	// fmt.Println(u.Encode())
	// now := time.Now().AddDate(0, 0, -1)
	// fmt.Println(now)
	// tmp := now.Format("2006-01-02 15:04:05")
	// fmt.Println(tmp)
	// fmt.Println(now)

	// fmt.Println(float64(14) / float64(100))
	// fmt.Println(int64(math.Ceil(float64(14) / float64(100))))

	fmt.Println(Pro3([]int{1, 2, 3, -4, 0, 0, -6, 9, 1}))
}

func Pro3(list []int) []int {
	l := len(list)
	if l <= 0 {
		return nil
	}
	left := 0
	right := l - 1
	mid := 0
	tmp := 0
	for {
		if mid > right {
			break
		}
		if list[mid] > 0 {
			tmp = list[left]
			list[left] = list[mid]
			list[mid] = tmp
			left++
			mid++
		} else {
			if list[mid] == 0 {
				mid++
			} else {
				tmp = list[mid]
				list[mid] = list[right]
				list[right] = tmp
				right--
			}
		}

	}
	return list
}

func swp(a *int, b *int) {
	tmp := 0
	tmp = *a
	*a = *b
	*b = tmp
}
