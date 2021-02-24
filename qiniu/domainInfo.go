package qiniu

// DomainInfo 接口域名信息
type DomainInfo struct {
	Domain string `json:"domain"`

	// 存储空间名字
	Tbl string `json:"tbl"`
	// 用户UID
	Owner int `json:"uid"`
	Ctime int `json:"ctime"`
	Utime int `json:"utime"`
}
