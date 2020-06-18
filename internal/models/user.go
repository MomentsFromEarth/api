package models

// ProfileUser model
type ProfileUser struct {
	MFEKey          string `json:"-"`
	Email           string `json:"email"`
	UserID          string `json:"user_id"`
	UserName        string `json:"username"`
	Avatar          string `json:"avatar"`
	CognitoSub      string `json:"-"`
	Created         int64  `json:"created"`
	Updated         int64  `json:"updated"`
	JoinMailingList bool   `json:"join_mailing_list"`
	NewUser         bool   `json:"new_user"`
	QueryKey01      string `json:"-"`
	QueryKey02      string `json:"-"`
}

// User model
type User struct {
	MFEKey          string `json:"mfe_key"`
	Email           string `json:"email"`
	UserID          string `json:"user_id"`
	UserName        string `json:"username"`
	Avatar          string `json:"avatar"`
	CognitoSub      string `json:"cognito_sub"`
	Created         int64  `json:"created"`
	Updated         int64  `json:"updated"`
	JoinMailingList bool   `json:"join_mailing_list"`
	NewUser         bool   `json:"new_user"`
	QueryKey01      string `json:"query_key_01"`
	QueryKey02      string `json:"query_key_02"`
}

// NewUser model
type NewUser struct {
	Email           string `json:"email"`
	CognitoSub      string `json:"cognito_sub"`
	JoinMailingList bool   `json:"join_mailing_list"`
	Avatar          string `json:"avatar"`
}
