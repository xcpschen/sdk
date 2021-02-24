package ctrlvideo

// CallbackBody ctrlVideo callback request body
type CallbackBody struct {
	CommRequest
	Action    string `json:"action"`
	ResultStr string `json:"resultstr"`
}

// ReleaseResult CallbackBody action is release
type ReleaseResult struct {
	OpenID      string `json:"openid"`
	ProjectID   string `json:"project_id"`
	Status      string `json:"status"`
	ReleseTime  string `json:"release_time"`
	PlayerTitle string `json:"player_title"`
	PlayerCover string `json:"player_cover"`
	ShareTitle  string `json:"share_title"`
	ShareDesc   string `json:"share_desc"`
	ShareThumb  string `json:"share_thumb"`
}

// AuditResult CallbackBody action is audit
type AuditResult struct {
	OpenID    string `json:"openid"`
	ProjectID string `json:"project_id"`
	Status    string `json:"status"`
}

// AuditResult Status infometion
const (
	AduitingStatus    string = "auditing"
	AduitPassStatus   string = "pass"
	AduitRejectStatus string = "reject"
)

// DownlineResult CallbackBody action is downline
type DownlineResult struct {
	Downline  int    `json:"downline"`
	OpenID    string `json:"openid"`
	ProjectID string `json:"project_id"`
}

// Downline infometion
const (
	UpLineCode int = iota
	DownLineCode
)
