package webservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"stats/src/kudos"
	"stats/src/model"
)

type usersWebService struct {
	host string
	port string
}

func NewUsersWebService(host string, port string) kudos.WebServiceUsers {
	return &usersWebService{host, port}
}

func (uws *usersWebService) UpdateQuantityKudos(userName string, quantity int) error {

	message := map[string]interface{}{
		"user_name": userName,
		"quantity":  quantity,
	}
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s:%s/v1/user/kudos", uws.host, uws.port)
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return model.WebServiceError
		}
		return errors.New(string(body))
	}
	return nil

}
