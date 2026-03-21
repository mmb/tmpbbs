package tmpbbs

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"time"

	"github.com/mmb/tmpbbs/internal/tmpbbs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	grpcKeepAliveTime          = 25 * time.Second
	grpcKeepAliveTimeout       = 10 * time.Second
	grpcServerKeepAliveMinTime = 20 * time.Second
)

// ServeGRPC creates and configures a grpc.Server then starts listening.
func ServeGRPC(ctx context.Context, listenAddress string, tlsCertFile string, tlsKeyFile string,
	postSyncServer *PostSyncServer,
) error {
	var listenConfig net.ListenConfig

	listener, err := listenConfig.Listen(ctx, "tcp", listenAddress)
	if err != nil {
		return err
	}
	defer listener.Close()

	serverOptions := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcServerKeepAliveMinTime,
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepAliveTime,
			Timeout: grpcKeepAliveTimeout,
		}),
	}

	if tlsCertFile != "" && tlsKeyFile != "" {
		var certificate tls.Certificate

		certificate, err = tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
		if err != nil {
			return err
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{certificate},
			MinVersion:   tls.VersionTLS13,
		}
		serverOptions = append(serverOptions, grpc.Creds(credentials.NewTLS(config)))
	}

	grpcServer := grpc.NewServer(serverOptions...)
	proto.RegisterPostSyncServer(grpcServer, postSyncServer)
	slog.InfoContext(ctx, "listening for gRPC", "address", listenAddress, "tlsCertFile", tlsCertFile, "tlsKeyFile",
		tlsKeyFile)

	return grpcServer.Serve(listener)
}
