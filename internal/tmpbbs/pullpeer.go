package tmpbbs

import (
	"context"
	"crypto/tls"
	"log/slog"
	"strings"
	"time"

	"github.com/mmb/tmpbbs/internal/tmpbbs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type pullPeer struct {
	client         proto.PostSyncClient
	postStore      *PostStore
	logger         *slog.Logger
	address        string
	rootUUID       string
	lastUUIDSynced string
	interval       time.Duration
}

func newPullPeer(address string, interval time.Duration, postStore *PostStore) (*pullPeer, error) {
	var creds credentials.TransportCredentials

	if strings.HasPrefix(address, "tls://") {
		address = address[6:]
		creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}) // #nosec G402 -- should work with self-signed certs
	} else {
		creds = insecure.NewCredentials()
	}

	clientConn, err := grpc.NewClient(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	return &pullPeer{
		interval:  interval,
		client:    proto.NewPostSyncClient(clientConn),
		postStore: postStore,
		logger:    slog.Default().With("serverAddress", address),
		address:   address,
	}, nil
}

func (pp *pullPeer) run(wait time.Duration) {
	time.Sleep(wait)
	pp.sync()
	ticker := time.NewTicker(pp.interval)

	for {
		<-ticker.C
		pp.sync()
	}
}

func (pp *pullPeer) sync() {
	response, err := pp.client.Get(context.Background(), &proto.PostSyncRequest{Uuid: pp.lastUUIDSynced})
	if err != nil {
		pp.logger.Error(err.Error())
	}

	protoPosts := response.GetPosts()
	pp.logger.Info("received response to peer sync request", "lastUUIDSynced", pp.lastUUIDSynced,
		"numPosts", len(protoPosts))

	for _, protoPost := range protoPosts {
		// Root post of peer
		if protoPost.GetParentUuid() == "" {
			pp.rootUUID = protoPost.GetUuid()
			pp.lastUUIDSynced = protoPost.GetUuid()

			continue
		}
		// We already have this post, do not add
		if pp.postStore.getPostByUUID(protoPost.GetUuid()) != nil {
			pp.lastUUIDSynced = protoPost.GetUuid()

			continue
		}

		post := &post{
			Title:    protoPost.GetTitle(),
			Author:   protoPost.GetAuthor(),
			Tripcode: protoPost.GetTripcode(),
			Body:     protoPost.GetBody(),
			uuid:     protoPost.GetUuid(),
			time:     protoPost.GetTime().AsTime(),
		}

		// If the parent is the peer's root, add it to our root
		if protoPost.GetParentUuid() == pp.rootUUID {
			pp.postStore.put(post, pp.postStore.posts[0].uuid)
			pp.lastUUIDSynced = protoPost.GetUuid()

			continue
		}

		// If we have the parent, add it to the parent
		if pp.postStore.getPostByUUID(protoPost.GetParentUuid()) != nil {
			pp.postStore.put(post, protoPost.GetParentUuid())
			pp.lastUUIDSynced = protoPost.GetUuid()

			continue
		}

		// We don't have the parent, start a resync from the peer root.
		pp.lastUUIDSynced = ""
		pp.logger.Warn("resync from root", "missingParentUUID", protoPost.GetParentUuid())

		return
	}
}

func RunPullPeers(addresses []string, interval time.Duration, postStore *PostStore) error {
	var waitBetween time.Duration
	if len(addresses) > 0 {
		waitBetween = interval / time.Duration(len(addresses))
	}

	for index, address := range addresses {
		pullPeer, err := newPullPeer(address, interval, postStore)
		if err != nil {
			return err
		}

		go pullPeer.run(time.Duration(index) * waitBetween)
	}

	return nil
}
