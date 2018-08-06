package kudos

type WebServiceUsers interface {
	UpdateQuantityKudos(userName string, quantity int) error
}
