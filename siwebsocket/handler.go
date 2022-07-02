package siwebsocket

import "log"

type MessageHandler interface {
	Handle(b []byte)
}

type DefaultMessageHandler struct{}

func (o *DefaultMessageHandler) Handle(b []byte) {
	log.Println(string(b))
}
