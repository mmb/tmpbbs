package tmpbbs

import (
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// A PostStore stores Posts in memory and provides safety for concurrent
// access.
type PostStore struct {
	idMap  map[string]int
	rootID string
	posts  []*post
	mutex  sync.RWMutex
}

// NewPostStore returns a new PostStore. It also creates the root Post.
func NewPostStore(title string) *PostStore {
	rootPost := newPost(title, "", "", nil)
	postStore := &PostStore{
		idMap:  make(map[string]int),
		rootID: rootPost.id,
	}
	postStore.put(rootPost, "")

	return postStore
}

func (ps *PostStore) put(post *post, parentID string) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	post.Parent = ps.getPostByID(parentID)
	if post.Parent != nil {
		post.Parent.Replies.PushFront(post)
	}

	ps.idMap[post.id] = len(ps.posts)
	ps.posts = append(ps.posts, post)

	if (post.IsOriginalPoster() || post.IsSuperuser()) && strings.HasPrefix(post.Body, "!delete") {
		post.Parent.delete()
	} else {
		post.bump()
	}
}

func (ps *PostStore) get(postID string, callback func(*post)) bool {
	var post *post

	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	if postID == "" {
		post = ps.posts[0]
	} else if post = ps.getPostByID(postID); post == nil {
		return false
	}

	callback(post)

	return true
}

func (ps *PostStore) getSince(postID string, maxPosts int) []*post {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	sinceIndex, found := ps.idMap[postID]
	if !found {
		sinceIndex = -1
	}

	start := sinceIndex + 1
	end := min(len(ps.posts), start+maxPosts)

	return ps.posts[start:end]
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
		ps.put(newPost(posts[i].Title, posts[i].Author, posts[i].Body, tripcoder), ps.posts[0].id)
	}

	return nil
}

func (ps *PostStore) getPostByID(postID string) *post {
	postIndex, found := ps.idMap[postID]
	if !found {
		return nil
	}

	return ps.posts[postIndex]
}
