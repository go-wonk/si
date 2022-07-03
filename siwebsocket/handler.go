package siwebsocket

import "log"

type MessageHandler interface {
	Handle(b []byte)
}

type NopMessageHandler struct{}

func (o *NopMessageHandler) Handle(b []byte) {
	// do nothing
}

type DefaultMessageHandler struct{}

func (o *DefaultMessageHandler) Handle(b []byte) {
	log.Println(string(b))
}
