package model

import "time"

type User struct {
	UserName      string    `json:"user_name"`
	Name          string    `json:"name"`
	QuantityKudos int32     `json:"quantity_kudos"`
	CreateAt      time.Time `json:"create_at"`
	UpdateAt      time.Time `json:"update_at"`
	Kudos         []*Kudos  `json:"kudos,omitempty"`
}
