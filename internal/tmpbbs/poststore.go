package tmpbbs

import (
	"container/list"
	"context"
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
	idMap         map[string]*list.Element
	posts         *list.List
	rootPost      *post
	rootID        string
	pruneInterval time.Duration
	pruneMaxAge   time.Duration
	mutex         sync.RWMutex
}

// NewPostStore returns a new PostStore. It also creates the root Post.
func NewPostStore(title string, pruneInterval time.Duration, pruneMaxAge time.Duration) *PostStore {
	rootPost := newPost(title, "", "", nil)
	postStore := &PostStore{
		idMap:         make(map[string]*list.Element),
		posts:         list.New(),
		rootPost:      rootPost,
		rootID:        rootPost.id,
		pruneInterval: pruneInterval,
		pruneMaxAge:   pruneMaxAge,
	}
	postStore.put(rootPost, "")

	if pruneInterval > 0 {
		go postStore.startPruner()
	}

	return postStore
}

func (ps *PostStore) put(post *post, parentID string) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if post.Parent = ps.getPostByID(parentID); post.Parent != nil {
		post.parentRepliesElement = post.Parent.Replies.PushFront(post)
	}

	post.postsElement = ps.posts.PushBack(post)
	ps.idMap[post.id] = post.postsElement

	if (post.IsOriginalPoster() || post.IsSuperuser()) && strings.HasPrefix(post.Body, "!delete") && post.Parent != nil &&
		post.Parent != ps.rootPost {
		post.Parent.delete()
	}

	post.bump()
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

// delete removes all references to a post from the PostStore.
// Write lock must be obtained before calling.
func (ps *PostStore) delete(postToDelete *post) {
	ctx := context.Background()
	slog.InfoContext(ctx, "delete post", "id", postToDelete.id)

	for reply := postToDelete.Replies.Front(); reply != nil; reply = reply.Next() {
		ps.delete(reply.Value.(*post)) //nolint:errcheck,forcetypeassert // only one type
	}

	postToDelete.Parent.Replies.Remove(postToDelete.parentRepliesElement)

	delete(ps.idMap, postToDelete.id)
	ps.posts.Remove(postToDelete.postsElement)
}

func (ps *PostStore) getSince(postID string, maxPosts int) []*post {
	var start *list.Element

	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	if sinceElement, found := ps.idMap[postID]; found {
		start = sinceElement.Next()
	} else {
		start = ps.posts.Front()
	}

	result := make([]*post, 0, maxPosts)
	for count, element := 0, start; count < maxPosts && element != nil; count, element = count+1, element.Next() {
		result = append(result, element.Value.(*post)) //nolint:errcheck,forcetypeassert // only one type
	}

	return result
}

func (ps *PostStore) hasPost(postID string) bool {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	_, found := ps.idMap[postID]

	return found
}

// LoadYAML loads Posts from a YAML file on the filesystem into the PostStore.
func (ps *PostStore) LoadYAML(path string, tripcoder *Tripcoder) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var posts []post

	if yamlErr := yaml.Unmarshal(data, &posts); yamlErr != nil {
		return yamlErr
	}

	for i := range posts {
		ps.put(newPost(posts[i].Title, posts[i].Author, posts[i].Body, tripcoder), ps.rootID)
	}

	return nil
}

// getPostByID returns the post with ID postID or nil if not found.
// The caller must handle locking.
func (ps *PostStore) getPostByID(postID string) *post {
	element, found := ps.idMap[postID]
	if !found {
		return nil
	}

	return element.Value.(*post) //nolint:errcheck,forcetypeassert // only one type
}

func (ps *PostStore) startPruner() {
	ticker := time.NewTicker(ps.pruneInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		ps.prune()
	}
}

func (ps *PostStore) prune() {
	ctx := context.Background()
	slog.InfoContext(ctx, "prune start")

	var deletePosts []*post

	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	beforeCount := ps.posts.Len()

	for element := ps.posts.Front().Next(); element != nil; element = element.Next() {
		checkPost := element.Value.(*post) //nolint:errcheck,forcetypeassert // only one type
		if time.Since(checkPost.lastUpdate()) > ps.pruneMaxAge {
			deletePosts = append(deletePosts, checkPost)
		}
	}

	for _, p := range deletePosts {
		ps.delete(p)
	}

	afterCount := ps.posts.Len()
	slog.InfoContext(ctx, "prune end", "prunedCount", beforeCount-afterCount)
}
