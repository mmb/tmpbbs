package tmpbbs

import (
	"container/list"
	"context"
	"crypto/tls"
	"log/slog"
	"strings"
	"time"

	"github.com/mmb/tmpbbs/internal/tmpbbs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

type pullPeer struct {
	client       proto.PostSyncClient
	logger       *slog.Logger
	postStore    *PostStore
	address      string
	lastIDSynced string
	rootID       string
	interval     time.Duration
	maxResults   int32
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
		interval:   interval,
		client:     proto.NewPostSyncClient(clientConn),
		postStore:  postStore,
		logger:     slog.Default().With("serverAddress", address),
		address:    address,
		maxResults: maxMaxResults,
	}, nil
}

func (pp *pullPeer) run(initialWait time.Duration) {
	time.Sleep(initialWait)

	for {
		if pp.sync() < int(pp.maxResults) {
			time.Sleep(pp.interval)
		}
	}
}

func (pp *pullPeer) sync() int {
	ctx := context.Background()

	response, err := pp.client.Get(ctx, &proto.PostSyncRequest{Id: pp.lastIDSynced, MaxResults: pp.maxResults},
		grpc.UseCompressor(gzip.Name))
	if err != nil {
		pp.logger.ErrorContext(ctx, err.Error())

		return 0
	}

	protoPosts := response.GetPosts()
	pp.logger.InfoContext(ctx, "received response to peer sync request", "lastIDSynced", pp.lastIDSynced, "numPosts",
		len(protoPosts))

	for _, protoPost := range protoPosts {
		// Root post of peer
		if protoPost.GetParentId() == "" {
			pp.rootID = protoPost.GetId()
			pp.lastIDSynced = protoPost.GetId()

			continue
		}
		// We already have this post, do not add
		if pp.postStore.hasPost(protoPost.GetId()) {
			pp.lastIDSynced = protoPost.GetId()

			continue
		}

		post := &post{
			Title:     protoPost.GetTitle(),
			Author:    protoPost.GetAuthor(),
			Tripcode:  protoPost.GetTripcode(),
			Body:      protoPost.GetBody(),
			Replies:   list.New(),
			id:        protoPost.GetId(),
			time:      protoPost.GetTime().AsTime(),
			superuser: protoPost.GetSuperuser(),
		}

		// If the parent is the peer's root, add it to our root
		if protoPost.GetParentId() == pp.rootID {
			pp.postStore.put(post, pp.postStore.rootID)
			pp.lastIDSynced = protoPost.GetId()

			continue
		}

		// If we have the parent, add it to the parent
		if pp.postStore.hasPost(protoPost.GetParentId()) {
			pp.postStore.put(post, protoPost.GetParentId())
			pp.lastIDSynced = protoPost.GetId()

			continue
		}

		// We don't have the parent, start a resync from the peer root.
		pp.lastIDSynced = ""
		pp.logger.WarnContext(ctx, "resync from root", "missingParentID", protoPost.GetParentId())

		return 0
	}

	return len(protoPosts)
}

// RunPullPeers creates a PullPeer for each peer address and starts syncing.
// It calculates a time to wait before starting for each peer so they run
// staggered.
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
