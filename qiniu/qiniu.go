package qiniu

// 域名和响应代码
const (
	UPLOAD_HOST_HUADONG = "up.qiniup.com"     // 华东
	UPLOAD_HOST_HUANAN  = "up-z2.qiniup.com"  // 华南
	UPLOAD_HOST_HUABEI  = "up-z1.qiniup.com"  // 华北
	UPLOAD_HOST_NA      = "up-na0.qiniup.com" // 北美
	// UPLOAD_HOST_SA      = "up-as0.qiniup.com" // 东南亚
	UPLOAD_HOST_SA = "upload-as0.qiniup.com" // 东南亚

	RS_MG_HOST  = "rs.qiniu.com"  // 资源管理
	RS_API_HOST = "api.qiniu.com" // API
	RSF_MG_HOST = "rsf.qbox.me"   //
	RS_BOX_HOST = "rs.qbox.me"

	STATUS_SUCCESS      = 200 // 成功
	STATUS_FAILD        = 298 // 部分/所有操作失败
	STATUS_QUERY_FORMAT = 400 // 请求报文格式错误
	STATUS_TOKEN_ERROR  = 401 // TOKEN 错误
	STATUS_SERVER_ERROR = 599 // 服务器操作失败
	STATUS_NOT_EXISTS   = 612 // 文件资源不存在

	FILE_STATUS_ENABLE  = 0 // 启用文件
	FILE_STATUS_DISABLE = 1 // 禁用文件
)

var (
	status = map[int]string{
		STATUS_SUCCESS:      "成功",
		STATUS_FAILD:        "部分/所有操作失败",
		STATUS_QUERY_FORMAT: "请求报文格式错误",
		STATUS_TOKEN_ERROR:  "TOKEN 错误",
		STATUS_SERVER_ERROR: "服务器操作失败",
		STATUS_NOT_EXISTS:   "文件资源不存在",
	}
	hosts = map[string]string{
		UPLOAD_HOST_HUADONG: "华东",
		UPLOAD_HOST_HUANAN:  "华南",
		UPLOAD_HOST_HUABEI:  "华北",
		UPLOAD_HOST_NA:      "北美",
		UPLOAD_HOST_SA:      "东南亚",
	}
	upload_host = map[string]string{
		"hd": UPLOAD_HOST_HUADONG,
		"hn": UPLOAD_HOST_HUANAN,
		"hb": UPLOAD_HOST_HUABEI,
		"na": UPLOAD_HOST_NA,
		"sa": UPLOAD_HOST_NA,
	}

	ApiHostMap = map[string]string{
		"hd": "api-z0.qiniu.com",
		"hb": "api-z1.qiniu.com",
		"hn": "api-z2.qiniu.com",
		"na": "api-na0.qiniu.com",
		"sa": "api-as0.qiniu.com",
	}
)
