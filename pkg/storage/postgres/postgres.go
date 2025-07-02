package postgres

import (
	"context"
	"news-aggregator/pkg/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func New(dbURL string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Posts(n int) ([]storage.Post, error) {
	if n == 0 {
		n = 30
	}
	rows, err := s.db.Query(context.Background(), `
	SELECT p.id, p.title, p.content, p.pub_time, p.link 
	FROM posts p
	LIMIT $1
	`,
		n,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []storage.Post

	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, rows.Err()
}

func (s *Store) AddPost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), ` 
	INSERT INTO posts (title, content, pub_time, link)
	VALUES ($1, $2, $3, $4)
	`,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
	)

	return err
}

func (s *Store) AddPosts(p []storage.Post) error {
	for _, post := range p {
		err := s.AddPost(post)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), ` 
	UPDATE posts
	SET  title = $1, content = $2, pub_time = $3, link = $4
	WHERE id = $5
	`,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
		p.ID,
	)
	return err
}

func (s *Store) DeletePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
    DELETE FROM posts 
	WHERE id = $1
	`, p.ID)
	return err
}
