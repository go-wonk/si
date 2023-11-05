package sirabbitmq_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/go-wonk/si/v2/sirabbitmq"
)

func Test_Logger(t *testing.T) {
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	var logger sirabbitmq.Logger = sirabbitmq.NewDefaultLogger(l)

	ctx := context.Background()
	logger.Debug(ctx, "asdf")

	sirabbitmq.SetLogger(logger)
	sirabbitmq.Debug("msg")
}
