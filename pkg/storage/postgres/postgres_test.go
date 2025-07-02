package postgres

import (
	cryptoRand "crypto/rand"
	"math/rand"
	"news-aggregator/pkg/storage"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New("postgres://postgres:postgres@localhost:5432/posts")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDB_Posts(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	posts := []storage.Post{
		{
			Title:   "Test Post",
			Content: cryptoRand.Text(),
		},
	}
	db, err := New("postgres://postgres:postgres@localhost:5432/posts")
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddPosts(posts)
	if err != nil {
		t.Fatal(err)
	}
	allPosts, err := db.Posts(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", allPosts)
}
