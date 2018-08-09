package webservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"web-service-users/src/model"
	"web-service-users/src/user"
)

type kudosWebService struct {
	host string
	port string
}

func NewKudosWebService(host string, port string) user.WebServiceKudos {
	return &kudosWebService{host, port}
}

func (kws *kudosWebService) GetKudosByUserName(userName string) ([]*model.Kudos, error) {

	url := fmt.Sprintf("http://%s:%s/v1/user/kudos/%s", kws.host, kws.port, userName)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var kudos []*model.Kudos
	err = json.Unmarshal(body, &kudos)
	if err != nil {
		return nil, err
	}
	return kudos, nil

}
