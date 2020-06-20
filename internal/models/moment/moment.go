package models

// Moment model
type Moment struct {
	MFEKey      string `json:"mfe_key"`
	MomentID    string `json:"moment_id"`
	Creator     string `json:"creator"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Size        int64  `json:"size"`
	Filename    string `json:"filename"`
	QueueID     string `json:"queue_id"`
	HostID      string `json:"host_id"`
	Status      string `json:"status"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
	QueryKey01  string `json:"query_key_01"`
	QueryKey02  string `json:"query_key_02"`
}

// NewMoment model
type NewMoment struct {
	Creator     string `json:"creator"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Size        int64  `json:"size"`
	Filename    string `json:"filename"`
	QueueID     string `json:"queue_id"`
}
