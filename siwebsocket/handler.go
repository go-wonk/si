package siwebsocket

import (
	"io"

	"github.com/go-wonk/si/sicore"
)

type MessageHandler interface {
	Handle(r io.Reader, opts ...sicore.ReaderOption)
}

type NopMessageHandler struct{}

func (o *NopMessageHandler) Handle(r io.Reader, opts ...sicore.ReaderOption) {
	// do nothing
	io.Copy(io.Discard, r)
}

type DefaultMessageHandler struct{}

func (o *DefaultMessageHandler) Handle(r io.Reader, opts ...sicore.ReaderOption) {
	// log.Println(string(b))
	sr := sicore.GetReader(r, opts...)
	defer sicore.PutReader(sr)

	_, err := sr.ReadAll()
	if err != nil {
		// log.Println(err)
		return
	}

	// log.Println(string(b))
}
