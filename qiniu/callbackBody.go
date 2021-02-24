package qiniu

import (
	"fmt"
	"strings"
)

// MagicVar qiniu magic var for callback
type MagicVar []string

// Set var
func (m *MagicVar) Set(key string) {
	*m = append(*m, fmt.Sprintf("%s=$(%s)", key, key))
}

// ToString to string
func (m *MagicVar) ToString() string {
	return strings.Join(*m, "&")
}

// CustomeVar qiniu customer var for callback
type CustomeVar []string

// Set var
func (c *CustomeVar) Set(key string) {
	*c = append(*c, fmt.Sprintf("%s=$(x:%s)", key, key))
}

// ToString To String
func (c *CustomeVar) ToString() string {
	return strings.Join(*c, "&")
}
