package siwebsocket

import (
	"io"
	"log"

	"github.com/go-wonk/si/sicore"
)

type MessageHandler interface {
	Handle(r io.Reader, opts ...sicore.ReaderOption) error
}

type NopMessageHandler struct{}

func (o *NopMessageHandler) Handle(r io.Reader, opts ...sicore.ReaderOption) error {
	// discard reader data
	_, err := io.Copy(io.Discard, r)
	return err
}

type DefaultMessageHandler struct{}

func (o *DefaultMessageHandler) Handle(r io.Reader, opts ...sicore.ReaderOption) error {
	sr := sicore.GetReader(r, opts...)
	defer sicore.PutReader(sr)

	_, err := sr.ReadAll()
	if err != nil {
		return err
	}

	return nil
}

type DefaultMessageLogHandler struct{}

func (o *DefaultMessageLogHandler) Handle(r io.Reader, opts ...sicore.ReaderOption) error {
	sr := sicore.GetReader(r, opts...)
	defer sicore.PutReader(sr)

	b, err := sr.ReadAll()
	if err != nil {
		return err
	}
	log.Println(string(b))
	return nil
}
