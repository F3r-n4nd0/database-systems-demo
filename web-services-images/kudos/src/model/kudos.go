package model

import "time"

type Kudos struct {
	Id           string    `json:"id"`
	FromUserName string    `json:"from_user_name"`
	ToUserName   string    `json:"to_user_name"`
	Message      string    `json:"message"`
	CreateAt     time.Time `json:"create_at"`
	UpdateAt     time.Time `json:"update_at"`
}
