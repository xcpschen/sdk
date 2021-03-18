package ldap

import (
	"fmt"
	"strings"
)

type Filter interface {
	SetVal(interface{})
	ToString() string
}

type FilterOperation struct {
	opt string
	Val []interface{}
}

type FOVal struct {
	Opt string
	Key string
	Val string
}

func NewAddFilterOperation() Filter {
	return &FilterOperation{opt: "&"}
}
func NewORFilterOperation() Filter {
	return &FilterOperation{opt: "|"}
}

func (fo *FilterOperation) ToString() string {
	val := []string{}
	for _, item := range fo.Val {
		switch v := item.(type) {
		case string:
			val = append(val, v)
		case map[string]string:
			for k, m := range v {
				val = append(val, fmt.Sprintf("(%s=%s)", k, m))
			}
		case Filter:
			val = append(val, v.ToString())
		}
	}

	return fmt.Sprintf("(%s%s)", fo.opt, strings.Join(val, ""))
}

func (fo *FilterOperation) SetVal(Val interface{}) {
	fo.Val = append(fo.Val, Val)
}

func (fo *FilterOperation) SetEquality(key, val string) {
	fo.Val = append(fo.Val, fmt.Sprintf("(%s=%s)", key, val))
}

func (fo *FilterOperation) SetNegation(key, val string) {
	fo.Val = append(fo.Val, fmt.Sprintf("(!(%s=%s))", key, val))
}

func (fo *FilterOperation) SetPresence(key string) {
	fo.Val = append(fo.Val, fmt.Sprintf("(%s=*)", key))
}

func (fo *FilterOperation) SetAbsence(key string) {
	fo.Val = append(fo.Val, fmt.Sprintf("(!(%s=*))", key))
}

func (fo *FilterOperation) SetGreaterThan(key, val string) {
	fo.Val = append(fo.Val, fmt.Sprintf("(%s>=%s)", key, val))
}

func (fo *FilterOperation) SetLessThan(key, val string) {
	fo.Val = append(fo.Val, fmt.Sprintf("(%s<=%s)", key, val))
}

func (fo *FilterOperation) SetProximity(key, val string) {
	fo.Val = append(fo.Val, fmt.Sprintf("(%s~=%s)", key, val))
}
