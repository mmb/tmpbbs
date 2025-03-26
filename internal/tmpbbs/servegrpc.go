package tmpbbs

import (
	"crypto/tls"
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
	} else {
		grpcServer = grpc.NewServer()
	}

	proto.RegisterPostSyncServer(grpcServer, postSyncServer)

	return grpcServer.Serve(listener)
}
