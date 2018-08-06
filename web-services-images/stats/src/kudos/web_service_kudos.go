package kudos

type WebServiceKudos interface {
	GetQuantityKudos(userName string) (int, error)
}
