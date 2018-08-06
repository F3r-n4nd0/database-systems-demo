package kudos

import "stats/src/model"

type UseCase interface {
	QueueKudos(kudosQ *model.KudosQueue)
}
