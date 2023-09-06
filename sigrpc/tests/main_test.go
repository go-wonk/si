package sigrpc_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-wonk/si/v2/sigrpc"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/go-wonk/si/v2/sigrpc/tests/protos"
)

var (
	onlinetest, _ = strconv.ParseBool(os.Getenv("ONLINE_TEST"))

	server     *sigrpc.Server
	serverAddr = ":60000"

	client *grpc.ClientConn
)

// func openClient() *http.Client {
// 	tlsConfig := &tls.Config{
// 		InsecureSkipVerify: true,
// 	}

// 	dialer := &net.Dialer{Timeout: 5 * time.Second}

// 	tr := &http.Transport{
// 		MaxIdleConns:       300,
// 		IdleConnTimeout:    time.Duration(15) * time.Second,
// 		DisableCompression: false,
// 		TLSClientConfig:    tlsConfig,
// 		DisableKeepAlives:  false,
// 		Dial:               dialer.Dial,
// 	}

// 	return sihttp.NewStandardClient(time.Duration(30), tr)
// }

func setup() error {
	var err error
	if onlinetest {
		// build server
		enforcementPolicyUse := true
		enforcementPolicyMinTime := 15
		enforcementPolicyPermitWithoutStream := true
		certPem := "./certs/server.crt"
		certKey := "./certs/server.key"
		keepAliveMaxConnIdle := 300
		keepAliveMaxConnAge := 300
		keepAliveMaxConnAgeGrace := 6
		keepAliveTime := 60
		keepAliveTimeout := 1
		healthCheckUse := true

		server, err = sigrpc.NewServer(serverAddr,
			enforcementPolicyUse, enforcementPolicyMinTime, enforcementPolicyPermitWithoutStream,
			certPem, certKey,
			keepAliveMaxConnIdle, keepAliveMaxConnAge, keepAliveMaxConnAgeGrace, keepAliveTime, keepAliveTimeout,
			healthCheckUse)
		if err != nil {
			return err
		}
		pb.RegisterStudentServer(server.Svr, &studentGrpcServer{})

		go func() {
			server.Start()
		}()

		// build client
		resolveScheme := "student"
		resolveServiceName := "student-svc"
		keepAlivePermitWithoutStream := true
		certServername := "localhost"
		defaultServiceConfig := `{
			"loadBalancingConfig": [{"round_robin":{}}],
			"methodConfig": [{
				"name": [{}],
				"waitForReady": true,
				"retryPolicy": {
					"MaxAttempts": 4,
					"InitialBackoff": ".01s",
					"MaxBackoff": ".01s",
					"BackoffMultiplier": 1.0,
					"RetryableStatusCodes": [ "UNAVAILABLE" ]
				}
			}]
		}`
		dialBlock := false
		dialTimeoutSecond := 6
		client, err = sigrpc.NewClient(serverAddr, resolveScheme, resolveServiceName, keepAliveTime, keepAliveTimeout, keepAlivePermitWithoutStream,
			certPem, certServername, defaultServiceConfig, dialBlock, dialTimeoutSecond)
		if err != nil {
			return err
		}
	}

	return nil
}

func shutdown() {
	if server != nil {
		server.Close()
	}
	if client != nil {
		client.Close()
	}
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Println(err)
		shutdown()
		os.Exit(1)
	}

	exitCode := m.Run()

	shutdown()
	os.Exit(exitCode)
}

type studentGrpcServer struct {
	pb.StudentServer
}

func (d *studentGrpcServer) Read(ctx context.Context, in *pb.StudentRequest) (*pb.StudentReply, error) {
	docs := make([]*pb.StudentEntity, 0)
	docs = append(docs, &pb.StudentEntity{
		Name:        "wonk",
		Age:         10,
		DateTime:    timestamppb.New(time.Now()),
		DoubleValue: 10.1,
	})

	var count int64 = 1
	rep := pb.StudentReply{
		Status:    200,
		Documents: docs,
		Count:     &count,
	}

	return &rep, nil
}
