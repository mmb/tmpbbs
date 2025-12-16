package tmpbbs

import (
	"context"
	"log/slog"

	"github.com/mmb/tmpbbs/internal/tmpbbs/proto"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// A PostSyncServer is a gRPC server that returns Posts.
type PostSyncServer struct {
	proto.UnimplementedPostSyncServer

	postStore *PostStore
}

const maxMaxResults = 500

// NewPostSyncServer returns a new PostSyncServer.
func NewPostSyncServer(postStore *PostStore) *PostSyncServer {
	return &PostSyncServer{
		postStore: postStore,
	}
}

// Get returns Posts in order starting from an ID in the request.
func (pss *PostSyncServer) Get(ctx context.Context, request *proto.PostSyncRequest) (*proto.PostSyncResponse,
	error,
) {
	var clientAddress string

	p, exists := peer.FromContext(ctx)

	if exists {
		clientAddress = p.Addr.String()
	}

	logger := slog.Default().With("clientAddress", clientAddress)

	maxResults := min(int(request.GetMaxResults()), maxMaxResults)
	if maxResults == 0 {
		maxResults = maxMaxResults
	}

	posts := pss.postStore.getSince(request.GetId(), maxResults)
	protoPosts := make([]*proto.Post, len(posts))

	for index, post := range posts {
		var parentID string
		if post.Parent != nil {
			parentID = post.Parent.id
		}

		protoPosts[index] = &proto.Post{
			Time:      timestamppb.New(post.time),
			Title:     post.Title,
			Author:    post.Author,
			Tripcode:  post.Tripcode,
			Body:      post.Body,
			Id:        post.id,
			ParentId:  parentID,
			Superuser: post.IsSuperuser(),
		}
	}

	logger.InfoContext(ctx, "responded to peer sync request", "sinceID", request.GetId(),
		"maxResults", maxResults, "numResults", len(protoPosts))

	return &proto.PostSyncResponse{
		Posts: protoPosts,
	}, nil
}
