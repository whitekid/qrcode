package grpcserver

import (
	"context"
	"net"
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/pkg/errors"
	"github.com/whitekid/goxp/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"qrcodeapi/proto"
)

func Run(ctx context.Context, bindAddr string) error {
	logger := log.Zap(log.New(zap.AddCallerSkip(2)))

	var kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	var kasp = keepalive.ServerParameters{
		Time:    5 * time.Second, // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout: time.Second,     // Wait 1 second for the ping ack before assuming the connection is dead
	}

	opts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
	}

	// NOTE interceptor는 한번만 설정 가능
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_zap.UnaryServerInterceptor(logger),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_zap.StreamServerInterceptor(logger),
	}
	unaryInterceptors = append(unaryInterceptors)
	streamInterceptors = append(streamInterceptors)

	g := grpc.NewServer(opts...)

	service := New()
	proto.RegisterQRCodeServer(g, service)

	ln, err := net.Listen("tcp", bindAddr)
	if err != nil {
		return errors.Wrapf(err, "fail to listen")
	}

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	log.Infof("start grpc server at %s", ln.Addr().String())

	return g.Serve(ln)
}
