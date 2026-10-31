package tmpbbs

import (
	"container/list"
	"context"
	"errors"
	"iter"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// A PostStore stores Posts in memory and provides safety for concurrent
// access.
type PostStore struct {
	idMap         map[string]*post
	posts         *list.List
	rootPost      *post
	pruneInterval time.Duration
	pruneMaxAge   time.Duration
	mutex         sync.RWMutex
}

var errParentNotFound = errors.New("parent post not found")

// NewPostStore returns a new PostStore. It also creates the root Post.
func NewPostStore(title string, pruneInterval time.Duration, pruneMaxAge time.Duration) *PostStore {
	rootPost := newPost(title, "", "", nil)
	postStore := &PostStore{
		idMap:         make(map[string]*post),
		posts:         list.New(),
		pruneInterval: pruneInterval,
		pruneMaxAge:   pruneMaxAge,
		rootPost:      rootPost,
	}
	rootPost.postsElement = postStore.posts.PushBack(rootPost)
	postStore.idMap[rootPost.id] = rootPost

	if pruneInterval > 0 {
		go postStore.startPruner()
	}

	return postStore
}

// LoadYAML loads Posts from a YAML file on the filesystem into the PostStore.
func (ps *PostStore) LoadYAML(path string, tripcoder *Tripcoder) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var posts []post

	yamlErr := yaml.Unmarshal(data, &posts)
	if yamlErr != nil {
		return yamlErr
	}

	var putErr error
	for i := range posts {
		putErr = ps.put(newPost(posts[i].Title, posts[i].Author, posts[i].Body, tripcoder), ps.rootPost.id)
		if putErr != nil {
			return putErr
		}
	}

	return nil
}

// delete removes all references to a post from the PostStore.
// Write lock must be obtained before calling.
func (ps *PostStore) delete(postToDelete *post) {
	if ps.idMap[postToDelete.id] == nil {
		return
	}

	for reply := postToDelete.Replies.Front(); reply != nil; reply = reply.Next() {
		ps.delete(reply.Value.(*post)) //nolint:errcheck,forcetypeassert // only one type
	}

	postToDelete.Parent.Replies.Remove(postToDelete.parentRepliesElement)

	delete(ps.idMap, postToDelete.id)
	ps.posts.Remove(postToDelete.postsElement)
}

func (ps *PostStore) get(postID string, callback func(*post)) bool {
	var callbackPost *post

	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	if postID == "" {
		callbackPost = ps.rootPost
	} else if callbackPost = ps.getPostByID(postID); callbackPost == nil {
		return false
	}

	callback(callbackPost)

	return true
}

// getPostByID returns the post with ID postID or nil if not found.
// The caller must handle locking.
func (ps *PostStore) getPostByID(postID string) *post {
	return ps.idMap[postID]
}

func (ps *PostStore) getSince(postID string, maxPosts int) iter.Seq[*post] {
	return func(yield func(*post) bool) {
		ps.mutex.RLock()
		defer ps.mutex.RUnlock()

		start := ps.posts.Front()
		if p := ps.idMap[postID]; p != nil {
			start = p.postsElement.Next()
		}

		for count, element := 0, start; count < maxPosts && element != nil; count, element = count+1, element.Next() {
			if !yield(element.Value.(*post)) { //nolint:errcheck,forcetypeassert // only one type
				return
			}
		}
	}
}

func (ps *PostStore) prune(ctx context.Context) {
	slog.InfoContext(ctx, "prune start")

	var deletePosts []*post

	ps.mutex.Lock()

	beforeCount := ps.posts.Len()
	for element := ps.posts.Front().Next(); element != nil; element = element.Next() {
		checkPost := element.Value.(*post) //nolint:errcheck,forcetypeassert // only one type
		if time.Since(checkPost.lastActivity) > ps.pruneMaxAge {
			deletePosts = append(deletePosts, checkPost)
		}
	}

	for _, p := range deletePosts {
		ps.delete(p)
	}

	afterCount := ps.posts.Len()
	ps.mutex.Unlock()

	slog.InfoContext(ctx, "prune end", "prunedCount", beforeCount-afterCount)
}

func (ps *PostStore) put(post *post, parentID string) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if ps.idMap[post.id] != nil {
		return nil
	}

	if post.Parent = ps.getPostByID(parentID); post.Parent == nil {
		return errParentNotFound
	}

	post.parentRepliesElement = post.Parent.Replies.PushFront(post)
	post.postsElement = ps.posts.PushBack(post)
	ps.idMap[post.id] = post

	if (post.IsOriginalPoster() || post.IsSuperuser()) && strings.HasPrefix(post.Body, "!delete") &&
		post.Parent != ps.rootPost {
		post.Parent.delete()
	}

	post.bump()

	return nil
}

func (ps *PostStore) startPruner() {
	ticker := time.NewTicker(ps.pruneInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		ps.prune(context.Background())
	}
}
