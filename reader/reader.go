package reader

type Reader interface {
	GetInputChan() chan []byte
}
