package repository

import (
	"web-service-kudos/src/kudos"
	"web-service-kudos/src/model"

	"encoding/json"

	"fmt"

	"github.com/go-redis/redis"
)

type redisKudosRepository struct {
	client *redis.Client
}

func NewRedisKudosRepository(client *redis.Client) kudos.Repository {
	return &redisKudosRepository{client}
}

func (r *redisKudosRepository) kudosToString(kudos *model.Kudos) (string, error) {

	bytes, err := json.Marshal(kudos)
	if err != nil {
		return "", err
	}
	return string(bytes), nil

}

func (r *redisKudosRepository) stringToKudos(kudosString string) (*model.Kudos, error) {

	kudosModel := model.Kudos{}
	err := json.Unmarshal([]byte(kudosString), &kudosModel)
	if err != nil {
		return nil, err
	}
	return &kudosModel, nil

}

func (r *redisKudosRepository) Store(kudos *model.Kudos) error {

	compId := fmt.Sprintf("%s:%s", kudos.ToUserName, kudos.Id)
	kudosString, err := r.kudosToString(kudos)
	if err != nil {
		return err
	}
	err = r.client.Set(compId, kudosString, 0).Err()
	if err != nil {
		return err
	}
	return nil

}

func (r *redisKudosRepository) GetByID(id string) (*model.Kudos, error) {

	filter := fmt.Sprintf("*:%s", id)
	keys, err := r.client.Keys(filter).Result()
	if err != nil {
		return nil, err
	}
	if len(keys) != 1 {
		return nil, model.NotFoundError
	}
	val, err := r.client.Get(keys[0]).Result()
	if err != nil {
		return nil, err
	}
	kudosFound, err := r.stringToKudos(val)
	if err != nil {
		return nil, err
	}
	return kudosFound, nil

}

func (r *redisKudosRepository) DeleteByID(id string) error {

	filter := fmt.Sprintf("*:%s", id)
	keys, err := r.client.Keys(filter).Result()
	if err != nil {
		return err
	}
	if len(keys) != 1 {
		return model.NotFoundError
	}
	val, err := r.client.Del(keys[0]).Result()
	if err != nil {
		return err
	}
	if val != 1 {
		return model.NotFoundError
	}
	return nil

}

func (r *redisKudosRepository) FetchAllKudos(pageSize int64, numberPage int64) ([]*model.Kudos, error) {

	keys, err := r.client.Keys("*").Result()
	if err != nil {
		return nil, err
	}
	skip := pageSize * (numberPage - 1)
	filterKeys := paginateKeys(keys, int(skip), int(pageSize))
	var allKudos []*model.Kudos
	for _, key := range filterKeys {
		kudo, err := r.getByAbsoluteID(key)
		if err != nil {
			return nil, err
		}
		allKudos = append(allKudos, kudo)
	}
	return allKudos, nil

}

func paginateKeys(x []string, skip int, size int) []string {

	if skip > len(x) {
		skip = len(x)
	}
	end := skip + size
	if end > len(x) {
		end = len(x)
	}
	return x[skip:end]

}

func (r *redisKudosRepository) GetQuantityByUserName(userName string) (int, error) {

	filter := fmt.Sprintf("%s:*", userName)
	keys, err := r.client.Keys(filter).Result()
	if err != nil {
		return 0, err
	}
	return len(keys), nil

}

func (r *redisKudosRepository) GetByUserName(userName string) ([]*model.Kudos, error) {

	filter := fmt.Sprintf("%s:*", userName)
	keys, err := r.client.Keys(filter).Result()
	if err != nil {
		return nil, err
	}
	var allKudos []*model.Kudos
	for _, key := range keys {
		kudo, err := r.getByAbsoluteID(key)
		if err != nil {
			return nil, err
		}
		allKudos = append(allKudos, kudo)
	}
	return allKudos, nil

}

func (r *redisKudosRepository) getByAbsoluteID(id string) (*model.Kudos, error) {

	val, err := r.client.Get(id).Result()
	if err != nil {
		return nil, err
	}
	kudosFound, err := r.stringToKudos(val)
	if err != nil {
		return nil, err
	}
	return kudosFound, nil

}
