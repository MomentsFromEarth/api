package models

// Tag model
type Tag struct {
	MFEKey     string `json:"mfe_key"`
	TagID      string `json:"tag_id"`
	Name       string `json:"name"`
	Count      int32  `json:"count"`
	Created    int64  `json:"created"`
	Updated    int64  `json:"updated"`
	QueryKey01 string `json:"query_key_01"`
	QueryKey02 string `json:"query_key_02"`
}

// NewTag model
type NewTag struct {
	Name string `json:"name"`
}
