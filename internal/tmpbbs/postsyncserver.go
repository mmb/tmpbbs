package tmpbbs

import (
	"context"

	"github.com/mmb/tmpbbs/internal/tmpbbs/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostSyncServer struct {
	proto.UnimplementedPostSyncServer
	postStore *PostStore
}

const maxMaxResults = 1024

func NewPostSyncServer(postStore *PostStore) *PostSyncServer {
	return &PostSyncServer{
		postStore: postStore,
	}
}

func (pss *PostSyncServer) Get(_ context.Context, request *proto.PostSyncRequest) (*proto.PostSyncResponse, error) {
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

	return &proto.PostSyncResponse{
		Posts: protoPosts,
	}, nil
}
