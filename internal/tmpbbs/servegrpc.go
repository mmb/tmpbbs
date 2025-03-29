package tmpbbs

import (
	"crypto/tls"
	"log/slog"
	"net"

	"github.com/mmb/tmpbbs/internal/tmpbbs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func ServeGRPC(listenAddress string, tlsCertFile string, tlsKeyFile string, postSyncServer *PostSyncServer) error {
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return err
	}

	var grpcServer *grpc.Server

	tlsEnabled := false

	if tlsCertFile != "" && tlsKeyFile != "" {
		var certificate tls.Certificate

		certificate, err = tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
		if err != nil {
			return err
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{certificate},
			ClientAuth:   tls.NoClientCert,
			MinVersion:   tls.VersionTLS12,
		}
		grpcServer = grpc.NewServer(grpc.Creds(credentials.NewTLS(config)))
		tlsEnabled = true
	} else {
		grpcServer = grpc.NewServer()
	}

	proto.RegisterPostSyncServer(grpcServer, postSyncServer)

	slog.Info("listening for gRPC", "address", listenAddress, "tlsEnabled", tlsEnabled)

	return grpcServer.Serve(listener)
}
