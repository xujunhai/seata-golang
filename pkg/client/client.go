package client

import (
	"google.golang.org/grpc/backoff"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/opentrx/seata-golang/v2/pkg/apis"
	"github.com/opentrx/seata-golang/v2/pkg/client/config"
	"github.com/opentrx/seata-golang/v2/pkg/client/rm"
	"github.com/opentrx/seata-golang/v2/pkg/client/tcc"
	"github.com/opentrx/seata-golang/v2/pkg/client/tm"
)

// Init init resource manager，init transaction manager, expose a port to listen tc
// call back request.
func Init(config *config.Configuration) {
	conn, err := grpc.Dial(config.ServerAddressing,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(config.GetClientParameters()),
		grpc.WithReturnConnectionError(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   120 * time.Second,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
		grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	resourceManagerClient := apis.NewResourceManagerServiceClient(conn)
	transactionManagerClient := apis.NewTransactionManagerServiceClient(conn)

	rm.InitResourceManager(config.Addressing, resourceManagerClient)
	tm.InitTransactionManager(config.Addressing, transactionManagerClient)
	rm.RegisterTransactionServiceServer(tcc.GetTCCResourceManager())
}
