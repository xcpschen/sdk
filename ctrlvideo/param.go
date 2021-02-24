package ctrlvideo

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Param hanlder request param
type Param struct {
	data map[string]string
}

// Set kv key will be turn lower and key no have "X-Qiniu-"
func (p *Param) Set(key string, val string) {
	if key == "" {
		return
	}
	if p.data == nil {
		p.data = map[string]string{}
	}
	key = strings.ToLower(key)
	p.data[key] = val
}

// ToSignURLString turn to sign string
func (p *Param) ToSignURLString() string {
	keys := []string{}
	for k, _ := range p.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	str := []string{}
	for _, k := range keys {
		str = append(str, fmt.Sprintf("%s=%s", k, p.data[k]))
	}
	return strings.Join(str, "&")
}

// ToUrlEncodeString return urlencode string
func (p *Param) ToUrlEncodeString() string {
	u := url.Values{}
	for k, v := range p.data {
		u.Add(k, v)
	}
	return u.Encode()
}
