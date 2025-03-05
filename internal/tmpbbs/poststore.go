package tmpbbs

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type PostStore struct {
	posts []*post
	mutex sync.RWMutex
}

func NewPostStore(title string) *PostStore {
	return &PostStore{
		posts: []*post{newPost(title, "", "", nil)},
	}
}

func (ps *PostStore) put(post *post, parentID int) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	post.id = len(ps.posts)
	post.Parent = ps.posts[parentID]
	post.Parent.Replies = append(post.Parent.Replies, post)
	ps.posts = append(ps.posts, post)
}

func (ps *PostStore) get(postID int, callback func(*post)) bool {
	if postID < 0 {
		return false
	}

	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	if postID > len(ps.posts)-1 {
		return false
	}

	callback(ps.posts[postID])

	return true
}

func (ps *PostStore) LoadYAML(path string, tripCoder *TripCoder) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var posts []post

	if yamlErr := yaml.Unmarshal(data, &posts); yamlErr != nil {
		return yamlErr
	}

	for i := range posts {
		ps.put(newPost(posts[i].Title, posts[i].Author, posts[i].Body, tripCoder), 0)
	}

	return nil
}
