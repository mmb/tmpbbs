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
	listenConfig := net.ListenConfig{}

	listener, err := listenConfig.Listen(ctx, "tcp", listenAddress)
	if err != nil {
		return err
	}
	defer listener.Close()

	var grpcServer *grpc.Server

	keepAliveEnforcementPolicy := grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             grpcServerKeepAliveMinTime,
		PermitWithoutStream: true,
	})
	keepAliveParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    grpcKeepAliveTime,
		Timeout: grpcKeepAliveTimeout,
	})

	if tlsCertFile != "" && tlsKeyFile != "" {
		var certificate tls.Certificate

		certificate, err = tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
		if err != nil {
			return err
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{certificate},
			ClientAuth:   tls.NoClientCert,
			MinVersion:   tls.VersionTLS13,
		}
		grpcServer = grpc.NewServer(keepAliveEnforcementPolicy, keepAliveParams, grpc.Creds(credentials.NewTLS(config)))
	} else {
		grpcServer = grpc.NewServer(keepAliveEnforcementPolicy, keepAliveParams)
	}

	proto.RegisterPostSyncServer(grpcServer, postSyncServer)

	slog.InfoContext(ctx, "listening for gRPC", "address", listenAddress, "tlsCertFile", tlsCertFile, "tlsKeyFile",
		tlsKeyFile)

	return grpcServer.Serve(listener)
}
