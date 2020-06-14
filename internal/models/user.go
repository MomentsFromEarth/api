package models

// User model
type User struct {
	Email           string `json:"email"`
	UserID          string `json:"user_id"`
	UserName        string `json:"username"`
	Created         int64  `json:"created"`
	Updated         int64  `json:"updated"`
	JoinMailingList bool   `json:"join_mailing_list"`
	NewUser         bool   `json:"new_user"`
}
