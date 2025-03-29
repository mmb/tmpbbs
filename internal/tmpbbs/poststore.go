package tmpbbs

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type PostStore struct {
	uuidMap map[string]int
	posts   []*post
	mutex   sync.RWMutex
}

func NewPostStore(title string) *PostStore {
	postStore := &PostStore{
		uuidMap: make(map[string]int),
	}
	postStore.put(newPost(title, "", "", nil), "")

	return postStore
}

func (ps *PostStore) put(post *post, parentUUID string) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	post.Parent = ps.getPostByUUID(parentUUID)
	if post.Parent != nil {
		post.Parent.Replies = append(post.Parent.Replies, post)
	}

	ps.uuidMap[post.uuid] = len(ps.posts)
	ps.posts = append(ps.posts, post)
}

func (ps *PostStore) get(uuid string, callback func(*post)) bool {
	var post *post

	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	if uuid == "" {
		post = ps.posts[0]
	} else if post = ps.getPostByUUID(uuid); post == nil {
		return false
	}

	callback(post)

	return true
}

func (ps *PostStore) getSince(uuid string, maxPosts int) []*post {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	sinceIndex, found := ps.uuidMap[uuid]
	if !found {
		sinceIndex = -1
	}

	start := sinceIndex + 1
	end := min(len(ps.posts), start+maxPosts)

	return ps.posts[start:end]
}

func (ps *PostStore) hasPost(uuid string) bool {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	_, found := ps.uuidMap[uuid]

	return found
}

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
		ps.put(newPost(posts[i].Title, posts[i].Author, posts[i].Body, tripcoder), ps.posts[0].uuid)
	}

	return nil
}

func (ps *PostStore) getPostByUUID(uuid string) *post {
	postIndex, found := ps.uuidMap[uuid]
	if !found {
		return nil
	}

	return ps.posts[postIndex]
}
