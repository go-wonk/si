package sigrpc_test

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/go-wonk/si/v2/sigrpc/tests/protos"
	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	if !onlinetest {
		t.Skip("skipping online tests")
	}

	c := pb.NewStudentClient(client)
	rep, err := c.Read(context.Background(), &pb.StudentRequest{
		Name: "wonka",
	})
	assert.Nil(t, err)

	fmt.Println(rep.String())
}
