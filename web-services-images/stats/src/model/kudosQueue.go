package model

import "encoding/json"

type KudosQueue struct {
	UserName string `json:"user_name"`
	Type     string `json:"type"`
}

func UnMarshal(bytes []byte) (*KudosQueue, error) {

	var kudosQ KudosQueue
	err := json.Unmarshal(bytes, &kudosQ)
	if err != nil {
		return nil, err
	}
	return &kudosQ, nil

}
