package webservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"stats/src/kudos"
)

type kudosWebService struct {
	host string
	port string
}

func NewKudosWebService(host string, port string) kudos.WebServiceKudos {
	return &kudosWebService{host, port}
}

type responseQuantity struct {
	Quantity int `json:"quantity"`
}

func (kws *kudosWebService) GetQuantityKudos(userName string) (int, error) {

	url := fmt.Sprintf("http://%s:%s/v1/quantity/kudos/%s", kws.host, kws.port, userName)
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	var data responseQuantity
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}
	return data.Quantity, nil

}
