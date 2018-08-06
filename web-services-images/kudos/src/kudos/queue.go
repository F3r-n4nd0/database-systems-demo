package kudos

import "web-service-kudos/src/model"

type Queue interface {
	IncreasesKudos(kudos *model.Kudos) error
	DecreasesKudos(kudos *model.Kudos) error
}
