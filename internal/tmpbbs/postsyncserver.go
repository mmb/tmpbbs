package tmpbbs

import (
	"context"
	"log/slog"

	"github.com/mmb/tmpbbs/internal/tmpbbs/proto"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostSyncServer struct {
	proto.UnimplementedPostSyncServer
	postStore *PostStore
}

const maxMaxResults = 500

func NewPostSyncServer(postStore *PostStore) *PostSyncServer {
	return &PostSyncServer{
		postStore: postStore,
	}
}

func (pss *PostSyncServer) Get(context context.Context, request *proto.PostSyncRequest) (*proto.PostSyncResponse,
	error,
) {
	var clientAddress string

	peer, exists := peer.FromContext(context)

	if exists {
		clientAddress = peer.Addr.String()
	}

	logger := slog.Default().With("clientAddress", clientAddress)

	maxResults := min(int(request.GetMaxResults()), maxMaxResults)
	if maxResults == 0 {
		maxResults = maxMaxResults
	}

	posts := pss.postStore.getSince(request.GetUuid(), maxResults)
	protoPosts := make([]*proto.Post, len(posts))

	for index, post := range posts {
		var parentUUID string
		if post.Parent != nil {
			parentUUID = post.Parent.uuid
		}

		protoPosts[index] = &proto.Post{
			Time:       timestamppb.New(post.time),
			Title:      post.Title,
			Author:     post.Author,
			Tripcode:   post.Tripcode,
			Body:       post.Body,
			Uuid:       post.uuid,
			ParentUuid: parentUUID,
		}
	}

	logger.Info("responded to peer sync request", "sinceUUID", request.GetUuid(),
		"maxResults", request.GetMaxResults(), "numResults", len(protoPosts))

	return &proto.PostSyncResponse{
		Posts: protoPosts,
	}, nil
}
