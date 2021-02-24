package jh

// IJh 鉴黄类接口
type IJh interface {
	CheckOne(fs JHFile)
	CheckMore(fs []JHFile)
}

type JhConfig struct {
	JType string
}

type Jh struct {
	conf   JhConfig
	Driver IJh
}

// RiskLabel 命中标签
type RiskLabel struct {
	Score      float64 `json:"score"`
	Label      string  `json:"label"`
	Suggestion string  `json:"suggestion"`
}

// JHFile 鉴黄文件接口
type JHFile interface {
	// GetURL return file url
	GetURL() string
	// GetHash return file hash value
	GetHash() string
	// GetName return file name
	GetName() string
	// GetMime return file mime
	GetMime() string
	//GetSize return file size
	GetSize() int64
	// GetImgPixelSize return imgs piexl size
	GetImgPixelSize() int64
	// SaveJhResult 保存鉴黄
	// param string driverName 鉴黄供应商
	// param bool isPass 是否通过鉴黄
	// param string rspData 鉴黄结果
	// param string errorMsg 异常错误
	SaveJhResult(driverName string, HandelCode int, ReqID string, RiskDescription string, rspData string, rl []RiskLabel, errorMsg string)
}

// 检测处理方式
const (
	//JhIgnore 检测不成功,忽略
	JhIgnore int = 0
	//JhPass 通过审核
	JhPass int = 1
	//JhReview 需要人工审核
	JhReview int = 2
	//JhReject 拒绝
	JhReject int = 3
)

func (j *Jh) Check(fs ...JHFile) {

	j.Driver.CheckMore(fs)
}
