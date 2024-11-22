package tmpbbs

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type postStore struct {
	posts []*post
	mutex sync.RWMutex
}

func NewPostStore(title string) *postStore {
	return &postStore{
		posts: []*post{newPost(title, "", "", nil)},
	}
}

func (ps *postStore) put(post *post, parentID int) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	post.id = len(ps.posts)
	post.Parent = ps.posts[parentID]
	post.Parent.Replies = append(post.Parent.Replies, post)
	ps.posts = append(ps.posts, post)
}

func (ps *postStore) get(id int, callback func(*post, *post)) bool {
	if id < 0 {
		return false
	}

	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	if id > len(ps.posts)-1 {
		return false
	}
	callback(ps.posts[0], ps.posts[id])

	return true
}

func (ps *postStore) LoadYAML(path string, tripCoder *tripCoder) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var posts []post
	err = yaml.Unmarshal(data, &posts)
	if err != nil {
		return err
	}

	for _, post := range posts {
		ps.put(newPost(post.Title, post.Author, post.Body, tripCoder), 0)
	}

	return nil
}
